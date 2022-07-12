package service

import (
	"encoding/hex"
	"fmt"
	"hms/gateway/pkg/compressor"
	"hms/gateway/pkg/docs/model/base"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer/service/docAccess"
	"hms/gateway/pkg/indexer/service/docs"
	"hms/gateway/pkg/indexer/service/ehrs"
	"hms/gateway/pkg/indexer/service/groupAccess"
	"hms/gateway/pkg/indexer/service/subject"
	"hms/gateway/pkg/keystore"
	"hms/gateway/pkg/storage"
)

type IndexData struct {
	userUUID        *uuid.UUID
	ehrUUID         *uuid.UUID
	objectVersionID base.ObjectVersionID
	docType         types.DocumentType
	docIndexes      *[]*model.DocumentMeta
}

type DefaultDocumentService struct {
	Storage            storage.Storager
	Keystore           *keystore.KeyStore
	EhrsIndex          *ehrs.Index
	DocsIndex          *docs.Index
	DocAccessIndex     *docAccess.Index
	SubjectIndex       *subject.Index
	GroupAccessIndex   *groupAccess.Index
	Compressor         compressor.Interface
	CompressionEnabled bool
	IndexData
}

func NewDefaultDocumentService(cfg *config.Config) *DefaultDocumentService {
	ks := keystore.New(cfg.KeystoreKey)

	return &DefaultDocumentService{
		Storage:            storage.Storage(),
		Keystore:           ks,
		EhrsIndex:          ehrs.New(),
		DocsIndex:          docs.New(),
		DocAccessIndex:     docAccess.New(ks),
		SubjectIndex:       subject.New(),
		GroupAccessIndex:   groupAccess.New(ks),
		Compressor:         compressor.New(cfg.CompressionLevel),
		CompressionEnabled: cfg.CompressionEnabled,
	}
}

func (d *DefaultDocumentService) GetObjectVersionIDByUID(UID string) base.ObjectVersionID {
	documentUID := base.ObjectVersionID{}
	documentUID.New(UID, d.GetSystemID())

	return documentUID
}

func (d *DefaultDocumentService) Init(userUUID, ehrUUID *uuid.UUID, objectVersionID base.ObjectVersionID, docType types.DocumentType) *IndexData {
	data := IndexData{
		ehrUUID:         ehrUUID,
		userUUID:        userUUID,
		objectVersionID: objectVersionID,
		docType:         docType,
	}

	return &data
}

func (d *DefaultDocumentService) SetDocIndexes(data *IndexData, indexes *[]*model.DocumentMeta) {
	data.docIndexes = indexes
}

func (d *DefaultDocumentService) GetDocIndexes(data *IndexData) (*[]*model.DocumentMeta, error) {
	if data.docIndexes == nil {
		return nil, errors.ErrIsNotExist
	}

	return data.docIndexes, nil
}

func (d *DefaultDocumentService) SaveDocIndexes(data *IndexData) (err error) {
	docIndexes, _ := d.GetDocIndexes(data)
	err = d.DocsIndex.Replace(data.ehrUUID.String(), *docIndexes)

	return
}

func (d *DefaultDocumentService) GetDocIndexByObjectVersionID(data *IndexData) (doc *model.DocumentMeta, err error) {
	// Getting user privateKey
	userPubKey, userPrivKey, err := d.Keystore.Get(data.userUUID.String())
	if err != nil {
		return nil, err
	}

	docIndexes, err := d.GetDocIndexesByBaseID(data)
	if err != nil {
		return
	}

	objectVersionID := data.objectVersionID.String()

	for _, docIndex := range *docIndexes {
		// Getting access key
		indexKey := sha3.Sum256(append(docIndex.StorageID[:], data.userUUID[:]...))
		indexKeyStr := hex.EncodeToString(indexKey[:])

		keyEncrypted, err := d.DocAccessIndex.Get(indexKeyStr)
		if err != nil {
			return nil, err
		}

		keyDecrypted, err := keybox.OpenAnonymous(keyEncrypted, userPubKey, userPrivKey)
		if err != nil {
			return nil, err
		}

		if len(keyDecrypted) != 32 {
			return nil, fmt.Errorf("%w: document key length mismatch", errors.ErrEncryption)
		}

		key, err := chachaPoly.NewKeyFromBytes(keyDecrypted)
		if err != nil {
			return nil, err
		}

		docIDDecrypted, err := key.DecryptWithAuthData(docIndex.DocIDEncrypted, data.ehrUUID[:])
		if err != nil {
			continue
		}

		if objectVersionID == string(docIDDecrypted) {
			return docIndex, nil
		}
	}

	return nil, errors.ErrIsNotExist
}

