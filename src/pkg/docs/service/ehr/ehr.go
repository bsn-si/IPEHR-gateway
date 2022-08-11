package ehr

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/common"
	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/docs/types"
)

type Service struct {
	*service.DefaultDocumentService
}

func NewService(docService *service.DefaultDocumentService) *Service {
	return &Service{
		docService,
	}
}

func (s *Service) EhrCreate(ctx context.Context, userID string, ehrSystemID base.EhrSystemID, request *model.EhrCreateRequest) (*model.EHR, error) {
	ehrUUID := uuid.New()
	return s.EhrCreateWithID(ctx, userID, &ehrUUID, ehrSystemID, request)
}

func (s *Service) EhrCreateWithID(ctx context.Context, userID string, ehrUUID *uuid.UUID, ehrSystemID base.EhrSystemID, request *model.EhrCreateRequest) (*model.EHR, error) {
	var ehr model.EHR

	ehr.SystemID.Value = ehrSystemID.String()
	ehr.EhrID.Value = ehrUUID.String()

	ehr.EhrStatus.ID.Type = "OBJECT_VERSION_ID"
	ehr.EhrStatus.ID.Value = uuid.New().String() + "::" + ehrSystemID.String() + "::1"
	ehr.EhrStatus.Namespace = "local"
	ehr.EhrStatus.Type = "EHR_STATUS"

	ehr.EhrAccess.ID.Type = "OBJECT_VERSION_ID"
	ehr.EhrAccess.ID.Value = uuid.New().String() + "::" + ehrSystemID.String() + "::1"
	ehr.EhrAccess.Namespace = "local"
	ehr.EhrAccess.Type = "EHR_ACCESS"

	ehr.TimeCreated.Value = time.Now().Format(common.OpenEhrTimeFormat)

	err := s.SaveEhr(ctx, userID, &ehr)
	if err != nil {
		return nil, fmt.Errorf("ehr save error: %w", err)
	}

	// Creating EHR_STATUS
	ehrStatusID := ehr.EhrStatus.ID.Value
	subjectID := request.Subject.ExternalRef.ID.Value
	subjectNamespace := request.Subject.ExternalRef.Namespace

	_, err = s.CreateStatus(ctx, userID, ehrStatusID, subjectID, subjectNamespace, ehrUUID, ehrSystemID)
	if err != nil {
		return nil, fmt.Errorf("create status error: %w", err)
	}

	return &ehr, nil
}

