package composition

import (
	"encoding/json"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/crypto/chacha_poly"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
	"log"
	"time"
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
	key := chacha_poly.GenerateKey()

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

func (s CompositionService) GetById(userId, ehrId string, uid string) (composition *model.Composition, err error) {

	statuses, err := s.Doc.DocsIndex.GetByType(ehrId, types.COMPOSITION)
	if err != nil {
		return
	}

	for _, v := range statuses {
		composition, err := s.getCompositionFromStorage(userId, ehrId, v)
		if err != nil {
			return nil, err
		}
		if composition.Uid.Value == uid {
			return composition, nil
		}
	}

	return nil, errors.IsNotExist
}

func (s *CompositionService) getCompositionFromStorage(userId, ehrId string, documentMeta *model.DocumentMeta) (composition *model.Composition, err error) {

	encryptedData, err := s.Doc.GetDocFromStorageById(userId, documentMeta.StorageId, []byte(ehrId))
	if err == nil {
		composition, err = s.ParseJson(encryptedData)
	}

	return
}
