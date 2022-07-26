package service

import (
	"encoding/hex"
	"fmt"
	"hms/gateway/pkg/common"
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
	"hms/gateway/pkg/indexer"
	"hms/gateway/pkg/indexer/service/docAccess"
	"hms/gateway/pkg/indexer/service/docs"
	"hms/gateway/pkg/indexer/service/groupAccess"
	"hms/gateway/pkg/indexer/service/subject"
	"hms/gateway/pkg/keystore"
	"hms/gateway/pkg/storage"
)

type DefaultDocumentService struct {
	Storage  storage.Storager
	Keystore *keystore.KeyStore
	Index    *indexer.Index
	//EhrsIndex          *ehrs.Index
	DocsIndex          *docs.Index
	DocAccessIndex     *docAccess.Index
	SubjectIndex       *subject.Index
	GroupAccessIndex   *groupAccess.Index
	Compressor         compressor.Interface
	CompressionEnabled bool
}

func NewDefaultDocumentService(cfg *config.Config) *DefaultDocumentService {
	ks := keystore.New(cfg.KeystoreKey)

	return &DefaultDocumentService{
		Storage:  storage.Storage(),
		Keystore: ks,
		Index:    indexer.New(cfg.Contract.Address, cfg.Contract.Endpoint, cfg.Contract.PrivKeyPath),
		//EhrsIndex:          ehrs.New(),
		DocsIndex:          docs.New(),
		DocAccessIndex:     docAccess.New(ks),
		SubjectIndex:       subject.New(),
		GroupAccessIndex:   groupAccess.New(ks),
		Compressor:         compressor.New(cfg.CompressionLevel),
		CompressionEnabled: cfg.CompressionEnabled,
	}
}

func (d *DefaultDocumentService) GetDocIndexByObjectVersionID(userID string, ehrUUID *uuid.UUID, objectVersionID *base.ObjectVersionID) (doc *model.DocumentMeta, err error) {
	// Getting user privateKey
	userPubKey, userPrivKey, err := d.Keystore.Get(userID)
	if err != nil {
		return nil, err
	}

	docIndexes, err := d.DocsIndex.Get(ehrUUID.String())
	if err != nil {
		return nil, err
	}

	objVersionIDString := objectVersionID.String()

	for _, docIndex := range docIndexes {
		// Getting access key
		indexKey := sha3.Sum256(append(docIndex.StorageID[:], []byte(userID)...))
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

		docIDDecrypted, err := key.DecryptWithAuthData(docIndex.DocIDEncrypted, ehrUUID[:])
		if err != nil {
			continue
		}

		if objVersionIDString == string(docIDDecrypted) {
			return docIndex, nil
		}
	}

	return nil, errors.ErrIsNotExist
}

func (d *DefaultDocumentService) GetDocIndexesByBaseID(ehrUUID *uuid.UUID, objectVersionID *base.ObjectVersionID, docType types.DocumentType) ([]*model.DocumentMeta, error) {
	docIndexes, err := d.DocsIndex.Get(ehrUUID.String())
	if err != nil {
		return nil, err
	}

	var (
		docsMeta            []*model.DocumentMeta
		basedID             = objectVersionID.BasedID()
		baseDocumentUIDHash = sha3.Sum256([]byte(basedID))
	)

	for _, docIndex := range docIndexes {
		if docType > 0 && docIndex.TypeCode != docType {
			continue
		}

		if docIndex.BaseDocumentUIDHash == nil {
			continue
		}

		if *docIndex.BaseDocumentUIDHash != baseDocumentUIDHash {
			continue
		}

		docsMeta = append(docsMeta, docIndex)
	}

	return docsMeta, nil
}

func (d *DefaultDocumentService) GetDocIndexByBaseIDAndVersion(ehrUUID *uuid.UUID, objectVersionID *base.ObjectVersionID, docType types.DocumentType) (*model.DocumentMeta, error) {
	docIndexes, err := d.GetDocIndexesByBaseID(ehrUUID, objectVersionID, docType)
	if err != nil {
		return nil, err
	}

	for _, docIndex := range docIndexes {
		if docIndex.Version == objectVersionID.VersionTreeID() {
			return docIndex, nil
		}
	}

	return nil, errors.ErrIsNotExist
}

func (d *DefaultDocumentService) GetLastVersionDocIndexByBaseID(ehrUUID *uuid.UUID, objectVersionID *base.ObjectVersionID, docType types.DocumentType) (*model.DocumentMeta, error) {
	docIndexes, err := d.GetDocIndexesByBaseID(ehrUUID, objectVersionID, docType)
	if err != nil {
		return nil, fmt.Errorf("GetDocIndexesByBaseID error: %w", err)
	}

	for _, docIndex := range docIndexes {
		if docIndex.IsLastVersion {
			return docIndex, nil
		}
	}

	return nil, errors.ErrIsNotExist
}

func (d *DefaultDocumentService) GetDocFromStorageByID(userID string, storageID *[32]byte, authData []byte) (docBytes []byte, err error) {
	// Getting access key
	indexKey := sha3.Sum256(append(storageID[:], []byte(userID)...))
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

func (d *DefaultDocumentService) UpdateCollection(ehrUUID *uuid.UUID, docIndexes, toUpdate []*model.DocumentMeta, action func(*model.DocumentMeta) error) (err error) {
	changed := false

	for _, docIndex := range toUpdate {
		err := action(docIndex)
		if err != nil {
			return err
		}

		changed = true
	}

	if changed {
		if err = d.DocsIndex.Replace(ehrUUID.String(), docIndexes); err != nil {
			return err
		}
	}

	return
}

func (d *DefaultDocumentService) GenerateID() string {
	return uuid.New().String()
}

func (d *DefaultDocumentService) GetSystemID() base.EhrSystemID {
	ehrSystemID, _ := base.NewEhrSystemID(common.EhrSystemID)
	return ehrSystemID
}

func (d *DefaultDocumentService) ValidateID(id string, systemID base.EhrSystemID, docType types.DocumentType) bool {
	if docType == types.Composition {
		_, err := base.NewObjectVersionID(id, systemID)
		return err == nil
	}

	return true
}
