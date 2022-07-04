package service

import (
	"encoding/hex"
	"fmt"
	"github.com/Masterminds/semver"
	"hms/gateway/pkg/compressor"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/status"
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

func (d *DefaultDocumentService) GetLastVersionDocIndexByBaseId(userId, ehrId, baseDocumentId string, documentType types.DocumentType) (documentMeta *model.DocumentMeta, err error) {

	documentsMeta, err := d.getDocIndexesByDocId(userId, ehrId, baseDocumentId, documentType)
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

func (d *DefaultDocumentService) GetDocIndexByBaseIdAndVersion(userId, ehrId, baseDocumentId, version string, documentType types.DocumentType) (documentMeta *model.DocumentMeta, err error) {
	documentsMeta, err := d.getDocIndexesByDocId(userId, ehrId, baseDocumentId, documentType)
	if err != nil {
		return nil, err
	}

	targetVersion, err := semver.NewVersion(version)
	if err != nil {
		return nil, err
	}

	for _, composition := range documentsMeta {
		v, err := semver.NewVersion(composition.Version)
		if err != nil {
			return nil, err
		}

		if v.Equal(targetVersion) {
			documentMeta = composition
			break
		}
	}

	return documentMeta, nil
}

func (d *DefaultDocumentService) getDocIndexesByDocId(userId, ehrId, docId string, docType types.DocumentType) (docs []*model.DocumentMeta, err error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	ehrUUID, err := uuid.Parse(ehrId)
	if err != nil {
		return nil, err
	}

	// Getting user privateKey
	userPubKey, userPrivKey, err := d.Keystore.Get(userId)
	if err != nil {
		return nil, err
	}

	docIndexes, err := d.DocsIndex.Get(ehrId)
	if err != nil {
		return nil, err
	}

	for _, docIndex := range docIndexes {
		if docType > 0 && docIndex.TypeCode != docType {
			continue
		}

		// Getting access key
		indexKey := sha3.Sum256(append(docIndex.StorageId[:], userUUID[:]...))
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
			return nil, fmt.Errorf("document key length mismatch")
		}

		key, err := chacha_poly.NewKeyFromBytes(keyDecrypted)
		if err != nil {
			return nil, err
		}

		docIdDecrypted, err := key.DecryptWithAuthData(docIndex.DocIdEncrypted, ehrUUID[:])
		if err != nil {
			continue
		}

		if docId == string(docIdDecrypted) {
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

func (d *DefaultDocumentService) UpdateDocStatus(userID, ehrID, docID string, docType types.DocumentType, old, new status.DocumentStatus) (err error) {
	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return err
	}

	ehrUUID, err := uuid.Parse(ehrID)
	if err != nil {
		return err
	}

	// Getting user privateKey
	userPubKey, userPrivKey, err := d.Keystore.Get(userID)
	if err != nil {
		return err
	}

	docIndexes, err := d.DocsIndex.Get(ehrID)
	if err != nil {
		return err
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
			return err
		}

		keyDecrypted, err := keybox.OpenAnonymous(keyEncrypted, userPubKey, userPrivKey)
		if err != nil {
			return err
		}

		if len(keyDecrypted) != 32 {
			return fmt.Errorf("%w: document key length mismatch", errors.ErrEncryption)
		}

		key, err := chachaPoly.NewKeyFromBytes(keyDecrypted)
		if err != nil {
			return err
		}

		docIDDecrypted, err := key.DecryptWithAuthData(docIndex.DocIDEncrypted, ehrUUID[:])
		if err != nil {
			continue
		}

		if docID == string(docIDDecrypted) {
			if docIndex.Status == new {
				return errors.ErrAlreadyUpdated
			}

			docIndex.Status = new

			if err = d.DocsIndex.Replace(ehrID, docIndexes); err != nil {
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
