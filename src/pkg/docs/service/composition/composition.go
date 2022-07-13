package composition

import (
	"encoding/json"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/status"
	"hms/gateway/pkg/errors"
	"log"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/google/uuid"

	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/indexer/service/dataSearch"
)

type Service struct {
	cfg             *config.Config
	Doc             *service.DefaultDocumentService
	DataSearchIndex *dataSearch.Index
}

func NewCompositionService(docService *service.DefaultDocumentService, cfg *config.Config) *Service {
	return &Service{
		Doc:             docService,
		DataSearchIndex: dataSearch.New(),
		cfg:             cfg,
	}
}

func (s *Service) Create(userUUID, ehrUUID, groupAccessUUID *uuid.UUID, request *model.Composition) (composition *model.Composition, err error) {
	composition = request

	groupAccess, err := s.Doc.GroupAccessIndex.Get(userUUID.String(), groupAccessUUID)
	if err != nil {
		log.Println("GroupAccessIndex.Get error:", err)
		return
	}

	err = s.save(userUUID, ehrUUID, groupAccess, composition)

	return
}

func (s *Service) Update(userUUID, ehrUUID, groupAccessUUID *uuid.UUID, composition *model.Composition) (*model.Composition, error) {
	groupAccess, err := s.Doc.GroupAccessIndex.Get(userUUID.String(), groupAccessUUID)
	if err != nil {
		return nil, err
	}

	if err = s.increaseCompositionVersion(composition); err != nil {
		return nil, err
	}

	// TODO should it be replaced with update method?
	err = s.save(userUUID, ehrUUID, groupAccess, composition)

	// TODO what we should do with prev composition?
	return composition, err
}

func (s *Service) increaseCompositionVersion(c *model.Composition) (err error) {
	cUID := s.Doc.GetObjectVersionIDByUID(c.UID.Value)
	if _, err := cUID.IncreaseUIDVersion(); err != nil {
		return err
	}

	c.UID.Value = cUID.String()

	return
}

func (s *Service) save(userUUID, ehrUUID *uuid.UUID, groupAccess *model.GroupAccess, doc *model.Composition) (err error) {
	objectVersionID := s.Doc.GetObjectVersionIDByUID(doc.UID.Value)
	baseDocumentUID := objectVersionID.BasedID()

	params := s.Doc.SetBaseParams(userUUID, ehrUUID, objectVersionID, types.Composition)

	// Checking the existence of the Composition
	if docMeta, err := s.Doc.GetDocIndexByObjectVersionID(params); err == nil {
		if docMeta != nil {
			return errors.ErrAlreadyExist
		}
	}

	docBytes, err := json.Marshal(doc)
	if err != nil {
		log.Println(err)
		return
	}

	if s.Doc.CompressionEnabled {
		docBytes, err = s.Doc.Compressor.Compress(docBytes)
		if err != nil {
			return err
		}
	}

	// Document encryption key generation
	key := chachaPoly.GenerateKey()

	// Document encryption
	docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(objectVersionID.String()))
	if err != nil {
		log.Println(err)
		return
	}

	// Storage saving
	docStorageID, err := s.Doc.Storage.Add(docEncrypted)
	if err != nil {
		log.Println(err)
		return
	}

	docIDEncrypted, err := key.EncryptWithAuthData([]byte(objectVersionID.String()), ehrUUID[:])
	if err != nil {
		return err
	}

	// Index Docs ehr_id -> doc_meta
	docIndex := &model.DocumentMeta{
		TypeCode:            types.Composition,
		DocIDEncrypted:      docIDEncrypted,
		IsLastVersion:       true,
		Version:             objectVersionID.VersionTreeID(),
		BaseDocumentUIDHash: sha3.Sum256([]byte(baseDocumentUID)),
		StorageID:           docStorageID,
		Timestamp:           uint64(time.Now().UnixNano()),
		Status:              status.ACTIVE,
	}

	docIndexes, err := s.Doc.GetDocIndexesByBaseID(params)
	if err != nil {
		return
	}

	//TODO need replace it to transaction model, e.g. what will happen if saving in storage failed (strange I know, but...), but docIndexes was already updated?
	if err = s.Doc.UpdateCollection(params, *docIndexes, func(meta *model.DocumentMeta) (err error) {
		meta.IsLastVersion = false
		return
	}); err != nil {
		return
	}

	if err = s.Doc.DocsIndex.Add(ehrUUID.String(), docIndex); err != nil {
		log.Println(err)
		return
	}

	docStorageIDEncrypted, err := groupAccess.Key.EncryptWithAuthData(docStorageID[:], groupAccess.GroupUUID[:])
	if err != nil {
		log.Println(err)
		return
	}

	// Index DataSearch
	if err = s.DataSearchIndex.UpdateIndexWithNewContent(doc.Content, groupAccess, docStorageIDEncrypted); err != nil {
		log.Println(err)
		return
	}

	// Index Access
	if err = s.Doc.DocAccessIndex.Add(userUUID.String(), docStorageID, key.Bytes()); err != nil {
		log.Println(err)
		return
	}

	return nil
}

