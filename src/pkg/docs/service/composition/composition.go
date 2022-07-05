package composition

import (
	"encoding/json"
	"github.com/google/uuid"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/docs/status"
	"hms/gateway/pkg/errors"
	"log"
	"time"

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

func (s *Service) CompositionCreate(userID string, ehrUUID, groupAccessUUID *uuid.UUID, request *model.Composition) (composition *model.Composition, err error) {
	composition = request

	groupAccess, err := s.Doc.GroupAccessIndex.Get(userID, groupAccessUUID)
	if err != nil {
		log.Println("GroupAccessIndex.Get error:", err)
		return
	}

	err = s.save(userID, ehrUUID, groupAccess, composition)

	return
}

func (s *Service) CompositionUpdate(userID string, ehrUUID, groupAccessUUID *uuid.UUID, composition *model.Composition) (*model.Composition, error) {
	groupAccess, err := s.Doc.GroupAccessIndex.Get(userID, groupAccessUUID)
	if err != nil {
		return nil, err
	}

	if err = s.increaseCompositionVersion(composition); err != nil {
		return nil, err
	}

	// TODO should it be replaced with update method?
	err = s.save(userID, ehrUUID, groupAccess, composition)

	// TODO what we should do with prev composition?
	return composition, err
}

func (s *Service) increaseCompositionVersion(c *model.Composition) (err error) {
	cUid := s.GetObjectVersionIdByUid(c.UID.Value)
	if err, _ := cUid.IncreaseUidVersion(); err != nil {
		return err
	}

	c.UID.Value = cUid.String()
	return
}

func (s *Service) GetObjectVersionIdByUid(uid string) base.ObjectVersionId {
	documentUid := base.ObjectVersionId{}
	documentUid.New(uid, s.cfg.CreatingSystemId)
	return documentUid
}

func (s *Service) save(userID string, ehrUUID *uuid.UUID, groupAccess *model.GroupAccess, doc *model.Composition) (err error) {
	documentUID := doc.UID.Value

	// Checking the existence of the Composition
	if docMeta, err := s.Doc.GetDocIndexByDocID(userID, documentUID, ehrUUID, types.Composition); err == nil {
		if docMeta != nil {
			return errors.ErrAlreadyExist
		}
	}

	docBytes, err := json.Marshal(doc)
	if err != nil {
		log.Println(err)
		return
	}

	documentUid := s.GetObjectVersionIdByUid(doc.UID.Value)

	baseDocumentUid := documentUid.BasedId()

	if s.Doc.CompressionEnabled {
		docBytes, err = s.Doc.Compressor.Compress(docBytes)
		if err != nil {
			return err
		}
	}

	// Document encryption key generation
	key := chachaPoly.GenerateKey()

	// Document encryption
	docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(baseDocumentUid))
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

	docIDEncrypted, err := key.EncryptWithAuthData([]byte(baseDocumentUid), ehrUUID[:])
	if err != nil {
		return err
	}

	// Index Docs ehr_id -> doc_meta
	docIndex := &model.DocumentMeta{
		TypeCode:       types.Composition,
		DocIDEncrypted: docIDEncrypted,
		Version:        documentUid.VersionTreeId(),
		StorageID:      docStorageID,
		Timestamp:      uint64(time.Now().UnixNano()),
		Status:         status.ACTIVE,
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
	if err = s.Doc.DocAccessIndex.Add(userID, docStorageID, key.Bytes()); err != nil {
		log.Println(err)
		return
	}

	return nil
}

func (s *Service) GetLastCompositionByBaseId(userID, ehrId, versionUID string) (composition *model.Composition, err error) {
	var documentMeta *model.DocumentMeta
	documentType := types.Composition

	documentUid := s.GetObjectVersionIdByUid(versionUID)
	baseDocumentUid := documentUid.BasedId()

	documentMeta, err = s.Doc.GetLastVersionDocIndexByBaseId(userID, ehrId, baseDocumentUid, documentType)
	if err != nil {
		return nil, errors.ErrIsNotExist
	}

	if documentMeta.Status == status.DELETED {
		return nil, errors.ErrAlreadyDeleted
	}

	decryptedData, err := s.Doc.GetDocFromStorageByID(userID, documentMeta.StorageID, []byte(versionUID))
	if err != nil {
		log.Println("GroupAccessIndex.Get error:", err)
		return nil, err
	}

	err = json.Unmarshal(decryptedData, &composition)

	return
}
func (s *Service) GetCompositionByID(userID string, ehrUUID *uuid.UUID, versionUID string) (composition *model.Composition, err error) {
	documentType := types.Composition

	documentUid := s.GetObjectVersionIdByUid(versionUID)
	baseDocumentUid := documentUid.BasedId()

	documentMeta, err := s.Doc.GetDocIndexByBaseIdAndVersion(userID, ehrUUID, baseDocumentUid, documentUid.VersionTreeId(), documentType)
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

	decryptedData, err := s.Doc.GetDocFromStorageByID(userID, documentMeta.StorageID, []byte(baseDocumentUid))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(decryptedData, &composition)

	return
}

func (s *Service) DeleteCompositionByID(userID, ehrID, versionUID string) (newUID string, err error) {
	documentType := types.Composition

	documentUid := s.GetObjectVersionIdByUid(versionUID)
	baseDocumentUid := documentUid.BasedId()

	// TODO i dont like it, too much arguments
	err = s.Doc.Update(
		userID,
		ehrID,
		baseDocumentUid,
		documentUid.VersionTreeId(),
		documentType,
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

	if err, _ := documentUid.IncreaseUidVersion(); err != nil {
		return "", err
	}

	return documentUid.String(), nil
}
