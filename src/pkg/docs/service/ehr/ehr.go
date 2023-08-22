package ehr

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	"github.com/vmihailenco/msgpack/v5"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/internal/observability/tracer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/access"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service"
	docService "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service"
	docGroupService "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/docGroup"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/groupAccess"
	proc "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	docStatus "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/status"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/types"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/ehrIndexer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/infrastructure"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/storage/treeindex"
	userModel "github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
	userService "github.com/bsn-si/IPEHR-gateway/src/pkg/user/service"
)

type Service struct {
	Infra       *infrastructure.Infra
	Doc         *service.DefaultDocumentService
	DocGroup    *docGroupService.Service
	User        *userService.Service
	GroupAccess *groupAccess.Service
}

func NewService(docSvc *docService.DefaultDocumentService, userSvc *userService.Service, docGroupSvc *docGroupService.Service, gaSvc *groupAccess.Service) *Service {
	return &Service{
		Infra:       docSvc.Infra,
		Doc:         docSvc,
		DocGroup:    docGroupSvc,
		User:        userSvc,
		GroupAccess: gaSvc,
	}
}

func (s *Service) EhrCreate(ctx context.Context, userID, systemID string, ehrUUID, groupAccessUUID *uuid.UUID, request *model.EhrCreateRequest, procRequest *proc.Request) (*model.EHR, error) {
	return s.EhrCreateWithID(ctx, userID, systemID, ehrUUID, groupAccessUUID, request, procRequest)
}