func (s *Service) GetLastCompositionByBaseID(userUUID, ehrUUID *uuid.UUID, versionUID string) (composition *model.Composition, err error) {
	var documentMeta *model.DocumentMeta

	objectVersionID := s.Doc.GetObjectVersionIDByUID(versionUID)
	params := s.Doc.SetBaseParams(userUUID, ehrUUID, objectVersionID, types.Composition)

	documentMeta, err = s.Doc.GetLastVersionDocIndexByBaseID(params)
	if err != nil {
		return
	}

	if documentMeta.Status == status.DELETED {
		return nil, errors.ErrAlreadyDeleted
	}

	decryptedData, err := s.Doc.GetDocFromStorageByID(userUUID.String(), documentMeta.StorageID, []byte(objectVersionID.String()))
	if err != nil {
		log.Println("GroupAccessIndex.Get error:", err)
		return nil, err
	}

	err = json.Unmarshal(decryptedData, &composition)

	return
}

func (s *Service) GetCompositionByID(userUUID, ehrUUID *uuid.UUID, versionUID string) (composition *model.Composition, err error) {
	objectVersionID := s.Doc.GetObjectVersionIDByUID(versionUID)
	params := s.Doc.SetBaseParams(userUUID, ehrUUID, objectVersionID, types.Composition)

	documentMeta, err := s.Doc.GetDocIndexByBaseIDAndVersion(params)
	if err != nil {
		return nil, err
	}

	if documentMeta == nil {
		return nil, errors.ErrIsNotExist
	}

	if documentMeta.Status == status.DELETED {
		err = errors.ErrAlreadyDeleted

		return
	}

	decryptedData, err := s.Doc.GetDocFromStorageByID(userUUID.String(), documentMeta.StorageID, []byte(objectVersionID.String()))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(decryptedData, &composition)

	return
}

func (s *Service) DeleteCompositionByID(userUUID, ehrUUID *uuid.UUID, versionUID string) (newUID string, err error) {
	objectVersionID := s.Doc.GetObjectVersionIDByUID(versionUID)
	params := s.Doc.SetBaseParams(userUUID, ehrUUID, objectVersionID, types.Composition)

	c, err := s.Doc.GetDocIndexByBaseIDAndVersion(params)
	if err != nil {
		return
	}

	err = s.Doc.Update(
		params,
		c,
		func(meta *model.DocumentMeta) error {
			if meta.Status == status.DELETED {
				return errors.ErrAlreadyDeleted
			}

			meta.Status = status.DELETED
			return nil
		})

	if err != nil {
		return
	}

	if _, err := objectVersionID.IncreaseUIDVersion(); err != nil {
		return "", err
	}

	return objectVersionID.String(), nil
}
