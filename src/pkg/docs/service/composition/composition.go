package composition

import (
	"encoding/json"
	"hms/gateway/pkg/docs/status"
	"hms/gateway/pkg/errors"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"hms/gateway/pkg/crypto/chacha_poly"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/service"
	"hms/gateway/pkg/docs/types"
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
		Status:         status.ACTIVE,
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

func (s CompositionService) delete(userId string, ehrUUID uuid.UUID, docIndex *model.DocumentMeta) (err error) {
	docIndex.Status = status.DELETED

	var docIndexes []*model.DocumentMeta
	docIndexes = append(docIndexes, docIndex)

	if err = s.Doc.DocsIndex.Replace(ehrUUID.String(), docIndexes); err != nil {
		log.Println(err)
	}

	return
}

func (c CompositionService) GetCompositionById(userId, ehrId, versionUid string, documentType types.DocumentType) (composition *model.Composition, err error) {
	documentMeta, err := c.Doc.GetDocIndexByDocId(userId, ehrId, versionUid, documentType)
	if err != nil {
		return nil, errors.IsNotExist
	}

	if documentMeta.Status == status.DELETED {
		err = errors.AlreadyDeleted
		return
	}

	decryptedData, err := c.Doc.GetDocFromStorageById(userId, documentMeta.StorageId, []byte(versionUid))
	if err != nil {
		return nil, err
	}

	return c.ParseJson(decryptedData)
}

func (c CompositionService) increaseUidVersion(uid string) string {
	base, ver := c.parseUidByVersion(uid)
	ver++

	return strings.Join(base, "::") + "::" + strconv.Itoa(ver)
}

func (c CompositionService) parseUidByVersion(uid string) (base []string, ver int) {
	base, verPart := c.parseUid(uid)

	ver = 0
	if verInt, err := strconv.Atoi(verPart); err == nil {
		ver = verInt
	}
	return
}

func (c CompositionService) parseUid(uid string) (base []string, last string) {
	re := regexp.MustCompile(`::`)
	parts := re.Split(uid, -1)
	length := len(parts) - 1
	if length == 0 {
		return parts, ""
	}

	return parts[:length], parts[length]
}

func (c CompositionService) DeleteCompositionById(userId, ehrId, versionUid string) (newUid string, err error) {
	err = c.Doc.UpdateDocStatus(userId, ehrId, versionUid, types.COMPOSITION, status.ACTIVE, status.DELETED)
	if err != nil {
		if errors.Is(err, errors.AlreadyUpdated) {
			return "", errors.AlreadyDeleted
		} else {
			return "", err
		}

	}
	newUid = c.increaseUidVersion(versionUid)

	return
}