func (d *DefaultDocumentService) GetDocIndexesByEhrID(data *IndexData) (*[]*model.DocumentMeta, error) {
	docsIndexesNew, err := d.DocsIndex.Get(data.ehrUUID.String())
	if err != nil {
		return nil, err
	}

	d.SetDocIndexes(data, &docsIndexesNew)

	return &docsIndexesNew, nil
}

func (d *DefaultDocumentService) GetDocIndexesByBaseID(data *IndexData) (*[]*model.DocumentMeta, error) {
	docIndexes, err := d.GetDocIndexesByEhrID(data)
	if err != nil {
		return nil, err
	}

	var docsMeta []*model.DocumentMeta

	basedID := data.objectVersionID.BasedID()
	baseDocumentUIDHash := sha3.Sum256([]byte(basedID))

	for _, docIndex := range *docIndexes {
		if docIndex.BaseDocumentUIDHash != baseDocumentUIDHash {
			continue
		}

		if data.docType > 0 && docIndex.TypeCode != data.docType {
			continue
		}

		docsMeta = append(docsMeta, docIndex)
	}

	return &docsMeta, nil
}

func (d *DefaultDocumentService) GetLastVersionDocIndexByBaseID(data *IndexData) (*model.DocumentMeta, error) {
	documentsMeta, err := d.GetDocIndexesByBaseID(data)
	if err != nil {
		return nil, err
	}

	for _, currentDocumentMeta := range *documentsMeta {
		if currentDocumentMeta.IsLastVersion {
			return currentDocumentMeta, nil
		}
	}

	return nil, errors.ErrIsNotExist
}

func (d *DefaultDocumentService) GetDocIndexByBaseIDAndVersion(data *IndexData) (*model.DocumentMeta, error) {
	documentsMeta, err := d.GetDocIndexesByBaseID(data)
	if err != nil {
		return nil, err
	}

	for _, documentMeta := range *documentsMeta {
		// TODO OR uid===uid, but really compare strings is not good
		if data.objectVersionID.Equal(documentMeta.Version) {
			return documentMeta, nil
		}
	}

	return nil, errors.ErrIsNotExist
}

func (d *DefaultDocumentService) GetDocFromStorageByID(userID string, storageID *[32]byte, authData []byte) (docBytes []byte, err error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	// Getting access key
	indexKey := sha3.Sum256(append(storageID[:], userUUID[:]...))
	indexKeyStr := hex.EncodeToString(indexKey[:])

	keyEncrypted, err := d.DocAccessIndex.Get(indexKeyStr)
	if err != nil {
		return nil, err
	}

	// Getting user privateKey
	userPubKey, userPrivKey, err := d.Keystore.Get(userID)
	if err != nil {
		return nil, err
	}

	keyDecrypted, err := keybox.OpenAnonymous(keyEncrypted, userPubKey, userPrivKey)
	if err != nil {
		return nil, err
	}

	if len(keyDecrypted) != 32 {
		return nil, fmt.Errorf("%w: document key length mismatch", errors.ErrEncryption)
	}

	var docKey chachaPoly.Key

	copy(docKey[:], keyDecrypted)

	docEncrypted, err := d.Storage.Get(storageID)
	if err != nil {
		return nil, err
	}

	// Doc decryption
	docDecrypted, err := docKey.DecryptWithAuthData(docEncrypted, authData)
	if err != nil {
		return nil, err
	}

	if d.CompressionEnabled {
		docDecrypted, err = d.Compressor.Decompress(docDecrypted)
		if err != nil {
			return nil, err
		}
	}

	return docDecrypted, nil
}

func (d *DefaultDocumentService) UpdateCollection(data *IndexData, documentsMeta []*model.DocumentMeta, action func(*model.DocumentMeta) error) (err error) {
	changed := false

	for _, documentMeta := range documentsMeta {
		err := action(documentMeta)
		if err != nil {
			return err
		}

		changed = true
	}

	if changed {
		// TODO !!! check should be save all collection
		if err = d.SaveDocIndexes(data); err != nil {
			return err
		}
	}

	return
}

func (d *DefaultDocumentService) Update(data *IndexData, documentsMeta *model.DocumentMeta, action func(*model.DocumentMeta) error) (err error) {
	var col []*model.DocumentMeta
	col = append(col, documentsMeta)

	return d.UpdateCollection(data, col, action)
}

func (d *DefaultDocumentService) GenerateID() string {
	return uuid.New().String()
}

func (d *DefaultDocumentService) GetSystemID() string {
	// TODO how will we use it dynamically?
	return "openEHRSys.example.com"
}

func (d *DefaultDocumentService) ValidateID(id string, docType types.DocumentType) bool {
	//TODO
	return true
}
