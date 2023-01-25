package service

import (
	"context"
	"fmt"
	"io"

	"github.com/google/uuid"
	"github.com/ipfs/go-cid"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/common"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/compressor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/config"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model/base"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	proc "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/infrastructure"
)

type DefaultDocumentService struct {
	Infra *infrastructure.Infra
	Proc  *processing.Proc
}

func NewDefaultDocumentService(cfg *config.Config, infra *infrastructure.Infra) *DefaultDocumentService {
	proc := processing.New(
		infra.LocalDB,
		infra.EthClient,
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

func (d *DefaultDocumentService) GetDocFromStorageByID(ctx context.Context, userID, systemID string, CID *cid.Cid, authData, docIDEncrypted []byte) ([]byte, error) {
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

	// Get doc access key
	docKey, err := d.GetDocAccessKey(ctx, userID, systemID, CID)
	if err != nil {
		return nil, fmt.Errorf("GetDocAccessKey error: %w", err)
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
	var docUID, docDecrypted []byte
	{
		if authData != nil {
			docUID, err = docKey.Decrypt(docIDEncrypted)
			if err != nil {
				return nil, fmt.Errorf("DocIDEncrypted DecryptWithAuthData error: %w", err)
			}

			docDecrypted, err = docKey.DecryptWithAuthData(docEncrypted, docUID)
			if err != nil {
				return nil, fmt.Errorf("docEncrypted DecryptWithAuthData error: %w", err)
			}
		} else {
			docDecrypted, err = docKey.Decrypt(docEncrypted)
			if err != nil {
				return nil, fmt.Errorf("docEncrypted Decrypt error: %w", err)
			}
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

func (d *DefaultDocumentService) GetDocAccessKey(ctx context.Context, userID, systemID string, CID *cid.Cid) (*chachaPoly.Key, error) {
	docKeyEncr, err := d.Infra.Index.GetDocKeyEncrypted(ctx, userID, systemID, CID.Bytes())
	if err != nil {
		return nil, fmt.Errorf("Index.GetDocKeyEncrypted error: %w", err)
	}

	docKey, err := d.DecryptKey(userID, docKeyEncr)
	if err != nil {
		return nil, fmt.Errorf("decryptKey error: %w", err)
	}

	return docKey, nil
}

func (d *DefaultDocumentService) GenerateID() string {
	return uuid.New().String()
}

func (d *DefaultDocumentService) GetSystemID() string {
	ehrSystemID, _ := base.NewEhrSystemID(common.EhrSystemID)
	return ehrSystemID.String()
}

func (d *DefaultDocumentService) DecryptKey(userID string, encryptedKey []byte) (*chachaPoly.Key, error) {
	userPubKey, userPrivateKey, err := d.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("keystore.Get error: %w userID %s", err, userID)
	}

	keyBytes, err := keybox.OpenAnonymous(encryptedKey, userPubKey, userPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("keybox.OpenAnonymous error: %w", err)
	}

	key, err := chachaPoly.NewKeyFromBytes(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("chachaPoly.NewKeyFromBytes error: %w", err)
	}

	return key, nil
}

func (d *DefaultDocumentService) KeyFromAttribures(docMeta *model.DocumentMeta, userPubKey, userPrivKey *[32]byte) (*chachaPoly.Key, error) {
	keyEncr := docMeta.GetAttr(model.AttributeKeyEncr)
	if keyEncr == nil {
		return nil, errors.ErrFieldIsEmpty("AttributeKeyEncr")
	}

	keyBytes, err := keybox.OpenAnonymous(keyEncr, userPubKey, userPrivKey)
	if err != nil {
		return nil, fmt.Errorf("keybox.Open error: %w", err)
	}

	key, err := chachaPoly.NewKeyFromBytes(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("chachaPoly.NewKeyFromBytes error: %w", err)
	}

	return key, nil
}

func (d *DefaultDocumentService) ContentFromAttributes(docMeta *model.DocumentMeta, key *chachaPoly.Key) ([]byte, error) {
	var err error

	content := docMeta.GetAttr(model.AttributeContentEncr)
	if content == nil {
		return nil, errors.ErrFieldIsEmpty("AttributeContentEncr")
	}

	content, err = key.Decrypt(content)
	if err != nil {
		return nil, fmt.Errorf("Content decryption error: %w", err)
	}

	content, err = compressor.New(compressor.BestCompression).Decompress(content)
	if err != nil {
		return nil, fmt.Errorf("Content decompression error: %w", err)
	}

	return content, nil
}
