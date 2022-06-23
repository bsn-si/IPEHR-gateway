package service

import (
	"encoding/hex"
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"
	"hms/gateway/pkg/crypto"

	"hms/gateway/pkg/config"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer/service/doc_access"
	"hms/gateway/pkg/indexer/service/docs"
	"hms/gateway/pkg/indexer/service/ehrs"
	"hms/gateway/pkg/indexer/service/group_access"
	"hms/gateway/pkg/indexer/service/subject"
	"hms/gateway/pkg/keystore"
	"hms/gateway/pkg/storage"
)

type DefaultDocumentService struct {
	Storage          storage.Storager
	EhrsIndex        *ehrs.EhrsIndex
	DocsIndex        *docs.DocsIndex
	DocAccessIndex   *doc_access.DocAccessIndex
	SubjectIndex     *subject.SubjectIndex
	GroupAccessIndex *group_access.GroupAccessIndex
	Keystore         *keystore.KeyStore
}

func NewDefaultDocumentService(cfg *config.Config) *DefaultDocumentService {
	ks := keystore.New(cfg.KeystoreKey)
	return &DefaultDocumentService{
		EhrsIndex:        ehrs.New(),
		DocsIndex:        docs.New(),
		DocAccessIndex:   doc_access.New(ks),
		SubjectIndex:     subject.New(),
		GroupAccessIndex: group_access.New(ks),
		Storage:          storage.Storage(),
		Keystore:         ks,
	}
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

		key, err := crypto.NewKeyFromBytes(keyDecrypted)
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

func (d *DefaultDocumentService) GetDocFromStorageById(userId string, storageId *[32]byte, authData []byte) (docBytes []byte, err error) {
	userUUID, err := uuid.Parse(userId)
	if err != nil {
		return nil, err
	}

	// Getting access key
	indexKey := sha3.Sum256(append(storageId[:], userUUID[:]...))
	indexKeyStr := hex.EncodeToString(indexKey[:])
	keyEncrypted, err := d.DocAccessIndex.Get(indexKeyStr)
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

	docKey, err := crypto.NewKeyFromBytes(keyDecrypted)
	if err != nil {
		return nil, err
	}

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

func (d *DefaultDocumentService) GenerateId() string {
	return uuid.New().String()
}

func (d *DefaultDocumentService) GetSystemId() string {
	return ""
}

func (d *DefaultDocumentService) ValidateId(id string, docType types.DocumentType) bool {
	//TODO

	return true
}
