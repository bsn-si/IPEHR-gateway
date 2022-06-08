package ehr

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/common"
	"hms/gateway/pkg/crypto/chacha_poly"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/service/ehr_status"
	"hms/gateway/pkg/docs/types"
)

type EhrService struct {
	Doc *service.DefaultDocumentService
}

func NewEhrService(docService *service.DefaultDocumentService) *EhrService {
	return &EhrService{
		Doc: docService,
	}
}

func (s EhrService) ParseJson(data []byte) (*model.EHR, error) {
	var doc model.EHR
	err := json.Unmarshal(data, &doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (s EhrService) MarshalJson(doc *model.EHR) ([]byte, error) {
	return json.Marshal(doc)
}

func (s EhrService) Create(userId string, request *model.EhrCreateRequest) (*model.EHR, error) {
	return s.CreateWithId(userId, uuid.New().String(), request)
}

func (s EhrService) CreateWithId(userId, ehrId string, request *model.EhrCreateRequest) (*model.EHR, error) {
	var ehr model.EHR

	ehr.SystemId.Value = s.Doc.GetSystemId() //TODO
	ehr.EhrId.Value = ehrId

	ehr.EhrStatus.Id.Type = "OBJECT_VERSION_ID"
	ehr.EhrStatus.Id.Value = uuid.New().String() + "::openEHRSys.example.com::1"
	ehr.EhrStatus.Namespace = "local"
	ehr.EhrStatus.Type = "EHR_STATUS"

	ehr.EhrAccess.Id.Type = "OBJECT_VERSION_ID"
	ehr.EhrAccess.Id.Value = uuid.New().String() + "::openEHRSys.example.com::1"
	ehr.EhrAccess.Namespace = "local"
	ehr.EhrAccess.Type = "EHR_ACCESS"

	ehr.TimeCreated.Value = time.Now().Format(common.OPENEHR_TIME_FORMAT)

	// Creating EHR_STATUS base
	ehrStatusService := ehr_status.NewEhrStatusService(s.Doc)

	ehrStatusId := ehr.EhrStatus.Id.Value
	subjectId := request.Subject.ExternalRef.Id.Value
	subjectNamespace := request.Subject.ExternalRef.Namespace

	ehrStatusDoc := ehrStatusService.Create(ehrStatusId, subjectId, subjectNamespace)

	err := s.save(userId, &ehr, ehrStatusDoc)

	return &ehr, err
}

func (s EhrService) save(userId string, doc *model.EHR, ehrStatusDoc *model.EhrStatus) error {
	docBytes, err := s.MarshalJson(doc)
	if err != nil {
		log.Println(err)
		return err
	}

	// Document encryption key generationg
	key := chacha_poly.GenerateKey()

	// Document encryption
	docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(doc.EhrId.Value))
	if err != nil {
		log.Println(err)
		return err
	}

	// Storage saving
	docStorageId, err := s.Doc.Storage.Add(docEncrypted)
	if err != nil {
		log.Println(err)
		return err
	}

	// Index EHR userId -> docStorageId
	if err = s.Doc.EhrsIndex.Add(userId, docStorageId); err != nil {
		log.Println(err)
		return err
	}

	// Index Docs ehr_id -> doc_meta
	docIndex := &model.DocumentMeta{
		TypeCode:  types.EHR,
		StorageId: docStorageId,
		Timestamp: uint32(time.Now().Unix()),
	}
	// First record in doc index
	if err = s.Doc.DocsIndex.Add(doc.EhrId.Value, docIndex); err != nil {
		log.Println(err)
		return err
	}

	// Index Access
	if err = s.Doc.DocAccessIndex.Add(userId, docStorageId, key.Bytes()); err != nil {
		log.Println(err)
		return err
	}

	// Saving EHR status
	ehrStatusService := ehr_status.NewEhrStatusService(s.Doc)
	if err = ehrStatusService.Save(doc.EhrId.Value, userId, ehrStatusDoc); err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// GetDocBySubject Get decrypted document by subject
func (s EhrService) GetDocBySubject(userId, subjectId, namespace string) (docDecrypted []byte, err error) {
	ehrId, err := s.Doc.SubjectIndex.GetEhrBySubject(subjectId, namespace)
	if err != nil {
		log.Println("Can't get ehrId", "subjectId", subjectId, err)
		return
	}

	// Getting docStorageId
	doc, err := s.Doc.DocsIndex.GetLastByType(ehrId, types.EHR)
	if err != nil {
		log.Println("Can't get docStorageId by ehrId", "ehrId", ehrId, err)
		return
	}

	// Getting doc from storage
	docDecrypted, err = s.Doc.GetDocFromStorageById(userId, doc.StorageId, []byte(ehrId))
	if err != nil {
		log.Println("Can't get encrypted doc", err)
	}
	return
}