func (s *Service) EhrCreateWithID(ctx context.Context, userID, systemID string, ehrUUID, groupAccessUUID *uuid.UUID, request *model.EhrCreateRequest, procRequest *proc.Request) (*model.EHR, error) {
	ctx, span := tracer.Start(ctx, "ehr_service.EhrCreateWithID", trace.WithAttributes([]attribute.KeyValue{
		attribute.String("userID", userID),
		attribute.String("ehrUUID", ehrUUID.String()),
		attribute.String("systemID", systemID),
	}...))
	defer span.End()

	ehr := model.EHR{
		SystemID: base.HierObjectID{
			UIDBasedID: base.UIDBasedID{
				ObjectID: base.ObjectID{
					Type:  "HIER_OBJECT_ID",
					Value: systemID,
				},
			},
		},
		EhrID: base.HierObjectID{
			UIDBasedID: base.UIDBasedID{
				ObjectID: base.ObjectID{
					Type:  "HIER_OBJECT_ID",
					Value: ehrUUID.String(),
				},
			},
		},
		EhrAccess: base.ObjectRef{
			Namespace: "local",
			Type:      "EHR_ACCESS",
			ID: base.ObjectID{
				Type:  "OBJECT_VERSION_ID",
				Value: uuid.New().String() + "::" + systemID + "::1",
			},
		},

		TimeCreated: base.DvDateTime{
			Value: time.Now().Format(common.OpenEhrTimeFormat),
		},
	}

	// Creating EHR_STATUS
	ehrStatusID := uuid.New().String() + "::" + systemID + "::1"
	subjectID := request.Subject.ExternalRef.ID.Value
	subjectNamespace := request.Subject.ExternalRef.Namespace

	subject := s.CreateSubject(subjectID, subjectNamespace, "PERSON")

	doc, err := s.CreateStatus(ehrStatusID, subject)
	if err != nil {
		return nil, fmt.Errorf("create status error: %w", err)
	}

	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	multiCallTx := s.Infra.Index.MultiCallEhrNew()

	// Index EHR userIDHash -> ehrUUID
	packed, err := s.Infra.Index.SetEhrUser(ctx, userID, systemID, ehrUUID, userPrivKey)
	if err != nil {
		return nil, fmt.Errorf("Index.SetEhrUser error: %w", err)
	}

	multiCallTx.Add(uint8(proc.TxSetEhrUser), packed)

	// 'All documents' docGroup creating
	allDocsGroup := &model.DocumentGroup{
		GroupID:  uuid.New(),
		GroupKey: chachaPoly.GenerateKey(),
		Name:     common.DefaultGroupAllDocuments,
	}

	var (
		allDocsGroupIDHash = indexer.Keccak256(allDocsGroup.GroupID[:])
		groupIDEncr        []byte
		groupKeyEncr       []byte
	)

	{
		groupIDEncr, err = allDocsGroup.GroupKey.Encrypt(allDocsGroup.GroupID[:])
		if err != nil {
			return nil, fmt.Errorf("allDocsGroupID encryption error: %w", err)
		}

		groupNameEncr, err := allDocsGroup.GroupKey.Encrypt([]byte(allDocsGroup.Name))
		if err != nil {
			return nil, fmt.Errorf("groupName encryption error: %w", err)
		}

		groupKeyEncr, err = keybox.SealAnonymous(allDocsGroup.GroupKey.Bytes(), userPubKey)
		if err != nil {
			return nil, fmt.Errorf("keybox.SealAnonymous error: %w", err)
		}

		packed, err := s.Infra.Index.DocGroupCreate(ctx, &allDocsGroup.GroupID, groupIDEncr, groupKeyEncr, groupNameEncr, userPrivKey)
		if err != nil {
			return nil, fmt.Errorf("Index.DocGroupCreate error: %w", err)
		}

		multiCallTx.Add(uint8(proc.TxDocGroupCreate), packed)
	}

	// Granting access to the group 'All documents' for user
	{
		userIDHash := sha3.Sum256([]byte(userID + systemID))

		accessObj := indexer.AccessObject{
			Kind:    access.DocGroup,
			IdHash:  *allDocsGroupIDHash,
			IdEncr:  groupIDEncr,
			KeyEncr: groupKeyEncr,
			Level:   access.Owner,
		}

		setAccessPacked, err := s.Infra.Index.SetAccessWrapper(ctx, &userIDHash, &accessObj, userPrivKey)
		if err != nil {
			return nil, fmt.Errorf("Index.SetAccess user to allDocsGroup error: %w", err)
		}

		multiCallTx.Add(uint8(proc.TxSetDocAccess), setAccessPacked)
	}

	err = s.SaveStatus(ctx, multiCallTx, procRequest, userID, systemID, ehrUUID, doc, allDocsGroup)
	if err != nil {
		return nil, fmt.Errorf("SaveStatus error: %w. ehrID: %s userID: %s", err, ehrUUID.String(), userID)
	}

	ehr.EhrStatus.ID = doc.UID.ObjectID
	ehr.EhrStatus.Type = "EHR_STATUS"

	err = s.SaveEhr(ctx, multiCallTx, procRequest, userID, systemID, &ehr, allDocsGroup)
	if err != nil {
		return nil, fmt.Errorf("SaveEhr error: %w", err)
	}

	// Granting access to the group 'All documents' for the 'Doctors' group
	{
		userGroupList, err := s.User.GroupGetList(ctx, userID, systemID)
		if err != nil {
			return nil, fmt.Errorf("DocAccess.List error: %w", err)
		}

		var doctorsGroup *userModel.UserGroup

		for _, ug := range userGroupList {
			if ug.Name == common.DefaultGroupDoctors {
				doctorsGroup = ug
				break
			}
		}

		if doctorsGroup == nil {
			return nil, fmt.Errorf("user default group 'doctors' %w", errors.ErrNotFound)
		}

		doctorsGroupIDHash := indexer.Keccak256(doctorsGroup.GroupID[:])

		doctorsGroupKey, err := chachaPoly.NewKeyFromBytes(doctorsGroup.Key[:])
		if err != nil {
			return nil, fmt.Errorf("chachaPoly.NewKeyFromBytes error: %w", err)
		}

		groupIDEncr, err := allDocsGroup.GroupKey.Encrypt(allDocsGroup.GroupID[:])
		if err != nil {
			return nil, fmt.Errorf("doctorsGroup.GroupID encrypt error: %w", err)
		}

		groupKeyEncr, err := doctorsGroupKey.Encrypt(allDocsGroup.GroupKey.Bytes())
		if err != nil {
			return nil, fmt.Errorf("allDocsGroupKey encryption error: %w", err)
		}

		accessObj := indexer.AccessObject{
			Kind:    access.DocGroup,
			IdHash:  *allDocsGroupIDHash,
			IdEncr:  groupIDEncr,
			KeyEncr: groupKeyEncr,
			Level:   access.Read,
		}

		setAccessPacked, err := s.Infra.Index.SetAccessWrapper(ctx, doctorsGroupIDHash, &accessObj, userPrivKey)
		if err != nil {
			return nil, fmt.Errorf("Index.SetAccess doctorsGroup to allDocsGroup error: %w", err)
		}

		multiCallTx.Add(uint8(proc.TxSetDocAccess), setAccessPacked)
	}

	txHash, txNonce, err := multiCallTx.Commit()
	if err != nil {
		return nil, fmt.Errorf("EhrCreateWithID commit error: %w", err)
	}

	for _, txKind := range multiCallTx.GetTxKinds() {
		procRequest.AddEthereumTx(proc.TxKind(txKind), txHash, txNonce)
	}

	// Adding dataStore index
	err = s.addDataIndex(ctx, ehrUUID, groupAccessUUID, &uuid.UUID{}, &ehr, procRequest)
	if err != nil {
		return nil, fmt.Errorf("addDataIndex error: %w", err)
	}

	return &ehr, nil
}

