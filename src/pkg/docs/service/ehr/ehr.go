package ehr

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"

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

func (s EhrService) Create(request *model.EhrCreateRequest) (ehrDoc *model.EHR, err error) {
	ehrDoc = s.CreateWithId(uuid.New().String(), request)

	err = s.addSubjectIndex(request, ehrDoc)
	if err != nil {
		log.Println(err)
	}

	return
}

// Add EHR creation request subject index
func (s EhrService) addSubjectIndex(request *model.EhrCreateRequest, ehrDoc *model.EHR) (err error) {
	subjectId := request.Subject.ExternalRef.Id.Value
	subjectNamespace := request.Subject.ExternalRef.Namespace
	err = s.Doc.SubjectIndex.AddEhrSubjectsIndex(ehrDoc.EhrId.Value, subjectId, subjectNamespace)
	return
}

func (s EhrService) CreateWithId(ehrId string, request *model.EhrCreateRequest) *model.EHR {
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

	ehr.TimeCreated.Value = time.Now().Format("2006-01-02T15:04:05.999-07:00")

	return &ehr
}

func (s EhrService) Save(userId string, doc *model.EHR) error {
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
	if err = s.Doc.AccessIndex.Add(userId, docStorageId, key.Bytes()); err != nil {
		log.Println(err)
		return err
	}

	// Creating EHR_STATUS base
	ehrStatusService := ehr_status.NewEhrStatusService(s.Doc)
	ehrStatusDoc := ehrStatusService.Create(doc.EhrId.Value, doc.EhrStatus.Id.Value)

	if err = ehrStatusService.Save(doc.EhrId.Value, userId, ehrStatusDoc); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
