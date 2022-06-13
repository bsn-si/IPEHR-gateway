package composition

import (
	"encoding/json"
	"github.com/google/uuid"
	"hms/gateway/pkg/config"
	"log"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
)

type CompositionService struct {
	Doc *service.DefaultDocumentService
	Cfg *config.Config
}

func NewCompositionService(docService *service.DefaultDocumentService, cfg *config.Config) *CompositionService {
	return &CompositionService{
		Doc: docService,
		Cfg: cfg,
	}
}

func (s CompositionService) ParseJson(data []byte) (*model.Composition, error) {
	var doc model.Composition
	err := json.Unmarshal(data, &doc)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (s CompositionService) MarshalJson(doc *model.Composition) ([]byte, error) {
	return json.Marshal(doc)
}

func (s CompositionService) Create(userId string, request *model.CompositionCreateRequest) (*model.Composition, error) {
	return s.CreateWithId(userId, uuid.New().String(), request) // TODO no new???
}

func (s CompositionService) CreateWithId(userId, ehrId string, request *model.CompositionCreateRequest) (*model.Composition, error) {
	var composition model.Composition

	//composition.SystemId.Value = s.Doc.GetSystemId() //TODO
	//composition.EhrId.Value = ehrId
	//
	//composition.EhrStatus.Id.Type = "OBJECT_VERSION_ID"
	//composition.EhrStatus.Id.Value = uuid.New().String() + "::openEHRSys.example.com::1"
	//composition.EhrStatus.Namespace = "local"
	//composition.EhrStatus.Type = "EHR_STATUS"
	//
	//composition.EhrAccess.Id.Type = "OBJECT_VERSION_ID"
	//composition.EhrAccess.Id.Value = uuid.New().String() + "::openEHRSys.example.com::1"
	//composition.EhrAccess.Namespace = "local"
	//composition.EhrAccess.Type = "EHR_ACCESS"
	//
	//composition.TimeCreated.Value = time.Now().Format("2006-01-02T15:04:05.999-07:00")

	// Creating EHR_STATUS base
	//ehrStatusService := ehr_status.NewEhrStatusService(s.Doc)
	//
	//ehrStatusId := ehr.EhrStatus.Id.Value
	//subjectId := request.Subject.ExternalRef.Id.Value
	//subjectNamespace := request.Subject.ExternalRef.Namespace
	//
	//ehrStatusDoc := ehrStatusService.Create(ehrStatusId, subjectId, subjectNamespace)

	err := s.save(userId, &composition)

	return &composition, err
}

func (s CompositionService) save(userId string, doc *model.Composition) error {
	//docBytes, err := s.MarshalJson(doc)
	_, err := s.MarshalJson(doc)
	if err != nil {
		log.Println(err)
		return err
	}

	// Document encryption key generationg
	//key := chacha_poly.GenerateKey()
	//
	//// Document encryption
	//docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(doc.EhrId.Value))
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	//
	//// Storage saving
	//docStorageId, err := s.Doc.Storage.Add(docEncrypted)
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	//
	//// Index EHR userId -> docStorageId
	//if err = s.Doc.EhrsIndex.Add(userId, docStorageId); err != nil {
	//	log.Println(err)
	//	return err
	//}
	//
	//// Index Docs ehr_id -> doc_meta
	//docIndex := &model.DocumentMeta{
	//	TypeCode:  types.EHR,
	//	StorageId: docStorageId,
	//	Timestamp: uint32(time.Now().Unix()),
	//}
	//// First record in doc index
	//if err = s.Doc.DocsIndex.Add(doc.EhrId.Value, docIndex); err != nil {
	//	log.Println(err)
	//	return err
	//}
	//
	//// Index Access
	//if err = s.Doc.AccessIndex.Add(userId, docStorageId, key.Bytes()); err != nil {
	//	log.Println(err)
	//	return err
	//}

	// Saving EHR status
	//ehrStatusService := ehr_status.NewEhrStatusService(s.Doc)
	//if err = ehrStatusService.Save(doc.EhrId.Value, userId, ehrStatusDoc); err != nil {
	//	log.Println(err)
	//	return err
	//}
	//
	//// Subject index
	//err = s.addSubjectIndex(doc.EhrId.Value, ehrStatusDoc)
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}

	return nil
}

// Add EHR status subject index
//func (s CompositionService) addSubjectIndex(ehrId string, ehrStatusDoc *model.EhrStatus) (err error) {
//	subjectId := ehrStatusDoc.Subject.ExternalRef.Id.Value
//	subjectNamespace := ehrStatusDoc.Subject.ExternalRef.Namespace
//	err = s.Doc.SubjectIndex.AddEhrSubjectsIndex(ehrId, subjectId, subjectNamespace)
//	return
//}
