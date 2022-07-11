package ehr

import (
	"encoding/json"
	"fmt"
	"golang.org/x/crypto/sha3"
	"log"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/common"
	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/types"
)

type Service struct {
	Doc *service.DefaultDocumentService
}

func NewService(docService *service.DefaultDocumentService) *Service {
	return &Service{
		Doc: docService,
	}
}

func (s *Service) EhrCreate(userID string, request *model.EhrCreateRequest) (*model.EHR, error) {
	return s.EhrCreateWithID(userID, uuid.New().String(), request)
}

func (s *Service) EhrCreateWithID(userID, ehrID string, request *model.EhrCreateRequest) (*model.EHR, error) {
	var ehr model.EHR

	ehr.SystemID.Value = s.Doc.GetSystemID()
	ehr.EhrID.Value = ehrID

	ehr.EhrStatus.ID.Type = "OBJECT_VERSION_ID"
	ehr.EhrStatus.ID.Value = uuid.New().String() + "::" + s.Doc.GetSystemID() + "::1"
	ehr.EhrStatus.Namespace = "local"
	ehr.EhrStatus.Type = "EHR_STATUS"

	ehr.EhrAccess.ID.Type = "OBJECT_VERSION_ID"
	ehr.EhrAccess.ID.Value = uuid.New().String() + "::" + s.Doc.GetSystemID() + "::1"
	ehr.EhrAccess.Namespace = "local"
	ehr.EhrAccess.Type = "EHR_ACCESS"

	ehr.TimeCreated.Value = time.Now().Format(common.OpenEhrTimeFormat)

	err := s.SaveEhr(userID, &ehr)
	if err != nil {
		return nil, fmt.Errorf("ehr save error: %w", err)
	}

	// Creating EHR_STATUS
	ehrStatusID := ehr.EhrStatus.ID.Value
	subjectID := request.Subject.ExternalRef.ID.Value
	subjectNamespace := request.Subject.ExternalRef.Namespace

	_, err = s.CreateStatus(userID, ehrID, ehrStatusID, subjectID, subjectNamespace)
	if err != nil {
		return nil, fmt.Errorf("create status error: %w", err)
	}

	return &ehr, nil
}

func (s *Service) SaveEhr(userID string, doc *model.EHR) error {
	docBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("ehr marshal error: %w", err)
	}

	if s.Doc.CompressionEnabled {
		docBytes, err = s.Doc.Compressor.Compress(docBytes)
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
	docStorageID, err := s.Doc.Storage.Add(docEncrypted)
	if err != nil {
		return fmt.Errorf("ehr storage saving error: %w", err)
	}

	// Index EHR userID -> docStorageID
	if err = s.Doc.EhrsIndex.Replace(userID, docStorageID); err != nil {
		return fmt.Errorf("ehr index replace error: %w. userID: %s docStorageID: %x", err, userID, docStorageID)
	}

	// Index Docs ehr_id -> doc_meta
	docIndex := &model.DocumentMeta{
		TypeCode:  types.Ehr,
		StorageID: docStorageID,
		Timestamp: uint64(time.Now().UnixNano()),
	}
	// First record in doc index
	if err = s.Doc.DocsIndex.Add(doc.EhrID.Value, docIndex); err != nil {
		return fmt.Errorf("docIndex add error: %w. ehrId: %s", err, doc.EhrID.Value)
	}

	// Index Access
	if err = s.Doc.DocAccessIndex.Add(userID, docStorageID, key.Bytes()); err != nil {
		return fmt.Errorf("docAccessIndex add error: %w. userID: %s docStorageID: %x", err, userID, docStorageID)
	}

	return nil
}

