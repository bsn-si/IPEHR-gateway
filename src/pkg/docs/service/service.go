package service

import (
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/crypto/chacha_poly"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer"
	"hms/gateway/pkg/keystore"
	"hms/gateway/pkg/storage"
)

type DefaultDocumentService struct {
	EhrsIndex   indexer.Indexer
	DocsIndex   indexer.Indexer
	AccessIndex indexer.Indexer
	Storage     storage.Storager
	Keystore    *keystore.KeyStore
}

func NewDefaultDocumentService() *DefaultDocumentService {
	return &DefaultDocumentService{
		EhrsIndex:   indexer.Init("ehrs"),
		DocsIndex:   indexer.Init("docs"),
		AccessIndex: indexer.Init("access"),
		Storage:     storage.Init(),
		Keystore:    keystore.New(),
	}
}

func (d *DefaultDocumentService) GetEhrDocIndexes(ehrId string) ([]*model.DocumentMeta, error) {
	var docIndexes []*model.DocumentMeta
	if err := d.DocsIndex.GetById(ehrId, &docIndexes); err != nil {
		return nil, err
	}
	return docIndexes, nil
}

func (d *DefaultDocumentService) GetLastDocIndexByType(ehrId string, docTypeCode types.DocumentType) (doc *model.DocumentMeta, err error) {
	var docIndexes []*model.DocumentMeta
	if err = d.DocsIndex.GetById(ehrId, &docIndexes); err != nil {
		return nil, err
	}

	for _, docIndex := range docIndexes {
		if docIndex.TypeCode == docTypeCode {
			if doc == nil || docIndex.Timestamp > doc.Timestamp {
				doc = docIndex
			}
		}
	}
	if doc == nil {
		return nil, errors.IsNotExist
	}
	return doc, nil
}

func (d *DefaultDocumentService) GetDocIndexByDocId(userId, ehrId, docId string, docType types.DocumentType) (doc *model.DocumentMeta, err error) {
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

	var docIndexes []*model.DocumentMeta
	if err = d.DocsIndex.GetById(ehrId, &docIndexes); err != nil {
		return nil, err
	}

	for _, docIndex := range docIndexes {
		if docType > 0 && docIndex.TypeCode != docType {
			continue
		}

		// Getting access key
		var keyEncrypted []byte
		indexKey := sha3.Sum256(append(docIndex.StorageId[:], userUUID[:]...))
		indexKeyStr := hex.EncodeToString(indexKey[:])
		err = d.AccessIndex.GetById(indexKeyStr, &keyEncrypted)
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
			return docIndex, nil
		}
	}
	return nil, errors.IsNotExist
}

func (d *DefaultDocumentService) AddEhrDocIndex(ehrId string, docIndex *model.DocumentMeta) error {
	docIndexes, err := d.GetEhrDocIndexes(ehrId)
	if err != nil {
		return err
	}
	docIndexes = append(docIndexes, docIndex)
	if err = d.DocsIndex.Replace(ehrId, docIndexes); err != nil {
		return err
	}
	return nil
}

func (d *DefaultDocumentService) GetDocFromStorageById(userId string, storageId *[32]byte, authData []byte) (docBytes []byte, err error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	// Getting access key
	var keyEncrypted []byte
	indexKey := sha3.Sum256(append(storageId[:], userUUID[:]...))
	indexKeyStr := hex.EncodeToString(indexKey[:])
	err = d.AccessIndex.GetById(indexKeyStr, &keyEncrypted)
	if err != nil {
		return nil, err
	}

	// Getting user privateKey
	userPubKey, userPrivKey, err := d.Keystore.Get(userId)
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

	var docKey chacha_poly.Key
	copy(docKey[:], keyDecrypted)

	docEncrypted, err := d.Storage.Get(storageId)
	if err != nil {
		return nil, err
	}

	// Doc decryption
	docDecrypted, err := docKey.DecryptWithAuthData(docEncrypted, authData)
	if err != nil {
		return nil, err
	}
	return docDecrypted, nil
}

func (d *DefaultDocumentService) AddAccessIndex(userId string, docStorageId *[32]byte, docKey []byte) error {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return err
	}

	// Getting user privateKey
	userPubKey, _, err := d.Keystore.Get(userId)
	if err != nil {
		return err
	}

	// Document key encryption
	keyEncrypted, err := keybox.SealAnonymous(docKey, userPubKey)
	if err != nil {
		return err
	}

	// Index doc_id -> encrypted_doc_key
	indexKey := sha3.Sum256(append(docStorageId[:], userUUID[:]...))
	indexKeyStr := hex.EncodeToString(indexKey[:])

	if err = d.AccessIndex.Add(indexKeyStr, keyEncrypted); err != nil {
		return err
	}

	return nil
}

func (d *DefaultDocumentService) GetSystemId() string {
	return ""
}

func (d *DefaultDocumentService) ValidateId(id string, docType types.DocumentType) bool {
	//TODO

	return true
}
