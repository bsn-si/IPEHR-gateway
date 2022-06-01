package ehr_status

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
	DocService *service.DefaultDocumentService
}

func NewEhrStatusService(docService *service.DefaultDocumentService) *EhrStatusService {
	return &EhrStatusService{
		DocService: docService,
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

func (s *EhrStatusService) Create(ehrId, ehrStatusId string) (doc *model.EhrStatus) {
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
			Value: ehrId,
		},
		Namespace: "DEMOGRAPHIC",
		Type:      "PERSON",
	}
	doc.IsQueryable = true
	doc.IsModifable = true

	return doc
}

func (s *EhrStatusService) Validate(doc *model.EhrStatus) bool {
	//TODO
	return true
}

func (s *EhrStatusService) Save(ehrId, userId string, doc *model.EhrStatus) error {
	docBytes, err := s.MarshalJson(doc)
	if err != nil {
		return err
	}

	// Document encryption key generationg
	key := chacha_poly.GenerateKey()

	// Document encryption
	docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(doc.Uid.Value))
	if err != nil {
		return err
	}

	// Storage saving
	docStorageId, err := s.DocService.Storage.Add(docEncrypted)
	if err != nil {
		return err
	}

	ehrUUID, err := uuid.Parse(ehrId)
	if err != nil {
		return err
	}
	// Doc id encryption
	docIdEncrypted, err := key.EncryptWithAuthData([]byte(doc.Uid.Value), ehrUUID[:]) //TODO should reduce doc.Uid.Value?
	if err != nil {
		return nil
	}

	// Appending EHR doc index
	docIndex := &model.DocumentMeta{
		TypeCode:       types.EHR_STATUS,
		StorageId:      docStorageId,
		DocIdEncrypted: docIdEncrypted,
		Timestamp:      uint32(time.Now().Unix()),
	}
	if err = s.DocService.AddEhrDocIndex(ehrId, docIndex); err != nil {
		return err
	}

	// Index Access
	if err = s.DocService.AccessIndex.Add(userId, docStorageId, key.Bytes()); err != nil {
		return err
	}
	return nil
}
