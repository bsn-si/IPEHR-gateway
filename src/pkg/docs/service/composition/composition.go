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
	"hms/gateway/pkg/indexer/service/data_search"
)

type CompositionService struct {
	Doc             *service.DefaultDocumentService
	DataSearchIndex *data_search.DataSearchIndex
}

func NewCompositionService(docService *service.DefaultDocumentService) *CompositionService {
	return &CompositionService{
		Doc:             docService,
		DataSearchIndex: data_search.New(),
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

func (s CompositionService) CompositionCreate(userId string, ehrUUID *uuid.UUID, request *model.Composition) (composition *model.Composition, err error) {
	composition = request

	err = s.save(userId, ehrUUID, composition)
	return
}

func (s CompositionService) save(userId string, ehrUUID *uuid.UUID, doc *model.Composition) (err error) {
	documentUid := doc.Uid.Value

	// Checking the existence of the Composition
	if docMeta, err := s.Doc.GetDocIndexByDocId(userId, documentUid, ehrUUID, types.COMPOSITION); err == nil {
		if docMeta != nil {
			return errors.AlreadyExist
		}
	}

	// Document encryption key generation
	key := chacha_poly.GenerateKey()

	docBytes, err := s.MarshalJson(doc)
	if err != nil {
		log.Println(err)
		return
	}

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

	// Index DataSearch
	groupId := uuid.New()
	docStorId := []byte{1, 2, 3}
	if err = s.DataSearchIndex.UpdateIndexWithNewContent(doc.Content, &groupId, docStorId); err != nil {
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

func (c CompositionService) GetCompositionById(userId, versionUid string, ehrUUID *uuid.UUID, documentType types.DocumentType) (composition *model.Composition, err error) {
	documentMeta, err := c.Doc.GetDocIndexByDocId(userId, versionUid, ehrUUID, documentType)
	if err != nil {
		return nil, errors.IsNotExist
	}

	decryptedData, err := c.Doc.GetDocFromStorageById(userId, documentMeta.StorageId, []byte(versionUid))
	if err != nil {
		return nil, err
	}

	return c.ParseJson(decryptedData)
}
