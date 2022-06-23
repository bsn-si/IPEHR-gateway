package composition

import (
	"encoding/json"
	"hms/gateway/pkg/crypto"
	"log"
	"time"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/types"
)

type CompositionService struct {
	Doc *service.DefaultDocumentService
}

func NewCompositionService(docService *service.DefaultDocumentService) *CompositionService {
	return &CompositionService{
		Doc: docService,
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

func (s CompositionService) CompositionCreate(userId, ehrId string, request *model.Composition) (composition *model.Composition, err error) {
	composition = request
	err = s.save(userId, ehrId, composition)
	return
}

func (s CompositionService) save(userId string, ehrId string, doc *model.Composition) (err error) {
	docBytes, err := s.MarshalJson(doc)
	if err != nil {
		log.Println(err)
		return
	}

	// Document encryption key generation
	key := crypto.GenerateKey()

	// Document encryption
	docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(ehrId))
	if err != nil {
		log.Println(err)
		return
	}

	// Storage saving
	docStorageId, err := s.Doc.Storage.Add(docEncrypted)
	if err != nil {
		log.Println(err)
		return
	}

	// Index Docs ehr_id -> doc_meta
	docIndex := &model.DocumentMeta{
		TypeCode:  types.COMPOSITION,
		StorageId: docStorageId,
		Timestamp: uint64(time.Now().UnixNano()),
	}

	// First record in doc index
	if err = s.Doc.DocsIndex.Add(ehrId, docIndex); err != nil {
		log.Println(err)
		return
	}

	// Index Access
	if err = s.Doc.DocAccessIndex.Add(userId, docStorageId, key.Bytes()); err != nil {
		log.Println(err)
		return
	}

	return nil
}
