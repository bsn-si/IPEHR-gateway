package composition

import (
	"encoding/json"
	"fmt"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/docs/status"
	"hms/gateway/pkg/errors"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/google/uuid"

	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/model/base"
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

func (s *Service) Create(userID string, ehrUUID, groupAccessUUID *uuid.UUID, ehrSystemID base.EhrSystemID, composition *model.Composition) (*model.Composition, error) {
	groupAccess, err := s.Doc.GroupAccessIndex.Get(userID, groupAccessUUID)
	if err != nil {
		return nil, fmt.Errorf("GroupAccessIndex.Get error: %w userID %s groupAccessUUID %s", err, userID, groupAccessUUID.String())
	}

	err = s.save(userID, ehrUUID, groupAccess, ehrSystemID, composition)
	if err != nil {
		return nil, fmt.Errorf("Composition %s save error: %w", composition.UID.Value, err)
	}

	return composition, nil
}

func (s *Service) Update(userID string, ehrUUID, groupAccessUUID *uuid.UUID, ehrSystemID base.EhrSystemID, composition *model.Composition) (*model.Composition, error) {
	groupAccess, err := s.Doc.GroupAccessIndex.Get(userID, groupAccessUUID)
	if err != nil {
		return nil, fmt.Errorf("GroupAccessIndex.Get error: %w userID %s groupAccessUUID %s", err, userID, groupAccessUUID.String())
	}

	if err = s.increaseVersion(composition, ehrSystemID); err != nil {
		return nil, fmt.Errorf("Composition increaseVersion error: %w composition.UID %s", err, composition.UID.Value)
	}

	err = s.save(userID, ehrUUID, groupAccess, ehrSystemID, composition)
	if err != nil {
		return nil, fmt.Errorf("Composition save error: %w userID %s ehrUUID %s composition.UID %s", err, userID, ehrUUID.String(), composition.UID.Value)
	}

	// TODO what we should do with prev composition?
	return composition, nil
}

func (s *Service) increaseVersion(c *model.Composition, ehrSystemID base.EhrSystemID) error {
	if c == nil || c.UID == nil || c.UID.Value == "" {
		return fmt.Errorf("%w Incorrect composition UID", errors.ErrIncorrectFormat)
	}

	objectVersionID, err := base.NewObjectVersionID(c.UID.Value, ehrSystemID)
	if err != nil {
		return fmt.Errorf("increaseVersion error: %w versionUID %s ehrSystemID %s", err, objectVersionID.String(), ehrSystemID.String())
	}

	if _, err := objectVersionID.IncreaseUIDVersion(); err != nil {
		return fmt.Errorf("Composition %s IncreaseUIDVersion error: %w", c.UID.Value, err)
	}

	c.UID.Value = objectVersionID.String()

	return nil
}

func (s *Service) save(userID string, ehrUUID *uuid.UUID, groupAccess *model.GroupAccess, ehrSystemID base.EhrSystemID, doc *model.Composition) error {
	objectVersionID, err := base.NewObjectVersionID(doc.UID.Value, ehrSystemID)
	if err != nil {
		return fmt.Errorf("saving error: %w versionUID %s ehrSystemID %s", err, objectVersionID.String(), ehrSystemID.String())
	}

	baseDocumentUID := objectVersionID.BasedID()
	baseDocumentUIDHash := sha3.Sum256([]byte(baseDocumentUID))

	// Checking the existence of the Composition
	if docIndex, err := s.Doc.GetDocIndexByObjectVersionID(userID, ehrUUID, objectVersionID); err == nil {
		if docIndex != nil {
			return fmt.Errorf("GetDocIndexByObjectVersionID error: %w userID %s ehrUUID %s objectVersionID %s", err, userID, ehrUUID.String(), objectVersionID.String())
		}
	}

	docBytes, err := json.Marshal(doc)
	if err != nil {
		return fmt.Errorf("Composition marshal error: %w", err)
	}

	if s.Doc.CompressionEnabled {
		docBytes, err = s.Doc.Compressor.Compress(docBytes)
		if err != nil {
			return fmt.Errorf("Compress error: %w", err)
		}
	}

	// Document encryption key generation
	key := chachaPoly.GenerateKey()

	// Document encryption
	docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(objectVersionID.String()))
	if err != nil {
		return fmt.Errorf("EncryptWithAuthData error: %w", err)
	}

	// Storage saving
	docStorageID, err := s.Doc.Storage.Add(docEncrypted)
	if err != nil {
		return fmt.Errorf("Storage.Add error: %w", err)
	}

	docIDEncrypted, err := key.EncryptWithAuthData([]byte(objectVersionID.String()), ehrUUID[:])
	if err != nil {
		return fmt.Errorf("EncryptWithAuthData error: %w", err)
	}

	// Set previous documents last version to false
	docIndexes, err := s.Doc.DocsIndex.Get(ehrUUID.String())
	if err != nil {
		return fmt.Errorf("DocsIndex.Get error: %w ehrUUID %s", err, ehrUUID.String())
	}

	var toUpdate []*model.DocumentMeta

	for _, docIndex := range docIndexes {
		if docIndex.TypeCode != types.Composition {
			continue
		}

		if *docIndex.BaseDocumentUIDHash == baseDocumentUIDHash {
			toUpdate = append(toUpdate, docIndex)
		}
	}

	if len(toUpdate) != 0 {
		if err = s.Doc.UpdateCollection(ehrUUID, docIndexes, toUpdate, func(meta *model.DocumentMeta) error {
			meta.IsLastVersion = false
			return nil
		}); err != nil {
			return fmt.Errorf("UpdateCollection error: %w", err)
		}
	}

	// Index Docs ehr_id -> doc_meta
	docIndex := &model.DocumentMeta{
		TypeCode:            types.Composition,
		DocIDEncrypted:      docIDEncrypted,
		IsLastVersion:       true,
		Version:             objectVersionID.VersionTreeID(),
		BaseDocumentUIDHash: &baseDocumentUIDHash,
		StorageID:           docStorageID,
		Timestamp:           uint64(time.Now().UnixNano()),
		Status:              status.ACTIVE,
	}

	if err = s.Doc.DocsIndex.Add(ehrUUID.String(), docIndex); err != nil {
		return fmt.Errorf("DocsIndex.Add error: %w ehrUUID %s", err, ehrUUID.String())
	}

	docStorageIDEncrypted, err := groupAccess.Key.EncryptWithAuthData(docStorageID[:], groupAccess.GroupUUID[:])
	if err != nil {
		return fmt.Errorf("EncryptWithAuthData error: %w", err)
	}

	// Index DataSearch
	if err = s.DataSearchIndex.UpdateIndexWithNewContent(doc.Content, groupAccess, docStorageIDEncrypted); err != nil {
		return fmt.Errorf("UpdateIndexWithNewContent error: %w", err)
	}

	// Index Access
	if err = s.Doc.DocAccessIndex.Add(userID, docStorageID, key.Bytes()); err != nil {
		return fmt.Errorf("DocAccessIndex.Add error: %w userID %s", err, userID)
	}

	return nil
}

