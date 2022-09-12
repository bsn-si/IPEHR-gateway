package service

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/ipfs/go-cid"

	"hms/gateway/pkg/common"
	"hms/gateway/pkg/config"
	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/model/base"
	"hms/gateway/pkg/docs/service/processing"
	proc "hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
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
	proc := processing.New(
		infra.LocalDB, infra.EthClient,
		infra.FilecoinClient,
		infra.IpfsClient,
		cfg.Storage.Localfile.Path,
	)
	proc.Start()

	return &DefaultDocumentService{
		Infra: infra,
		Proc:  proc,
	}
}

func (d *DefaultDocumentService) GetDocFromStorageByID(ctx context.Context, userID string, CID *cid.Cid, authData, docIDEncrypted []byte) ([]byte, error) {
	// Checking that the same request is not in processing
	{
		status, err := d.Proc.GetRetrieveStatus(CID)
		if err != nil {
			return nil, fmt.Errorf("Proc.GetRetrieveStatus error: %w CID: %s", err, CID.String())
		}

		switch status {
		case proc.StatusPending, proc.StatusProcessing:
			return nil, errors.ErrIsInProcessing
		case proc.StatusFailed:
			return nil, fmt.Errorf("%w Document retrieve failed CID: %s", errors.ErrCustom, CID.String())
		case proc.StatusSuccess, proc.StatusUnknown:
		}
	}

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
		reader, err := d.Infra.IpfsClient.Get(ctx, CID)
		if err != nil && errors.Is(err, errors.ErrNotFound) {
			// Request to recovery file from Filecoin
			if err = d.Proc.AddRetrieve(CID.String()); err != nil {
				return nil, fmt.Errorf("Proc.AddRetrieve error: %w CID %s", err, CID.String())
			}
			return nil, errors.ErrIsInProcessing
		} else if err != nil {
			return nil, fmt.Errorf("IpfsClient.Get error: %w CID %s", err, CID.String())
		}
		defer reader.Close()

		docEncrypted, err = io.ReadAll(reader)
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