func (s *Service) SaveEhr(ctx context.Context, multiCallTx *indexer.MultiCallTx, procRequest *proc.Request, userID, systemID string, doc *model.EHR, allDocsGroup *model.DocumentGroup) error {
	ctx, span := tracer.Start(ctx, "ehr_service.SaveEhr", trace.WithAttributes([]attribute.KeyValue{
		attribute.String("userID", userID),
		attribute.String("ehrID", doc.EhrID.Value),
		attribute.String("systemID", systemID),
	}...))
	defer span.End()

	ehrUUID, err := uuid.Parse(doc.EhrID.Value)
	if err != nil {
		return fmt.Errorf("ehrUUID parse error: %w ehrID.Value %s", err, doc.EhrID.Value)
	}

	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	docBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("ehr marshal error: %w", err)
	}

	if s.Infra.CompressionEnabled {
		docBytes, err = s.Infra.Compressor.Compress(docBytes)
		if err != nil {
			return fmt.Errorf("ehr compress error: %w", err)
		}
	}

	// Document encryption key generation
	key := chachaPoly.GenerateKey()

	// Document encryption
	docEncrypted, err := key.EncryptWithAuthData(docBytes, ehrUUID[:])
	if err != nil {
		return fmt.Errorf("ehr encryption error: %w", err)
	}

	// IPFS saving
	CID, err := s.Infra.IpfsClient.Add(ctx, docEncrypted)
	if err != nil {
		return fmt.Errorf("IpfsClient.Add error: %w", err)
	}

	// Filecoin saving
	dealCID, minerAddr, err := s.Infra.FilecoinClient.StartDeal(ctx, CID, uint64(len(docEncrypted)))
	if err != nil {
		return fmt.Errorf("FilecoinClient.StartDeal error: %w", err)
	}
	//dealCID := fakeData.Cid()
	//minerAddr := []byte("123")

	procRequest.AddFilecoinTx(proc.TxSaveEhr, CID.String(), dealCID.String(), minerAddr)

	err = s.addEhrMetaData(multiCallTx, key, &ehrUUID, CID, dealCID, minerAddr, userPubKey, userPrivKey)
	if err != nil {
		return fmt.Errorf("addMetaData error: %w", err)
	}

	err = s.setDocAccess(ctx, multiCallTx, userID, systemID, CID, key, access.Owner, userPubKey, userPrivKey)
	if err != nil {
		return fmt.Errorf("setDocAccess error: %w", err)
	}

	err = s.addDocumentToGroup(multiCallTx, CID, allDocsGroup, userPrivKey)
	if err != nil {
		return fmt.Errorf("addDocumentToGroup error: %w", err)
	}

	return nil
}

