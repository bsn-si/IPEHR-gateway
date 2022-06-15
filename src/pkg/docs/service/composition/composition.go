package composition

import (
	"encoding/json"
	"github.com/google/uuid"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/crypto/chacha_poly"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/types"
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

func (s CompositionService) Create(userId string, request *model.Composition) (*model.Composition, error) {
	return s.CreateWithId(userId, uuid.New().String(), request) // TODO no new???
}

func (s CompositionService) CreateWithId(userId, ehrId string, request *model.Composition) (composition *model.Composition, err error) {
	composition = &model.Composition{}

	composition.Type = types.COMPOSITION.String()
	// TODO cant comprehend should we use data from request (because validation) or create new one?
	composition.ArchetypeNodeId = request.ArchetypeNodeId
	composition.Name.Value = request.Name.Value
	composition.Uid = &base.ObjectId{
		Type:  request.Uid.Type,
		Value: request.Uid.Type,
	}
	composition.ArchetypeDetails = &base.Archetyped{
		ArchetypeId: base.ObjectId{
			Value: request.ArchetypeDetails.ArchetypeId.Type,
		},
		TemplateId: &base.ObjectId{Value: request.ArchetypeDetails.TemplateId.Value},
		RmVersion:  request.ArchetypeDetails.RmVersion,
	}

	// TODO fill others
	//composition.Language
	//composition.Territory
	//composition.Category
	//composition.Composer
	//composition.Context
	//composition.Content

	err = s.save(userId, ehrId, composition)

	return composition, err
}

func (s CompositionService) save(userId string, ehrId string, doc *model.Composition) error {
	docBytes, err := s.MarshalJson(doc)
	if err != nil {
		log.Println(err)
		return err
	}

	// Document encryption key generationg
	key := chacha_poly.GenerateKey()

	// Document encryption
	docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(ehrId))
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
	//
	// Index EHR userId -> docStorageId
	if err = s.Doc.EhrsIndex.Add(userId, docStorageId); err != nil {
		log.Println(err)
		return err
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
		return err
	}

	// Index Access
	if err = s.Doc.DocAccessIndex.Add(userId, docStorageId, key.Bytes()); err != nil {
		log.Println(err)
		return err
	}

	return nil
}