func (s *Service) SaveEhr(ctx context.Context, userID string, doc *model.EHR) error {
	ehrUUID, err := uuid.Parse(doc.EhrID.Value)
	if err != nil {
		return fmt.Errorf("ehrUUID parse error: %w ehrID.Value %s", err, doc.EhrID.Value)
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
	docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(doc.EhrID.Value))
	if err != nil {
		return fmt.Errorf("ehr encryption error: %w", err)
	}

	// Storage saving
	cidBytes, err := s.Infra.IpfsClient.Add(docEncrypted)
	if err != nil {
		return fmt.Errorf("IpfsClient.Add error: %w", err)
	}

	// Start processing request
	reqID := ctx.(*gin.Context).GetString("reqId")
	{
		procReq := &processing.Request{
			ReqID:   reqID,
			UserID:  userID,
			EhrUUID: ehrUUID.String(),
			Kind:    processing.RequestEhrCreate,
			CID:     hex.EncodeToString(cidBytes[:]),
		}
		if err = s.Proc.AddRequest(procReq); err != nil {
			return fmt.Errorf("Proc.AddRequest error: %w", err)
		}
	}

	// Index EHR userID -> docStorageID
	{
		ehrIndexTx, err := s.Infra.Index.SetEhrUser(userID, &ehrUUID)
		if err != nil {
			return fmt.Errorf("Index.SetEhrUser error: %w", err)
		}

		err = s.Proc.AddBlockchainTx(reqID, ehrIndexTx, "", processing.TxSetEhrUser, processing.StatusPending)
		if err != nil {
			return fmt.Errorf("Proc.AddBlockchainTx error: %w", err)
		}
	}

	// Index Docs ehr_id -> doc_meta
	{
		docMeta := &model.DocumentMeta{
			TypeCode:  types.Ehr,
			CID:       cidBytes,
			Timestamp: uint64(time.Now().UnixNano()),
		}

		docIndexTx, err := s.Infra.Index.AddEhrDoc(&ehrUUID, docMeta)
		if err != nil {
			return fmt.Errorf("Index.AddEhrDoc error: %w", err)
		}

		err = s.Proc.AddBlockchainTx(reqID, docIndexTx, "", processing.TxSetEhrDocs, processing.StatusPending)
		if err != nil {
			return fmt.Errorf("Proc.AddBlockchainTx error: %w", err)
		}
	}

	// Index Access
	{
		userPubKey, _, err := s.Infra.Keystore.Get(userID)
		if err != nil {
			return fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
		}

		docAccessValue, err := keybox.SealAnonymous(key.Bytes(), userPubKey)
		if err != nil {
			return fmt.Errorf("keybox.SealAnonymous error: %w", err)
		}

		docAccessKey := sha3.Sum256(append(cidBytes[:], []byte(userID)...))

		docAccessTx, err := s.Infra.Index.SetDocKeyEncrypted(&docAccessKey, docAccessValue)
		if err != nil {
			return fmt.Errorf("Index.SetDocAccess error: %w", err)
		}

		err = s.Proc.AddBlockchainTx(reqID, docAccessTx, "", processing.TxSetDocAccess, processing.StatusPending)
		if err != nil {
			return fmt.Errorf("Proc.AddBlockchainTx error: %w", err)
		}
	}

	return nil
}

// GetDocBySubject Get decrypted document by subject
func (s *Service) GetDocBySubject(ctx context.Context, userID, subjectID, namespace string) (docDecrypted []byte, err error) {
	ehrUUID, err := s.Infra.Index.GetEhrUUIDBySubject(ctx, subjectID, namespace)
	if err != nil {
		return nil, fmt.Errorf("Index.GetEhrUUIDBySubject error: %w. userID: %s subjectID: %s namespace: %s", err, userID, subjectID, namespace)
	}

	// Getting docStorageID
	docMeta, err := s.Infra.Index.GetDocLastByType(ctx, ehrUUID, types.Ehr)
	if err != nil {
		return nil, fmt.Errorf("Index.GetLastDocByType error: %w. ehrUUID: %s", err, ehrUUID.String())
	}

	// Getting doc from storage
	docDecrypted, err = s.GetDocFromStorageByID(ctx, userID, docMeta.CID, ehrUUID[:], docMeta.DocUIDEncrypted)
	if err != nil {
		return nil, fmt.Errorf("GetDocFromStorageByID error: %w. userID: %s, doc.CID: %x ehrUUID: %s", err, userID, docMeta.CID, ehrUUID.String())
	}

	return docDecrypted, nil
}

func (s *Service) CreateStatus(ctx context.Context, userID, ehrStatusID, subjectID, subjectNamespace string, ehrUUID *uuid.UUID, ehrSystemID base.EhrSystemID) (doc *model.EhrStatus, err error) {
	doc = &model.EhrStatus{}
	doc.Type = types.EhrStatus.String()
	doc.ArchetypeNodeID = "openEHR-EHR-EHR_STATUS.generic.v1"
	doc.Name = base.DvText{Value: "EHR Status"}

	// FIXIT
	doc.UID = &base.UIDBasedID{ObjectID: base.ObjectID{
		Type:  "OBJECT_VERSION_ID",
		Value: ehrStatusID,
	}}

	doc.Subject.ExternalRef = base.ObjectRef{
		ID: base.ObjectID{
			Type:  "HIER_OBJECT_ID",
			Value: subjectID,
		},
		Namespace: subjectNamespace,
		Type:      "PERSON",
	}
	doc.IsQueryable = true
	doc.IsModifable = true

	err = s.SaveStatus(ctx, userID, ehrUUID, ehrSystemID, doc, true)
	if err != nil {
		return nil, fmt.Errorf("SaveStatus error: %w. ehrID: %s userID: %s", err, ehrUUID.String(), userID)
	}

	return doc, nil
}

