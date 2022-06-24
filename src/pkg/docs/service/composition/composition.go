package composition

import (
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/crypto/chacha_poly"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
)

type CompositionService struct {
	Doc *service.DefaultDocumentService
}

func NewCompositionService(docService *service.DefaultDocumentService) *CompositionService {
	return &CompositionService{
		Doc: docService,
	}
}

func (s CompositionService) ParseJson(data []byte) (composition *model.Composition, err error) {
	composition = &model.Composition{}
	err = json.Unmarshal(data, composition)
	return
}

func (s CompositionService) MarshalJson(doc *model.Composition) ([]byte, error) {
	return json.Marshal(doc)
}

func (s CompositionService) CompositionCreate(userId, ehrId string, request *model.Composition) (composition *model.Composition, err error) {
	composition = request

	ehrUUID, err := uuid.Parse(ehrId)
	if err != nil {
		return
	}

	err = s.save(userId, ehrUUID, composition)
	return
}

func (s CompositionService) save(userId string, ehrUUID uuid.UUID, doc *model.Composition) (err error) {
	docBytes, err := s.MarshalJson(doc)
	if err != nil {
		log.Println(err)
		return
	}

	documentUid := doc.Uid.Value

	// Document encryption key generation
	key := chacha_poly.GenerateKey()

	// Document encryption
	docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(documentUid))
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

	docIdEncrypted, err := key.EncryptWithAuthData([]byte(documentUid), ehrUUID[:])
	if err != nil {
		return err
	}

	// Index Docs ehr_id -> doc_meta
	docIndex := &model.DocumentMeta{
		TypeCode:       types.COMPOSITION,
		DocIdEncrypted: docIdEncrypted,
		StorageId:      docStorageId,
		Timestamp:      uint64(time.Now().UnixNano()),
	}

	// First record in doc index
	if err = s.Doc.DocsIndex.Add(ehrUUID.String(), docIndex); err != nil {
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

func (c CompositionService) GetCompositionById(userId, ehrId, versionUid string, documentType types.DocumentType) (composition *model.Composition, err error) {
	documentMeta, err := c.Doc.GetDocIndexByDocId(userId, ehrId, versionUid, documentType)
	if err != nil {
		return nil, errors.IsNotExist
	}

	decryptedData, err := c.Doc.GetDocFromStorageById(userId, documentMeta.StorageId, []byte(versionUid))
	if err != nil {
		return nil, err
	}

	return c.ParseJson(decryptedData)
}