// GetDocBySubject Get decrypted document by subject
func (s *Service) GetDocBySubject(userID, subjectID, namespace string) (docDecrypted []byte, err error) {
	ehrID, err := s.Doc.SubjectIndex.GetEhrBySubject(subjectID, namespace)
	if err != nil {
		return nil, fmt.Errorf("SubjectIndex GetDocBySubject error: %w. userID: %s subjectID: %s namespace: %s", err, userID, subjectID, namespace)
	}

	// Getting docStorageID
	doc, err := s.Doc.DocsIndex.GetLastByType(ehrID, types.Ehr)
	if err != nil {
		return nil, fmt.Errorf("DocsIndex GetLastByType error: %w. ehrID: %s", err, ehrID)
	}

	// Getting doc from storage
	docDecrypted, err = s.Doc.GetDocFromStorageByID(userID, doc.StorageID, []byte(ehrID))
	if err != nil {
		log.Println("Can't get encrypted doc", err)
		return nil, fmt.Errorf("GetDocFromStorageByID error: %w. userID: %s, doc.StorageID: %x ehrID: %s", err, userID, doc.StorageID, ehrID)
	}

	return docDecrypted, nil
}

func (s *Service) CreateStatus(userID, ehrID, ehrStatusID, subjectID, subjectNamespace string) (doc *model.EhrStatus, err error) {
	doc = &model.EhrStatus{}
	doc.Type = types.EhrStatus.String()
	doc.ArchetypeNodeID = "openEHR-EHR-EHR_STATUS.generic.v1"
	doc.Name = base.DvText{Value: "EHR Status"}
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

	err = s.SaveStatus(ehrID, userID, doc)
	if err != nil {
		return nil, fmt.Errorf("SaveStatus error: %w. ehrID: %s userID: %s", err, ehrID, userID)
	}

	return doc, nil
}

func (s *Service) UpdateStatus(userID, ehrID string, status *model.EhrStatus) (err error) {
	docMeta, err := s.Doc.DocsIndex.GetLastByType(ehrID, types.Ehr)
	if err != nil {
		return fmt.Errorf("DocsIndex.GetLastByType error: %w. ehrID: %s", err, ehrID)
	}

	ehrDecrypted, err := s.Doc.GetDocFromStorageByID(userID, docMeta.StorageID, []byte(ehrID))
	if err != nil {
		return fmt.Errorf("GetDocFromStorageByID error: %w. userID: %s StorageID: %x ehrID: %s", err, userID, docMeta.StorageID, ehrID)
	}

	var ehr model.EHR
	if err = json.Unmarshal(ehrDecrypted, &ehr); err != nil {
		return fmt.Errorf("ehr unmarshal error: %w", err)
	}

	if status.UID.Value != ehr.EhrStatus.ID.Value {
		ehr.EhrStatus.ID.Value = status.UID.Value
		if err = s.SaveEhr(userID, &ehr); err != nil {
			return fmt.Errorf("ehr save error: %w", err)
		}
	}

	return nil
}

func (s *Service) SaveStatus(ehrID, userID string, status *model.EhrStatus) error {
	// Document encryption key generation
	key := chachaPoly.GenerateKey()

	objectVersionID := s.Doc.GetObjectVersionIDByUID(status.UID.Value)
	baseDocumentUID := objectVersionID.BasedID()

	statusStorageID, err := s.saveStatusToStorage(status, key)
	if err != nil {
		return fmt.Errorf("saveStatusToStorage error: %w", err)
	}

	subjectID := status.Subject.ExternalRef.ID.Value
	subjectNamespace := status.Subject.ExternalRef.Namespace

	err = s.Doc.SubjectIndex.AddEhrSubjectsIndex(ehrID, subjectID, subjectNamespace)
	if err != nil {
		return fmt.Errorf("SubjectIndex.AddEhrSubjectsIndex error: %w ehrID: %s subjectID: %s subjectNamespace: %s", err, ehrID, subjectID, subjectNamespace)
	}

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		return fmt.Errorf("ehr uuid parse error: %w ehrID: %s", err, ehrID)
	}

	// Doc id encryption
	statusIDEncrypted, err := key.EncryptWithAuthData([]byte(objectVersionID.String()), ehrUUID[:])
	if err != nil {
		return fmt.Errorf("EncryptWithAuthData error: %w ehrID: %s statusUid: %s", err, ehrID, status.UID.Value)
	}

	// Appending EHR doc index
	docIndex := &model.DocumentMeta{
		TypeCode:            types.EhrStatus,
		StorageID:           statusStorageID,
		DocIDEncrypted:      statusIDEncrypted,
		Version:             objectVersionID.VersionTreeID(),
		BaseDocumentUIDHash: sha3.Sum256([]byte(baseDocumentUID)),
		Timestamp:           uint64(time.Now().Unix()),
	}

	if err = s.Doc.DocsIndex.Add(ehrID, docIndex); err != nil {
		return fmt.Errorf("DocsIndex.Add error: %w ehrID: %s", err, ehrID)
	}

	// Index Access
	if err = s.Doc.DocAccessIndex.Add(userID, statusStorageID, key.Bytes()); err != nil {
		return fmt.Errorf(".DocAccessIndex.Add error: %w userID: %s statusStorageID: %x", err, userID, statusStorageID)
	}

	if err = s.UpdateStatus(userID, ehrID, status); err != nil {
		return fmt.Errorf("UpdateStatus error: %w userID: %s ehrID: %s", err, userID, ehrID)
	}

	return nil
}