func (s *Service) UpdateStatus(ctx context.Context, userID string, ehrUUID *uuid.UUID, status *model.EhrStatus) (err error) {
	docMeta, err := s.Infra.Index.GetDocLastByType(ctx, ehrUUID, types.Ehr)
	if err != nil {
		return fmt.Errorf("Index.GetLastEhrDocByType error: %w. ehrID: %s", err, ehrUUID.String())
	}

	ehrDecrypted, err := s.GetDocFromStorageByID(ctx, userID, docMeta.CID, ehrUUID[:], docMeta.DocUIDEncrypted)
	if err != nil {
		return fmt.Errorf("GetDocFromStorageByID error: %w. userID: %s StorageID: %x ehrID: %s", err, userID, docMeta.CID, ehrUUID.String())
	}

	var ehr model.EHR
	if err = json.Unmarshal(ehrDecrypted, &ehr); err != nil {
		return fmt.Errorf("ehr unmarshal error: %w", err)
	}

	if status.UID.Value != ehr.EhrStatus.ID.Value {
		ehr.EhrStatus.ID.Value = status.UID.Value
		if err = s.SaveEhr(ctx, userID, &ehr); err != nil {
			return fmt.Errorf("ehr save error: %w", err)
		}
	}

	return nil
}

func (s *Service) SaveStatus(ctx context.Context, userID string, ehrUUID *uuid.UUID, ehrSystemID base.EhrSystemID, status *model.EhrStatus, isNew bool) error {
	// Document encryption key generation
	key := chachaPoly.GenerateKey()

	objectVersionID, err := base.NewObjectVersionID(status.UID.Value, ehrSystemID)
	if err != nil {
		return fmt.Errorf("SaveStatus error: %w versionUID %s ehrSystemID %s", err, objectVersionID.String(), ehrSystemID.String())
	}

	baseDocumentUID := objectVersionID.BasedID()
	baseDocumentUIDHash := sha3.Sum256([]byte(baseDocumentUID))

	// Storage saving
	cidBytes, err := s.saveStatusToStorage(status, key)
	if err != nil {
		return fmt.Errorf("saveStatusToStorage error: %w", err)
	}

	// Start processing request
	requestKind := processing.RequestEhrStatusUpdate
	if isNew {
		requestKind = processing.RequestEhrStatusCreate
	}

	reqID := ctx.(*gin.Context).GetString("reqId") + "_" + string(requestKind)
	{
		procReq := &processing.Request{
			ReqID:   reqID,
			UserID:  userID,
			EhrUUID: ehrUUID.String(),
			Kind:    requestKind,
			CID:     hex.EncodeToString(cidBytes[:]),
		}
		err = s.Proc.AddRequest(procReq)
		if err != nil {
			return fmt.Errorf("Proc.AddRequest error: %w", err)
		}
	}

	// Index subject and namespace
	{
		subjectID := status.Subject.ExternalRef.ID.Value
		subjectNamespace := status.Subject.ExternalRef.Namespace

		setSubjectTx, err := s.Infra.Index.SetSubject(ctx, ehrUUID, subjectID, subjectNamespace)
		if err != nil {
			return fmt.Errorf("Index.SetSubject error: %w ehrID: %s subjectID: %s subjectNamespace: %s", err, ehrUUID.String(), subjectID, subjectNamespace)
		}

		err = s.Proc.AddBlockchainTx(reqID, setSubjectTx, "", processing.TxSetEhrBySubject, processing.StatusPending)
		if err != nil {
			return fmt.Errorf("Proc.AddBlockchainTx error: %w", err)
		}
	}

	// Index Docs ehr_id -> doc_meta
	{
		statusIDEncrypted, err := key.EncryptWithAuthData([]byte(objectVersionID.String()), ehrUUID[:])
		if err != nil {
			return fmt.Errorf("EncryptWithAuthData error: %w ehrID: %s statusUid: %s", err, ehrUUID.String(), status.UID.Value)
		}

		docMeta := &model.DocumentMeta{
			TypeCode:        types.EhrStatus,
			CID:             cidBytes,
			DocUIDEncrypted: statusIDEncrypted,
			Version:         objectVersionID.VersionTreeID(),
			DocBaseUIDHash:  &baseDocumentUIDHash,
			Timestamp:       uint64(time.Now().Unix()),
		}

		docIndexTx, err := s.Infra.Index.AddEhrDoc(ehrUUID, docMeta)
		if err != nil {
			return fmt.Errorf("Index.AddEhrDoc error: %w", err)
		}

		err = s.Proc.AddBlockchainTx(reqID, docIndexTx, "", processing.TxSetEhrDocs, processing.StatusPending)
		if err != nil {
			return fmt.Errorf("Proc.AddBlockchainTx error: %w", err)
		}
	}

	// Index Access
	{
		userPubKey, _, err := s.Infra.Keystore.Get(userID)
		if err != nil {
			return fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
		}

		docAccessValue, err := keybox.SealAnonymous(key.Bytes(), userPubKey)
		if err != nil {
			return fmt.Errorf("keybox.SealAnonymous error: %w", err)
		}

		docAccessKey := sha3.Sum256(append(cidBytes[:], []byte(userID)...))

		docAccessTx, err := s.Infra.Index.SetDocKeyEncrypted(&docAccessKey, docAccessValue)
		if err != nil {
			return fmt.Errorf("Index.SetDocAccess error: %w", err)
		}

		err = s.Proc.AddBlockchainTx(reqID, docAccessTx, "", processing.TxSetDocAccess, processing.StatusPending)
		if err != nil {
			return fmt.Errorf("Proc.AddBlockchainTx error: %w", err)
		}
	}

	if err = s.UpdateStatus(ctx, userID, ehrUUID, status); err != nil {
		return fmt.Errorf("UpdateStatus error: %w userID: %s ehrID: %s", err, userID, ehrUUID.String())
	}

	return nil
}