func (s *Service) addEhrMetaData(multiCallTx *indexer.MultiCallTx, key *chachaPoly.Key, ehrUUID *uuid.UUID, CID, dealCID *cid.Cid, minerAddr string, userPubKey, userPrivKey *[32]byte) error {
	ehrIDEncrypted, err := key.Encrypt(ehrUUID[:])
	if err != nil {
		return fmt.Errorf("Encrypt error: %w ehrID: %s", err, ehrUUID.String())
	}

	CIDEncr, err := key.Encrypt(CID.Bytes())
	if err != nil {
		return fmt.Errorf("CID encryption error: %w", err)
	}

	keyEncr, err := keybox.SealAnonymous(key.Bytes(), userPubKey)
	if err != nil {
		return fmt.Errorf("keybox.SealAnonymous error: %w", err)
	}

	docMeta := &model.DocumentMeta{
		Status:    uint8(docStatus.ACTIVE),
		Id:        CID.Bytes(),
		Version:   nil,
		Timestamp: uint32(time.Now().Unix()),
		IsLast:    true,
		Attrs: []ehrIndexer.AttributesAttribute{
			{Code: model.AttributeIDEncr, Value: CIDEncr},
			{Code: model.AttributeKeyEncr, Value: keyEncr},
			{Code: model.AttributeDocUIDHash, Value: make([]byte, 32)},
			{Code: model.AttributeDocUIDEncr, Value: ehrIDEncrypted},
			{Code: model.AttributeDealCid, Value: dealCID.Bytes()},
			{Code: model.AttributeMinerAddress, Value: []byte(minerAddr)},
		},
	}

	packed, err := s.Infra.Index.AddEhrDoc(types.Ehr, docMeta, userPrivKey)
	if err != nil {
		return fmt.Errorf("Index.AddEhrDoc error: %w", err)
	}

	multiCallTx.Add(uint8(proc.TxAddEhrDoc), packed)

	return nil
}

func (s *Service) setDocAccess(ctx context.Context, multiCallTx *indexer.MultiCallTx, userID, systemID string, CID *cid.Cid, key *chachaPoly.Key, accessLevel access.Level, userPubKey, userPrivKey *[32]byte) error {
	ctx, span := tracer.Start(ctx, "ehr_service.setDocAccess", trace.WithAttributes([]attribute.KeyValue{
		attribute.String("userID", userID),
		attribute.String("CID", CID.String()),
		attribute.String("systemID", systemID),
	}...))
	defer span.End()

	userIDHash := sha3.Sum256([]byte(userID + systemID))
	docIDHash := indexer.Keccak256(CID.Bytes())

	CIDEncr, err := key.Encrypt(CID.Bytes())
	if err != nil {
		return fmt.Errorf("CID encryption error: %w", err)
	}

	keyEncr, err := keybox.SealAnonymous(key.Bytes(), userPubKey)
	if err != nil {
		return fmt.Errorf("keybox.SealAnonymous error: %w", err)
	}

	accessObj := indexer.AccessObject{
		Kind:    access.Doc,
		IdHash:  *docIDHash,
		IdEncr:  CIDEncr,
		KeyEncr: keyEncr,
		Level:   accessLevel,
	}

	packed, err := s.Infra.Index.SetAccessWrapper(ctx, &userIDHash, &accessObj, userPrivKey)
	if err != nil {
		return fmt.Errorf("Index.SetAccess user to allDocsGroup error: %w", err)
	}

	multiCallTx.Add(uint8(proc.TxSetDocAccess), packed)

	return nil
}

