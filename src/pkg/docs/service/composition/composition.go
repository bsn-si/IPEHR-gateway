package composition

import (
	"encoding/json"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/docs/status"
	"hms/gateway/pkg/errors"
	"log"
	"time"

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
	cUID := s.GetObjectVersionIDByUID(c.UID.Value)
	if _, err := cUID.IncreaseUIDVersion(); err != nil {
		return err
	}

	c.UID.Value = cUID.String()

	return
}

func (s *Service) GetObjectVersionIDByUID(uid string) base.ObjectVersionID {
	documentUID := base.ObjectVersionID{}
	documentUID.New(uid, s.cfg.CreatingSystemID)

	return documentUID
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

	objectVersionID := s.GetObjectVersionIDByUID(doc.UID.Value)

	baseDocumentUID := objectVersionID.BasedID()

	if s.Doc.CompressionEnabled {
		docBytes, err = s.Doc.Compressor.Compress(docBytes)
		if err != nil {
			return err
		}
	}

	// Document encryption key generation
	key := chachaPoly.GenerateKey()

	// Document encryption
	docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(baseDocumentUID))
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

	docIDEncrypted, err := key.EncryptWithAuthData([]byte(baseDocumentUID), ehrUUID[:])
	if err != nil {
		return err
	}

	// Index Docs ehr_id -> doc_meta
	docIndex := &model.DocumentMeta{
		TypeCode:       types.Composition,
		DocIDEncrypted: docIDEncrypted,
		Version:        objectVersionID.VersionTreeID(),
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

func (s *Service) GetLastCompositionByBaseID(userID, ehrID, versionUID string) (composition *model.Composition, err error) {
	var documentMeta *model.DocumentMeta

	documentType := types.Composition

	documentUID := s.GetObjectVersionIDByUID(versionUID)
	baseDocumentUID := documentUID.BasedID()

	documentMeta, err = s.Doc.GetLastVersionDocIndexByBaseID(userID, ehrID, baseDocumentUID, documentType)
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

	documentUID := s.GetObjectVersionIDByUID(versionUID)
	baseDocumentUID := documentUID.BasedID()

	documentMeta, err := s.Doc.GetDocIndexByBaseIDAndVersion(userID, ehrUUID, baseDocumentUID, documentUID.VersionTreeID(), documentType)
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

	decryptedData, err := s.Doc.GetDocFromStorageByID(userID, documentMeta.StorageID, []byte(baseDocumentUID))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(decryptedData, &composition)

	return
}

func (s *Service) DeleteCompositionByID(userID, ehrID, versionUID string) (newUID string, err error) {
	documentType := types.Composition

	documentUID := s.GetObjectVersionIDByUID(versionUID)
	baseDocumentUID := documentUID.BasedID()

	// TODO i dont like it, too much arguments
	err = s.Doc.Update(
		userID,
		ehrID,
		baseDocumentUID,
		documentUID.VersionTreeID(),
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

	if _, err := documentUID.IncreaseUIDVersion(); err != nil {
		return "", err
	}

	return documentUID.String(), nil
}
