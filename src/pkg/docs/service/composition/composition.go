package composition

import (
	"encoding/json"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/status"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer/service/dataSearch"
)

type Service struct {
	Doc             *service.DefaultDocumentService
	DataSearchIndex *dataSearch.Index
}

func NewCompositionService(docService *service.DefaultDocumentService) *Service {
	return &Service{
		Doc:             docService,
		DataSearchIndex: dataSearch.New(),
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

	if s.Doc.CompressionEnabled {
		docBytes, err = s.Doc.Compressor.Compress(docBytes)
		if err != nil {
			return err
		}
	}

	// Document encryption key generation
	key := chachaPoly.GenerateKey()

	// Document encryption
	docEncrypted, err := key.EncryptWithAuthData(docBytes, []byte(documentUID))
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

	docIDEncrypted, err := key.EncryptWithAuthData([]byte(documentUID), ehrUUID[:])
	if err != nil {
		return err
	}

	// Index Docs ehr_id -> doc_meta
	docIndex := &model.DocumentMeta{
		TypeCode:       types.Composition,
		DocIDEncrypted: docIDEncrypted,
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

func (s *Service) GetCompositionByID(userID, versionUID string, ehrUUID *uuid.UUID, documentType types.DocumentType) (*model.Composition, error) {
	documentMeta, err := s.Doc.GetDocIndexByDocID(userID, versionUID, ehrUUID, documentType)
	if err != nil {
		return nil, errors.ErrIsNotExist
	}

	if documentMeta.Status == status.DELETED {
		return nil, errors.ErrAlreadyDeleted
	}

	decryptedData, err := s.Doc.GetDocFromStorageByID(userID, documentMeta.StorageID, []byte(versionUID))
	if err != nil {
		return nil, err
	}

	var composition model.Composition

	err = json.Unmarshal(decryptedData, &composition)
	if err != nil {
		return nil, err
	}

	return &composition, nil
}

func (s *Service) increaseUIDVersion(uid string) string {
	base, ver := s.parseUIDByVersion(uid)
	ver++

	return strings.Join(base, "::") + "::" + strconv.Itoa(ver)
}

func (s *Service) parseUIDByVersion(uid string) (base []string, ver int) {
	base, verPart := s.parseUID(uid)

	ver = 0
	if verInt, err := strconv.Atoi(verPart); err == nil {
		ver = verInt
	}

	return
}

func (s *Service) parseUID(uid string) (base []string, last string) {
	re := regexp.MustCompile(`::`)
	parts := re.Split(uid, -1)

	length := len(parts) - 1
	if length == 0 {
		return parts, ""
	}

	return parts[:length], parts[length]
}

func (s *Service) DeleteCompositionByID(userID, ehrID, versionUID string) (newUID string, err error) {
	err = s.Doc.UpdateDocStatus(userID, ehrID, versionUID, types.Composition, status.ACTIVE, status.DELETED)
	if err != nil {
		if errors.Is(err, errors.ErrAlreadyUpdated) {
			return "", errors.ErrAlreadyDeleted
		}

		return "", err
	}

	newUID = s.increaseUIDVersion(versionUID)

	return
}