func (s *Service) addDocumentToGroup(multiCallTx *indexer.MultiCallTx, CID *cid.Cid, groupToAdd *model.DocumentGroup, userPrivKey *[32]byte) error {
	docCIDHash := indexer.Keccak256(CID.Bytes())

	docCIDEncr, err := groupToAdd.GroupKey.Encrypt(CID.Bytes())
	if err != nil {
		return fmt.Errorf("EHR_STATUS CID encryption error: %w", err)
	}

	packed, err := s.Infra.Index.DocGroupAddDoc(&groupToAdd.GroupID, docCIDHash, docCIDEncr, userPrivKey)
	if err != nil {
		return fmt.Errorf("Index.DocGroupAddDoc error: %w", err)
	}

	multiCallTx.Add(uint8(proc.TxDocGroupAddDoc), packed)

	return nil
}

func (s *Service) GetByID(ctx context.Context, userID, systemID string, ehrUUID *uuid.UUID) ([]byte, error) {
	docMeta, err := s.Infra.Index.GetDocLastByType(ctx, ehrUUID, types.Ehr)
	if err != nil {
		return nil, fmt.Errorf("GetDocLastByType error: %w", err)
	}

	CID, err := cid.Parse(docMeta.Id)
	if err != nil {
		return nil, fmt.Errorf("cid.Parse error: %w", err)
	}

	docUIDEncrypted := docMeta.GetAttr(model.AttributeDocUIDEncr)
	if docUIDEncrypted == nil {
		return nil, errors.ErrFieldIsEmpty("DocUIDEncrypted")
	}

	docDecrypted, err := s.Doc.GetDocFromStorageByID(ctx, userID, systemID, &CID, ehrUUID[:], docUIDEncrypted)
	if err != nil && errors.Is(err, errors.ErrIsInProcessing) {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("GetDocFromStorageByID error: %w", err)
	}

	return docDecrypted, nil
}

// GetDocBySubject Get decrypted document by subject
func (s *Service) GetDocBySubject(ctx context.Context, userID, systemID, subjectID, namespace string) (docDecrypted []byte, err error) {
	ctx, span := tracer.Start(ctx, "ehr_service.GetDocBySubject", trace.WithAttributes([]attribute.KeyValue{
		attribute.String("userID", userID),
		attribute.String("systemID", systemID),
		attribute.String("subjectID", subjectID),
		attribute.String("namespace", namespace),
	}...))
	defer span.End()

	ehrUUID, err := s.Infra.Index.GetEhrUUIDBySubject(ctx, subjectID, namespace)
	if err != nil {
		return nil, fmt.Errorf("Index.GetEhrUUIDBySubject error: %w. userID: %s subjectID: %s namespace: %s", err, userID, subjectID, namespace)
	}

	// Getting docStorageID
	docMeta, err := s.Infra.Index.GetDocLastByType(ctx, ehrUUID, types.Ehr)
	if err != nil {
		return nil, fmt.Errorf("Index.GetLastDocByType error: %w. ehrUUID: %s", err, ehrUUID.String())
	}

	CID, err := cid.Parse(docMeta.Id)
	if err != nil {
		return nil, fmt.Errorf("cid.Parse error: %w", err)
	}

	docUIDEncrypted := docMeta.GetAttr(model.AttributeDocUIDEncr)
	if docUIDEncrypted == nil {
		return nil, errors.ErrFieldIsEmpty("DocUIDEncrypted")
	}

	// Getting doc from storage
	docDecrypted, err = s.Doc.GetDocFromStorageByID(ctx, userID, systemID, &CID, ehrUUID[:], docUIDEncrypted)
	if err != nil && errors.Is(err, errors.ErrIsInProcessing) {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("GetDocFromStorageByID error: %w. userID: %s, doc.CID: %s ehrUUID: %s", err, userID, CID.String(), ehrUUID.String())
	}

	return docDecrypted, nil
}

func (s *Service) CreateSubject(subjectID, subjectNamespace, subType string) (subject base.PartySelf) {
	subject.ExternalRef = &base.ObjectRef{
		ID: base.ObjectID{
			Type:  "HIER_OBJECT_ID", // TODO is it always eq with "HIER_OBJECT_ID"?
			Value: subjectID,
		},
		Namespace: subjectNamespace,
		Type:      subType,
	}

	return
}

