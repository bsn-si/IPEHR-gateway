package ehr

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/crypto/chacha_poly"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/types"
)

type EhrStatusService struct {
	Doc *service.DefaultDocumentService
}

func NewEhrStatusService(docService *service.DefaultDocumentService) *EhrStatusService {
	return &EhrStatusService{
		Doc: docService,
	}
}

func (s *EhrStatusService) ParseJson(data []byte) (*model.EhrStatus, error) {
	var doc model.EhrStatus
	err := json.Unmarshal(data, &doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (s *EhrStatusService) MarshalJson(doc *model.EhrStatus) ([]byte, error) {
	return json.Marshal(doc)
}

func (s *EhrStatusService) Create(userId, ehrId, ehrStatusId, subjectId, subjectNamespace string) (doc *model.EhrStatus, err error) {
	doc = &model.EhrStatus{}
	doc.Type = types.EHR_STATUS.String()
	doc.ArchetypeNodeId = "openEHR-EHR-EHR_STATUS.generic.v1"
	doc.Name = base.DvText{Value: "EHR Status"}
	doc.Uid = &base.ObjectId{
		Type:  "OBJECT_VERSION_ID",
		Value: ehrStatusId,
	}
	doc.Subject.ExternalRef = base.ObjectRef{
		Id: base.ObjectId{
			Type:  "HIER_OBJECT_ID",
			Value: subjectId,
		},
		Namespace: subjectNamespace,
		Type:      "PERSON",
	}
	doc.IsQueryable = true
	doc.IsModifable = true

	err = s.Save(ehrId, userId, doc)
	if err != nil {
		return
	}

	return
}

func (s *EhrStatusService) Validate(doc *model.EhrStatus) bool {
	//TODO
	return true
}

func (s *EhrStatusService) Save(ehrId, userId string, status *model.EhrStatus) error {
	// Document encryption key generation
	key := chacha_poly.GenerateKey()

	statusStorageId, err := s.saveStatusToStorage(status, key, ehrId)
	if err != nil {
		return err
	}

	subjectId := status.Subject.ExternalRef.Id.Value
	subjectNamespace := status.Subject.ExternalRef.Namespace
	err = s.Doc.SubjectIndex.AddEhrSubjectsIndex(ehrId, subjectId, subjectNamespace)
	if err != nil {
		return err
	}

	ehrUUID, err := uuid.Parse(ehrId)
	if err != nil {
		return err
	}
	// Doc id encryption
	statusIdEncrypted, err := key.EncryptWithAuthData([]byte(status.Uid.Value), ehrUUID[:]) //TODO should reduce doc.Uid.Value?
	if err != nil {
		return err
	}

	// Appending EHR doc index
	docIndex := &model.DocumentMeta{
		TypeCode:       types.EHR_STATUS,
		StorageId:      statusStorageId,
		DocIdEncrypted: statusIdEncrypted,
		Timestamp:      uint64(time.Now().UnixNano()),
	}
	if err = s.Doc.DocsIndex.Add(ehrId, docIndex); err != nil {
		return err
	}

	// Index Access
	if err = s.Doc.DocAccessIndex.Add(userId, statusStorageId, key.Bytes()); err != nil {
		return err
	}

	ehrService := NewEhrService(s.Doc)
	err = ehrService.UpdateDocumentStatus(userId, ehrId, status)
	if err != nil {
		return err
	}

	return nil
}

func (s *EhrStatusService) saveStatusToStorage(status *model.EhrStatus, key *chacha_poly.Key, ehrId string) (storageId *[32]byte, err error) {
	statusBytes, err := s.MarshalJson(status)
	if err != nil {
		return
	}

	// Document encryption
	statusEncrypted, err := key.EncryptWithAuthData(statusBytes, []byte(status.Uid.Value))
	if err != nil {
		return
	}

	// Storage saving
	storageId, err = s.Doc.Storage.Add(statusEncrypted)
	return
}

// Get current (last) status of EHR document
func (s *EhrStatusService) Get(userId, ehrId string) (status *model.EhrStatus, err error) {
	statusMeta, err := s.Doc.DocsIndex.GetLastByType(ehrId, types.EHR_STATUS)
	if err != nil {
		return
	}

	status, err = s.getStatusFromStorage(userId, ehrId, statusMeta)

	return
}

func (s *EhrStatusService) GetStatusBySubject(userId, subjectId, namespace string) (status *model.EhrStatus, err error) {
	ehrId, err := s.Doc.SubjectIndex.GetEhrBySubject(subjectId, namespace)
	if err != nil {
		return
	}

	statuses, err := s.Doc.DocsIndex.GetByType(ehrId, types.EHR_STATUS)
	if err != nil {
		return
	}

	for _, v := range statuses {
		status, err = s.getStatusFromStorage(userId, ehrId, v)
		if err != nil {
			return
		}
		if status.Subject.ExternalRef.Id.Value == subjectId && status.Subject.ExternalRef.Namespace == namespace {
			return
		}
	}
	return
}

func (s *EhrStatusService) getStatusFromStorage(userId, ehrId string, statusMeta *model.DocumentMeta) (status *model.EhrStatus, err error) {
	statusKeyBytes, err := s.Doc.DocAccessIndex.GetDocumentKey(userId, statusMeta.StorageId)
	if err != nil {
		return
	}

	statusKey, err := chacha_poly.NewKeyFromBytes(statusKeyBytes)
	if err != nil {
		return
	}

	encryptedStatus, err := s.Doc.Storage.Get(statusMeta.StorageId)
	if err != nil {
		return
	}

	ehrUUID, err := uuid.Parse(ehrId)
	statusId, err := statusKey.DecryptWithAuthData(statusMeta.DocIdEncrypted, ehrUUID[:])
	if err != nil {
		return
	}

	statusBytes, err := statusKey.DecryptWithAuthData(encryptedStatus, statusId)
	if err != nil {
		return
	}

	status, err = s.ParseJson(statusBytes)

	return
}
