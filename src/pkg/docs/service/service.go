package service

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/google/uuid"
	"github.com/ipfs/go-cid"

	"hms/gateway/pkg/common"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/infrastructure"
)

type DefaultDocumentService struct {
	Infra *infrastructure.Infra
	Proc  *processing.Proc
	//EhrsIndex          *ehrs.Index
	//DocsIndex        *docs.Index
	//DocAccessIndex   *docAccess.Index
	//SubjectIndex     *subject.Index
	//GroupAccessIndex *groupAccess.Index
}

func NewDefaultDocumentService(cfg *config.Config, infra *infrastructure.Infra) *DefaultDocumentService {
	proc := processing.New(infra.LocalDB, infra.EthClient, infra.FilecoinClient)
	proc.Start()

	return &DefaultDocumentService{
		Infra: infra,
		Proc:  proc,
		//EhrsIndex:          ehrs.New(),
		//DocsIndex:        docs.New(),
		//DocAccessIndex:   docAccess.New(infra.Keystore),
		//SubjectIndex:     subject.New(),
		//GroupAccessIndex: groupAccess.New(infra.Keystore),
	}
}

/* TODO брать из блокчейна
func (d *DefaultDocumentService) GetDocIndexByObjectVersionID(userID string, ehrUUID *uuid.UUID, objectVersionID *base.ObjectVersionID) (doc *model.DocumentMeta, err error) {
	// Getting user privateKey
	userPubKey, userPrivKey, err := d.Infra.Keystore.Get(userID)
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
		indexKey := sha3.Sum256(append(docIndex.CID[:], []byte(userID)...))
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
*/

/* TODO брать из блокчейна
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
*/

/* TODO брать из блокчейна
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
*/

/* TODO брать из блокчейна
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
*/

func (d *DefaultDocumentService) GetDocFromStorageByID(ctx context.Context, userID string, CID *cid.Cid, authData, docIDEncrypted []byte) ([]byte, error) {
	// Get doc key
	var docKey *chachaPoly.Key
	{
		docKeyEncr, err := d.Infra.Index.GetDocKeyEncrypted(ctx, userID, CID)
		if err != nil {
			return nil, fmt.Errorf("Index.GetDocKeyEncrypted error: %w", err)
		}

		userPubKey, userPrivateKey, err := d.Infra.Keystore.Get(userID)
		if err != nil {
			return nil, fmt.Errorf("keystore.Get error: %w userID %s", err, userID)
		}

		docKeyBytes, err := keybox.OpenAnonymous(docKeyEncr, userPubKey, userPrivateKey)
		if err != nil {
			return nil, fmt.Errorf("keybox.OpenAnonymous error: %w", err)
		}

		docKey, err = chachaPoly.NewKeyFromBytes(docKeyBytes)
		if err != nil {
			return nil, fmt.Errorf("chachaPoly.NewKeyFromBytes error: %w", err)
		}
	}

	// Get doc encrypted
	var docEncrypted []byte
	{
		reader, err := d.Infra.IpfsClient.Get(CID)
		if err != nil {
			return nil, fmt.Errorf("IpfsClient.Get error: %w CID %s", err, CID.String())
		}
		defer reader.Close()

		docEncrypted, err = ioutil.ReadAll(reader)
		if err != nil {
			return nil, fmt.Errorf("ipfs read error: %w", err)
		}
	}

	// Decrypt and decompress
	var docDecrypted []byte
	{
		docID, err := docKey.DecryptWithAuthData(docIDEncrypted, authData)
		if err != nil {
			return nil, fmt.Errorf("DocIDEncrypted DecryptWithAuthData error: %w", err)
		}

		docDecrypted, err = docKey.DecryptWithAuthData(docEncrypted, docID)
		if err != nil {
			return nil, fmt.Errorf("docEncrypted DecryptWithAuthData error: %w", err)
		}

		if d.Infra.CompressionEnabled {
			docDecrypted, err = d.Infra.Compressor.Decompress(docDecrypted)
			if err != nil {
				return nil, fmt.Errorf("Decompress error: %w", err)
			}
		}
	}

	return docDecrypted, nil
}

/* TODO будет на блокчейне
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
*/

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
