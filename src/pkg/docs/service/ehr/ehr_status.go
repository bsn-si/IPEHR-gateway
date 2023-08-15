package ehr

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/internal/observability/tracer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/access"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	proc "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	docStatus "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/status"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/types"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/ehrIndexer"
)

func (s *Service) CreateStatus(ehrStatusID string, subject base.PartySelf) (doc *model.EhrStatus, err error) {
	doc = &model.EhrStatus{
		Locatable: base.Locatable{
			Type:            base.EHRStatusItemType,
			ArchetypeNodeID: "openEHR-EHR-EHR_STATUS.generic.v1",
			Name:            base.NewDvText("EHR Status"),
			ObjectVersionID: base.ObjectVersionID{
				// todo FIXIT
				UID: &base.UIDBasedID{
					ObjectID: base.ObjectID{
						Type:  "OBJECT_VERSION_ID",
						Value: ehrStatusID,
					},
				},
			},
		},
		Subject:     subject,
		IsQueryable: true,
		IsModifable: true,
	}

	return doc, nil
}

func (s *Service) SaveStatus(ctx context.Context, multiCallTx *indexer.MultiCallTx, procRequest *proc.Request, userID, systemID string, ehrUUID *uuid.UUID, status *model.EhrStatus, allDocsGroup *model.DocumentGroup) error {
	ctx, span := tracer.Start(ctx, "ehr_service.SaveStatus")
	defer span.End()

	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	// Document encryption key generation
	key := chachaPoly.GenerateKey()

	objectVersionID, err := base.NewObjectVersionID(status.UID.Value, systemID)
	if err != nil {
		return fmt.Errorf("SaveStatus error: %w versionUID %s ehrSystemID %s", err, objectVersionID.String(), systemID)
	}

	statusBytes, err := json.Marshal(status)
	if err != nil {
		return fmt.Errorf("json.Marshal error: %w", err)
	}

	if s.Infra.CompressionEnabled {
		statusBytes, err = s.Infra.Compressor.Compress(statusBytes)
		if err != nil {
			return fmt.Errorf("Compress error: %w", err)
		}
	}

	// Document encryption
	statusEncrypted, err := key.EncryptWithAuthData(statusBytes, []byte(status.UID.Value))
	if err != nil {
		return fmt.Errorf("EncryptWithAuthData error: %w", err)
	}

	// IPFS saving
	CID, err := s.Infra.IpfsClient.Add(ctx, statusEncrypted)
	if err != nil {
		return fmt.Errorf("IpfsClient.Add error: %w", err)
	}

	// Filecoin saving
	dealCID, minerAddr, err := s.Infra.FilecoinClient.StartDeal(ctx, CID, uint64(len(statusEncrypted)))
	if err != nil {
		return fmt.Errorf("FilecoinClient.StartDeal error: %w", err)
	}
	//dealCID := fakeData.Cid()
	//minerAddr := []byte("123")

	procRequest.AddFilecoinTx(proc.TxSaveEhrStatus, CID.String(), dealCID.String(), minerAddr)

	// Index subject and namespace
	{
		subjectID := status.Subject.ExternalRef.ID.Value
		subjectNamespace := status.Subject.ExternalRef.Namespace

		setSubjectPacked, err := s.Infra.Index.SetEhrSubject(ctx, ehrUUID, subjectID, subjectNamespace, userPrivKey)
		if err != nil {
			return fmt.Errorf("Index.SetSubject error: %w ehrID: %s subjectID: %s subjectNamespace: %s", err, ehrUUID.String(), subjectID, subjectNamespace)
		}

		multiCallTx.Add(uint8(proc.TxSetEhrBySubject), setSubjectPacked)
	}

	err = s.addEhrStatusMetaData(multiCallTx, key, objectVersionID, CID, dealCID, minerAddr, userPubKey, userPrivKey)
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

func (s *Service) addEhrStatusMetaData(multiCallTx *indexer.MultiCallTx, key *chachaPoly.Key, objectVersionID *base.ObjectVersionID, CID, dealCID *cid.Cid, minerAddr string, userPubKey, userPrivKey *[32]byte) error {
	statusIDEncrypted, err := key.Encrypt([]byte(objectVersionID.String()))
	if err != nil {
		return fmt.Errorf("Encrypt error: %w objectVersionID %s", err, objectVersionID.String())
	}

	CIDEncr, err := key.Encrypt(CID.Bytes())
	if err != nil {
		return fmt.Errorf("CID encryption error: %w", err)
	}

	keyEncr, err := keybox.SealAnonymous(key.Bytes(), userPubKey)
	if err != nil {
		return fmt.Errorf("keybox.SealAnonymous error: %w", err)
	}

	baseDocumentUID := []byte(objectVersionID.BasedID())
	baseDocumentUIDHash := sha3.Sum256(baseDocumentUID)

	docMeta := &model.DocumentMeta{
		Status:    uint8(docStatus.ACTIVE),
		Id:        CID.Bytes(),
		Version:   objectVersionID.VersionBytes()[:],
		Timestamp: uint32(time.Now().Unix()),
		IsLast:    true,
		Attrs: []ehrIndexer.AttributesAttribute{
			{Code: model.AttributeIDEncr, Value: CIDEncr},
			{Code: model.AttributeKeyEncr, Value: keyEncr},
			{Code: model.AttributeDocUIDHash, Value: baseDocumentUIDHash[:]},
			{Code: model.AttributeDocUIDEncr, Value: statusIDEncrypted},
			{Code: model.AttributeDealCid, Value: dealCID.Bytes()},
			{Code: model.AttributeMinerAddress, Value: []byte(minerAddr)},
		},
	}

	packed, err := s.Infra.Index.AddEhrDoc(types.EhrStatus, docMeta, userPrivKey)
	if err != nil {
		return fmt.Errorf("Index.AddEhrDoc error: %w", err)
	}

	multiCallTx.Add(uint8(proc.TxAddEhrDoc), packed)

	return nil
}

func (s *Service) UpdateStatus(ctx context.Context, procRequest *proc.Request, userID, systemID string, ehrUUID *uuid.UUID, status *model.EhrStatus) error {
	ctx, span := tracer.Start(ctx, "ehr_service.UpdateStatus")
	defer span.End()

	// Searching 'all documents' group
	var allDocGroup *model.DocumentGroup
	{
		docGroups, err := s.DocGroup.GroupGetList(ctx, userID, systemID)
		if err != nil {
			return fmt.Errorf("DocGroup.GroupGetList error: %w", err)
		}

		for _, dg := range docGroups {
			if dg.Name == common.DefaultGroupAllDocuments {
				allDocGroup = dg
				break
			}
		}

		if allDocGroup == nil {
			return fmt.Errorf("user 'all documents' group not found: %w", errors.ErrNotFound)
		}
	}

	multiCallTx := s.Infra.Index.MultiCallEhrNew()

	err := s.SaveStatus(ctx, multiCallTx, procRequest, userID, systemID, ehrUUID, status, allDocGroup)
	if err != nil {
		return fmt.Errorf("SaveStatus error: %w", err)
	}

	// TODO i dont like this logic, because in method GetByID we always grab whole data from filecoin, which contain last status id. It need fix it.
	if err := s.UpdateEhr(ctx, multiCallTx, procRequest, userID, systemID, ehrUUID, status, allDocGroup); err != nil {
		return fmt.Errorf("UpdateEhr error: %w", err)
	}

	txHash, err := multiCallTx.Commit()
	if err != nil {
		return fmt.Errorf("UpdateStatus commit error: %w", err)
	}

	for _, txKind := range multiCallTx.GetTxKinds() {
		procRequest.AddEthereumTx(proc.TxKind(txKind), txHash)
	}

	return nil
}

// GetStatus Get current (last) status of EHR document
func (s *Service) GetStatus(ctx context.Context, userID, systemID string, ehrUUID *uuid.UUID) (*model.EhrStatus, error) {
	ctx, span := tracer.Start(ctx, "ehr_service.GetStatus")
	defer span.End()

	docMeta, err := s.Infra.Index.GetDocLastByType(ctx, ehrUUID, types.EhrStatus)
	if err != nil {
		return nil, fmt.Errorf("Index.GetLastEhrDocByType error: %w. ehrID: %s", err, ehrUUID.String())
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

	var status model.EhrStatus
	if err := json.Unmarshal(docDecrypted, &status); err != nil {
		return nil, fmt.Errorf("EHR status unmarshal error: %w", err)
	}

	return &status, nil
}

func (s *Service) GetStatusByVersionID(ctx context.Context, userID, systemID string, ehrUUID *uuid.UUID, versionID *base.ObjectVersionID) ([]byte, error) {
	ctx, span := tracer.Start(ctx, "ehr_service.GetStatusByVersionID")
	defer span.End()

	baseDocumentUID := versionID.BasedID()
	baseDocumentUIDHash := sha3.Sum256([]byte(baseDocumentUID))

	docMeta, err := s.Infra.Index.GetDocByVersion(ctx, ehrUUID, types.EhrStatus, &baseDocumentUIDHash, versionID.VersionBytes())
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("Index.GetDocByVersion error: %w", err)
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
	if err != nil {
		if errors.Is(err, errors.ErrIsInProcessing) {
			return nil, err
		} else if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("GetDocFromStorageByID error: %w", err)
	}

	return docDecrypted, nil
}

func (s *Service) GetStatusByNearestTime(ctx context.Context, userID, systemID string, ehrUUID *uuid.UUID, nearestTime time.Time) ([]byte, error) {
	ctx, span := tracer.Start(ctx, "ehr_service.GetStatusByNearestTime")
	defer span.End()

	docMeta, err := s.Infra.Index.GetDocByTime(ctx, ehrUUID, types.EhrStatus, uint32(nearestTime.Unix()))
	if err != nil && errors.Is(err, errors.ErrNotFound) {
		return nil, err
	} else if err != nil {
		return nil, fmt.Errorf("Index.GetDocByTime error: %w ehrID %s nearestTime %s docType %s", err, ehrUUID.String(), nearestTime.String(), types.EhrStatus)
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

func (s *Service) ValidateStatus(status *model.EhrStatus) bool {
	// TODO
	return true
}