func (s *Service) UpdateEhr(ctx context.Context, multiCallTx *indexer.MultiCallTx, procRequest *proc.Request, userID, systemID string, ehrUUID *uuid.UUID, status *model.EhrStatus, allDocsGroup *model.DocumentGroup) error {
	ctx, span := tracer.Start(ctx, "ehr_service.UpdateEhr", trace.WithAttributes([]attribute.KeyValue{
		attribute.String("userID", userID),
		attribute.String("ehrUUID", ehrUUID.String()),
		attribute.String("systemID", systemID),
	}...))
	defer span.End()

	docMeta, err := s.Infra.Index.GetDocLastByType(ctx, ehrUUID, types.Ehr)
	if err != nil {
		return fmt.Errorf("Index.GetLastEhrDocByType error: %w. ehrID: %s", err, ehrUUID.String())
	}

	CID, err := cid.Parse(docMeta.Id)
	if err != nil {
		return fmt.Errorf("cid.Parse error: %w", err)
	}

	docUIDEncrypted := docMeta.GetAttr(model.AttributeDocUIDEncr)
	if docUIDEncrypted == nil {
		return errors.ErrFieldIsEmpty("DocUIDEncrypted")
	}

	ehrDecrypted, err := s.Doc.GetDocFromStorageByID(ctx, userID, systemID, &CID, ehrUUID[:], docUIDEncrypted)
	if err != nil && errors.Is(err, errors.ErrIsInProcessing) {
		return err
	} else if err != nil {
		return fmt.Errorf("GetDocFromStorageByID error: %w. userID: %s StorageID: %x ehrID: %s", err, userID, &CID, ehrUUID.String())
	}

	var ehr model.EHR
	if err = json.Unmarshal(ehrDecrypted, &ehr); err != nil {
		return fmt.Errorf("ehr unmarshal error: %w", err)
	}

	if status.UID.Value != ehr.EhrStatus.ID.Value {
		ehr.EhrStatus.ID.Value = status.UID.Value
		if err = s.SaveEhr(ctx, multiCallTx, procRequest, userID, systemID, &ehr, allDocsGroup); err != nil {
			return fmt.Errorf("ehr save error: %w", err)
		}
	}

	return nil
}

func (s *Service) ValidateEhr(ehr *model.EHR) bool {
	// TODO
	return true
}

func (s *Service) addDataIndex(ctx context.Context, ehrUUID, groupAccessUUID, dataIndexUUID *uuid.UUID, ehr *model.EHR, procRequest *proc.Request) error {
	ctx, span := tracer.Start(ctx, "ehr_service.addDataIndex", trace.WithAttributes([]attribute.KeyValue{
		attribute.String("ehrUUID", ehrUUID.String()),
		attribute.String("groupAccessUUID", groupAccessUUID.String()),
		attribute.String("dataIndexUUID", dataIndexUUID.String()),
	}...))
	defer span.End()

	ehrNode, err := treeindex.ProcessEHR(ehr)
	if err != nil {
		return fmt.Errorf("treeindex.ProcessEHR error: %w", err)
	}

	data, err := msgpack.Marshal(ehrNode)
	if err != nil {
		return fmt.Errorf("msgpack.Marshal(ehrNode) error: %w", err)
	}

	compressed, err := s.Infra.Compressor.Compress(data)
	if err != nil {
		return fmt.Errorf("data compressinon error: %w", err)
	}

	txHash, txNonce, err := s.Infra.Index.DataUpdate(ctx, groupAccessUUID, dataIndexUUID, ehrUUID, compressed)
	if err != nil {
		return fmt.Errorf("Index.DataUpdate error: %w", err)
	}

	procRequest.AddEthereumTx(proc.TxIndexDataUpdate, txHash, txNonce)

	return nil
}