func (s *Service) GetLastByBaseID(userID string, ehrUUID *uuid.UUID, versionUID string, ehrSystemID base.EhrSystemID) (*model.Composition, error) {
	objectVersionID, err := base.NewObjectVersionID(versionUID, ehrSystemID)
	if err != nil {
		return nil, fmt.Errorf("GetLastByBaseID error: %w versionUID %s ehrSystemID %s", err, objectVersionID.String(), ehrSystemID.String())
	}

	documentMeta, err := s.Doc.GetLastVersionDocIndexByBaseID(ehrUUID, objectVersionID, types.Composition)
	if err != nil {
		return nil, fmt.Errorf("GetLastVersionDocIndexByBaseID error: %w userID %s objectVersionID %s", err, userID, objectVersionID)
	}

	if documentMeta.Status == status.DELETED {
		return nil, fmt.Errorf("GetLastByBaseID error: %w", errors.ErrAlreadyDeleted)
	}

	decryptedData, err := s.Doc.GetDocFromStorageByID(userID, documentMeta.StorageID, []byte(objectVersionID.String()))
	if err != nil {
		return nil, fmt.Errorf("GetDocFromStorageByID error: %w userID %s storageID %s", err, userID, documentMeta.StorageID)
	}

	var composition *model.Composition
	if err = json.Unmarshal(decryptedData, &composition); err != nil {
		return nil, fmt.Errorf("Composition unmarshal error: %w", err)
	}

	return composition, nil
}

