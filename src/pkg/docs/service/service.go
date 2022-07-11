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

type indexData struct {
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
	indexData
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

func (d *DefaultDocumentService) Init(userUUID, ehrUUID *uuid.UUID, objectVersionID base.ObjectVersionID, docType types.DocumentType) {
	d.ehrUUID = ehrUUID
	d.userUUID = userUUID
	d.objectVersionID = objectVersionID
	d.docType = docType
}

func (d *DefaultDocumentService) SetDocIndexes(indexes *[]*model.DocumentMeta) {
	d.docIndexes = indexes
}

func (d *DefaultDocumentService) GetDocIndexes() (*[]*model.DocumentMeta, error) {
	if d.docIndexes == nil {
		return nil, errors.ErrIsNotExist
	}

	return d.docIndexes, nil
}

func (d *DefaultDocumentService) SaveDocIndexes() (err error) {
	docIndexes, _ := d.GetDocIndexes()
	err = d.DocsIndex.Replace(d.ehrUUID.String(), *docIndexes)

	return
}

func (d *DefaultDocumentService) GetDocIndexByObjectVersionID() (doc *model.DocumentMeta, err error) {
	// Getting user privateKey
	userPubKey, userPrivKey, err := d.Keystore.Get(d.userUUID.String())
	if err != nil {
		return nil, err
	}

	docIndexes, err := d.GetDocIndexesByBaseID()
	if err != nil {
		return
	}

	objectVersionID := d.objectVersionID.String()

	for _, docIndex := range *docIndexes {
		// Getting access key
		indexKey := sha3.Sum256(append(docIndex.StorageID[:], d.userUUID[:]...))
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

		docIDDecrypted, err := key.DecryptWithAuthData(docIndex.DocIDEncrypted, d.ehrUUID[:])
		if err != nil {
			continue
		}

		if objectVersionID == string(docIDDecrypted) {
			return docIndex, nil
		}
	}

	return nil, errors.ErrIsNotExist
}

func (d *DefaultDocumentService) GetDocIndexesByEhrID() (*[]*model.DocumentMeta, error) {
	docsIndexesNew, err := d.DocsIndex.Get(d.ehrUUID.String())
	if err != nil {
		return nil, err
	}
	d.SetDocIndexes(&docsIndexesNew)

	return &docsIndexesNew, nil
}

func (d *DefaultDocumentService) GetDocIndexesByBaseID() (*[]*model.DocumentMeta, error) {
	docIndexes, err := d.GetDocIndexesByEhrID()
	if err != nil {
		return nil, err
	}

	var docsMeta []*model.DocumentMeta

	basedID := d.objectVersionID.BasedID()
	baseDocumentUIDHash := sha3.Sum256([]byte(basedID))

	for _, docIndex := range *docIndexes {
		if docIndex.BaseDocumentUIDHash != baseDocumentUIDHash {
			continue
		}

		if d.docType > 0 && docIndex.TypeCode != d.docType {
			continue
		}

		docsMeta = append(docsMeta, docIndex)
	}

	return &docsMeta, nil
}

func (d *DefaultDocumentService) GetLastVersionDocIndexByBaseID() (*model.DocumentMeta, error) {
	documentsMeta, err := d.GetDocIndexesByBaseID()
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

func (d *DefaultDocumentService) GetDocIndexByBaseIDAndVersion() (*model.DocumentMeta, error) {
	documentsMeta, err := d.GetDocIndexesByBaseID()
	if err != nil {
		return nil, err
	}

	for _, documentMeta := range *documentsMeta {
		if d.objectVersionID.Equal(documentMeta.Version) {
			return documentMeta, nil
		}
		// TODO OR uid===uid, but really compare strings is not good
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

func (d *DefaultDocumentService) UpdateCollection(documentsMeta []*model.DocumentMeta, action func(*model.DocumentMeta) error) (err error) {
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
		if err = d.SaveDocIndexes(); err != nil {
			return err
		}
	}

	return
}

func (d *DefaultDocumentService) Update(documentsMeta *model.DocumentMeta, action func(*model.DocumentMeta) error) (err error) {
	var col []*model.DocumentMeta
	col = append(col, documentsMeta)

	return d.UpdateCollection(col, action)
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