// GetStatus Get current (last) status of EHR document
func (s *Service) GetStatus(ctx context.Context, userID string, ehrUUID *uuid.UUID) (*model.EhrStatus, error) {
	docMeta, err := s.Infra.Index.GetDocLastByType(ctx, ehrUUID, types.EhrStatus)
	if err != nil {
		return nil, fmt.Errorf("Index.GetLastEhrDocByType error: %w. ehrID: %s", err, ehrUUID.String())
	}

	docDecrypted, err := s.GetDocFromStorageByID(ctx, userID, docMeta.CID, ehrUUID[:], docMeta.DocUIDEncrypted)
	if err != nil {
		return nil, fmt.Errorf("GetDocFromStorageByID error: %w", err)
	}

	var status model.EhrStatus
	if err := json.Unmarshal(docDecrypted, &status); err != nil {
		return nil, fmt.Errorf("EHR status unmarshal error: %w", err)
	}

	return &status, nil
}

/* Не используется
func (s *Service) GetStatusBySubject(ctx context.Context, userID, subjectID, namespace string) (*model.EhrStatus, error) {
	ehrUUID, err := s.Infra.Index.GetEhrUUIDBySubject(ctx, subjectID, namespace)
	if err != nil {
		return nil, fmt.Errorf("Index.GetEhrUUIDBySubject error: %w. userID: %s subjectID: %s namespace: %s", err, userID, subjectID, namespace)
	}

	statuses, err := s.DocsIndex.GetByType(ehrID, types.EhrStatus)
	if err != nil {
		return nil, fmt.Errorf("DocsIndex.GetByType error: %w ehrID %s docType %s", err, ehrID, types.EhrStatus.String())
	}

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		return nil, fmt.Errorf("UUID parse error: %w ehrID %s", err, ehrID)
	}

	for _, v := range statuses {
		status, err := s.getStatusFromStorage(ctx, userID, &ehrUUID, v)
		if err != nil {
			return nil, fmt.Errorf("getStatusFromStorage error: %w userID %s ehrUUID %s", err, userID, ehrID)
		}

		if status.Subject.ExternalRef.ID.Value == subjectID && status.Subject.ExternalRef.Namespace == namespace {
			return status, nil
		}
	}

	return nil, fmt.Errorf("GetStatusBySubject error: %w", errors.ErrIsNotExist)
}
*/