func (s *Service) GetByID(userID string, ehrUUID *uuid.UUID, versionUID string, ehrSystemID base.EhrSystemID) (*model.Composition, error) {
	objectVersionID, err := base.NewObjectVersionID(versionUID, ehrSystemID)
	if err != nil {
		return nil, fmt.Errorf("NewObjectVersionID error: %w versionUID %s ehrSystemID %s", err, versionUID, ehrSystemID.String())
	}

	docIndex, err := s.Doc.GetDocIndexByBaseIDAndVersion(ehrUUID, objectVersionID, types.Composition)
	if err != nil {
		return nil, fmt.Errorf("GetDocIndexByBaseIDAndVersion error: %w ehrUUID %s objectVersionID %s", err, ehrUUID.String(), objectVersionID.String())
	}

	if docIndex.Status == status.DELETED {
		return nil, fmt.Errorf("GetCompositionByID error: %w", errors.ErrAlreadyDeleted)
	}

	decryptedData, err := s.Doc.GetDocFromStorageByID(userID, docIndex.StorageID, []byte(objectVersionID.String()))
	if err != nil {
		return nil, fmt.Errorf("GetDocFromStorageByID error: %w userID %s StorageID %x", err, userID, *docIndex.StorageID)
	}

	var composition model.Composition
	if err = json.Unmarshal(decryptedData, &composition); err != nil {
		return nil, fmt.Errorf("Composition unmarshal error: %w", err)
	}

	return &composition, nil
}

func (s *Service) DeleteByID(userID string, ehrUUID *uuid.UUID, versionUID string, ehrSystemID base.EhrSystemID) (newUID string, err error) {
	objectVersionID, err := base.NewObjectVersionID(versionUID, ehrSystemID)
	if err != nil {
		return "", fmt.Errorf("NewObjectVersionID error: %w versionUID %s ehrSystemID %s", err, versionUID, ehrSystemID.String())
	}

	var (
		docIndex            *model.DocumentMeta
		basedID             = objectVersionID.BasedID()
		baseDocumentUIDHash = sha3.Sum256([]byte(basedID))
	)

	docIndexes, err := s.Doc.DocsIndex.Get(ehrUUID.String())
	if err != nil {
		return "", fmt.Errorf("DocsIndex.Get error: %w ehrUUID %s", err, ehrUUID.String())
	}

	for _, di := range docIndexes {
		if di.TypeCode != types.Composition {
			continue
		}

		if di.BaseDocumentUIDHash == nil {
			continue
		}

		if *di.BaseDocumentUIDHash != baseDocumentUIDHash {
			continue
		}

		if di.Version == objectVersionID.VersionTreeID() {
			docIndex = di
			break
		}
	}

	if docIndex == nil {
		return "", fmt.Errorf("composition with versionUUID is not found %w", errors.ErrIsNotExist)
	}

	toUpdate := []*model.DocumentMeta{docIndex}

	err = s.Doc.UpdateCollection(ehrUUID, docIndexes, toUpdate, func(meta *model.DocumentMeta) error {
		if meta.Status == status.DELETED {
			return fmt.Errorf("UpdateCollection error: %w userID %s ehrUUID %s composition versionUID %s", errors.ErrAlreadyDeleted, userID, ehrUUID.String(), versionUID)
		}

		meta.Status = status.DELETED

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("UpdateCollection error: %w", err)
	}

	if _, err = objectVersionID.IncreaseUIDVersion(); err != nil {
		return "", fmt.Errorf("IncreaseUIDVersion error: %w objectVersionID %s", err, objectVersionID.String())
	}

	return objectVersionID.String(), nil
}