// GetStatus Get current (last) status of EHR document
func (s *Service) GetStatus(userID, ehrID string) (status *model.EhrStatus, err error) {
	statusMeta, err := s.Doc.DocsIndex.GetLastByType(ehrID, types.EhrStatus)
	if err != nil {
		return
	}

	status, err = s.getStatusFromStorage(userID, ehrID, statusMeta)

	return
}

func (s *Service) GetStatusBySubject(userID, subjectID, namespace string) (status *model.EhrStatus, err error) {
	ehrID, err := s.Doc.SubjectIndex.GetEhrBySubject(subjectID, namespace)
	if err != nil {
		return
	}

	statuses, err := s.Doc.DocsIndex.GetByType(ehrID, types.EhrStatus)
	if err != nil {
		return
	}

	for _, v := range statuses {
		status, err = s.getStatusFromStorage(userID, ehrID, v)
		if err != nil {
			return
		}

		if status.Subject.ExternalRef.ID.Value == subjectID && status.Subject.ExternalRef.Namespace == namespace {
			return
		}
	}

	return
}

func (s *Service) GetStatusByNearestTime(userID, ehrID string, nearestTime time.Time, docType types.DocumentType) (status *model.EhrStatus, err error) {
	docIndex, err := s.Doc.DocsIndex.GetDocIndexByNearestTime(ehrID, nearestTime, docType)
	if err != nil {
		return
	}

	status, err = s.getStatusFromStorage(userID, ehrID, docIndex)
	if err != nil {
		return
	}

	return
}

func (s *Service) saveStatusToStorage(status *model.EhrStatus, key *chachaPoly.Key) (storageID *[32]byte, err error) {
	statusBytes, err := json.Marshal(status)
	if err != nil {
		return
	}

	if s.Doc.CompressionEnabled {
		statusBytes, err = s.Doc.Compressor.Compress(statusBytes)
		if err != nil {
			return nil, err
		}
	}

	// Document encryption
	statusEncrypted, err := key.EncryptWithAuthData(statusBytes, []byte(status.UID.Value))
	if err != nil {
		return
	}

	// Storage saving
	storageID, err = s.Doc.Storage.Add(statusEncrypted)

	return
}

func (s *Service) getStatusFromStorage(userID, ehrID string, statusMeta *model.DocumentMeta) (*model.EhrStatus, error) {
	encryptedStatus, err := s.Doc.Storage.Get(statusMeta.StorageID)
	if err != nil {
		return nil, err
	}

	statusKey, err := s.Doc.DocAccessIndex.GetDocumentKey(userID, statusMeta.StorageID)
	if err != nil {
		return nil, err
	}

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		return nil, err
	}

	statusID, err := statusKey.DecryptWithAuthData(statusMeta.DocIDEncrypted, ehrUUID[:])
	if err != nil {
		return nil, err
	}

	statusBytes, err := statusKey.DecryptWithAuthData(encryptedStatus, statusID)
	if err != nil {
		return nil, err
	}

	if s.Doc.CompressionEnabled {
		statusBytes, err = s.Doc.Compressor.Decompress(statusBytes)
		if err != nil {
			return nil, err
		}
	}

	var status model.EhrStatus
	if err = json.Unmarshal(statusBytes, &status); err != nil {
		return nil, err
	}

	return &status, nil
}

func (s *Service) ValidateEhr(ehr *model.EHR) bool {
	// TODO
	return true
}

func (s *Service) ValidateStatus(status *model.EhrStatus) bool {
	// TODO
	return true
}