func (s *Service) GetStatusByNearestTime(ctx context.Context, userID string, ehrUUID *uuid.UUID, nearestTime time.Time, docType types.DocumentType) (*model.EhrStatus, error) {
	//docIndex, err := s.DocsIndex.GetByNearestTime(ehrUUID.String(), nearestTime, docType)
	docMeta, err := s.Infra.Index.GetDocByTime(ctx, ehrUUID, types.EhrStatus, uint32(nearestTime.Unix()))
	if err != nil {
		return nil, fmt.Errorf("DocsIndex.GetByNearestTime error: %w ehrID %s nearestTime %s docType %s", err, ehrUUID.String(), nearestTime.String(), docType.String())
	}

	docDecrypted, err := s.GetDocFromStorageByID(ctx, userID, docMeta.CID, ehrUUID[:], docMeta.DocUIDEncrypted)
	if err != nil {
		return nil, fmt.Errorf("GetDocFromStorageByID error: %w", err)
	}

	var status model.EhrStatus
	if err := json.Unmarshal(docDecrypted, &status); err != nil {
		return nil, fmt.Errorf("EHR status unmarshal error: %w", err)
	}

	return &status, nil
}

func (s *Service) saveStatusToStorage(status *model.EhrStatus, key *chachaPoly.Key) (*[32]byte, error) {
	statusBytes, err := json.Marshal(status)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal error: %w", err)
	}

	if s.Infra.CompressionEnabled {
		statusBytes, err = s.Infra.Compressor.Compress(statusBytes)
		if err != nil {
			return nil, fmt.Errorf("Compress error: %w", err)
		}
	}

	// Document encryption
	statusEncrypted, err := key.EncryptWithAuthData(statusBytes, []byte(status.UID.Value))
	if err != nil {
		return nil, fmt.Errorf("EncryptWithAuthData error: %w", err)
	}

	// Storage saving
	cid, err := s.Infra.IpfsClient.Add(statusEncrypted)
	if err != nil {
		return nil, fmt.Errorf("IpfsClient.Add error: %w", err)
	}

	return cid, nil
}

/*
func (s *Service) getStatusFromStorage(ctx context.Context, userID string, ehrUUID *uuid.UUID, statusMeta *model.DocumentMeta) (*model.EhrStatus, error) {
	docDecrypted, err := s.GetDocFromStorageByID(ctx, userID, statusMeta.CID, ehrUUID[:], statusMeta.DocUIDEncrypted)
	if err != nil {
		return nil, fmt.Errorf("GetDocFromStorageByID error: %w", err)
	}

	// Unmarshal EHR_STATUS
	var status model.EhrStatus
	if err := json.Unmarshal(docDecrypted, &status); err != nil {
		return nil, fmt.Errorf("EHR status unmarshal error: %w", err)
	}

	return &status, nil
}
*/

func (s *Service) ValidateEhr(ehr *model.EHR) bool {
	// TODO
	return true
}

func (s *Service) ValidateStatus(status *model.EhrStatus) bool {
	// TODO
	return true
}
