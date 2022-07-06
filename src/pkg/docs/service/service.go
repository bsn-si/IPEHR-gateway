package service

import (
	"encoding/hex"
	"fmt"
	"hms/gateway/pkg/compressor"

	"github.com/Masterminds/semver"

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

func (d *DefaultDocumentService) GetDocIndexByDocID(userID, docID string, ehrUUID *uuid.UUID, docType types.DocumentType) (doc *model.DocumentMeta, err error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	// Getting user privateKey
	userPubKey, userPrivKey, err := d.Keystore.Get(userID)
	if err != nil {
		return nil, err
	}

	docIndexes, err := d.DocsIndex.Get(ehrUUID.String())
	if err != nil {
		return nil, err
	}

	for _, docIndex := range docIndexes {
		if docType > 0 && docIndex.TypeCode != docType {
			continue
		}

		// Getting access key
		indexKey := sha3.Sum256(append(docIndex.StorageID[:], userUUID[:]...))
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

		if docID == string(docIDDecrypted) {
			return docIndex, nil
		}
	}

	return nil, errors.ErrIsNotExist
}

func (d *DefaultDocumentService) GetLastVersionDocIndexByBaseID(userID, ehrID, baseDocumentID string, documentType types.DocumentType) (documentMeta *model.DocumentMeta, err error) {
	documentsMeta, err := d.getDocIndexesByDocID(userID, ehrID, baseDocumentID, documentType)
	if err != nil {
		return nil, err
	}

	var lastVersion *semver.Version

	for _, currentDocumentMeta := range documentsMeta {
		v, err := semver.NewVersion(currentDocumentMeta.Version)
		if err != nil {
			return nil, err
		}

		if documentMeta == nil || v.GreaterThan(lastVersion) {
			documentMeta = currentDocumentMeta
			lastVersion = v
		}
	}

	return documentMeta, nil
}

func (d *DefaultDocumentService) GetDocIndexByBaseIDAndVersion(userID string, ehrUUID *uuid.UUID, baseDocumentID, version string, documentType types.DocumentType) (documentMeta *model.DocumentMeta, err error) {
	documentsMeta, err := d.getDocIndexesByDocID(userID, ehrUUID.String(), baseDocumentID, documentType)
	if err != nil {
		return nil, err
	}

	targetVersion, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}

	for _, documentMeta := range documentsMeta {
		v, err := semver.NewVersion(documentMeta.Version)
		if err != nil {
			return nil, err
		}

		if v.Equal(targetVersion) {
			return documentMeta, nil
		}
	}

	return nil, nil
}

func (d *DefaultDocumentService) getDocIndexesByDocID(userID, ehrID, docID string, docType types.DocumentType) (docs []*model.DocumentMeta, err error) {
	// TODO replace args to *uuid.UUID type
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		return nil, err
	}

	// Getting user privateKey
	userPubKey, userPrivKey, err := d.Keystore.Get(userID)
	if err != nil {
		return nil, err
	}

	docIndexes, err := d.DocsIndex.Get(ehrID)
	if err != nil {
		return nil, err
	}

	for _, docIndex := range docIndexes {
		if docType > 0 && docIndex.TypeCode != docType {
			continue
		}

		// Getting access key
		indexKey := sha3.Sum256(append(docIndex.StorageID[:], userUUID[:]...))
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
			return nil, errors.ErrKeyLengthMismatch
		}

		key, err := chachaPoly.NewKeyFromBytes(keyDecrypted)
		if err != nil {
			return nil, err
		}

		docIDDecrypted, err := key.DecryptWithAuthData(docIndex.DocIDEncrypted, ehrUUID[:])
		if err != nil {
			continue
		}

		if docID == string(docIDDecrypted) {
			docs = append(docs, docIndex)
		}
	}

	return docs, nil
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

func (d *DefaultDocumentService) Update(userID, ehrID, baseDocumentID, version string, documentType types.DocumentType, action func(*model.DocumentMeta) error) (err error) {
	documentsMeta, err := d.getDocIndexesByDocID(userID, ehrID, baseDocumentID, documentType)
	if err != nil {
		return err
	}

	targetVersion, err := semver.NewVersion(version)
	if err != nil {
		return err
	}

	for _, documentMeta := range documentsMeta {
		v, err := semver.NewVersion(documentMeta.Version)
		if err != nil {
			return err
		}

		if v.Equal(targetVersion) {
			err := action(documentMeta)
			if err != nil {
				return err
			}

			if err = d.DocsIndex.Replace(ehrID, documentsMeta); err != nil {
				return err
			}

			return nil
		}
	}

	return errors.ErrIsNotExist
}

func (d *DefaultDocumentService) GenerateID() string {
	return uuid.New().String()
}

func (d *DefaultDocumentService) GetSystemID() string {
	return ""
}

func (d *DefaultDocumentService) ValidateID(id string, docType types.DocumentType) bool {
	//TODO
	return true
}
