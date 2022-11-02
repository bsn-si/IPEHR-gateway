package docAccess

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/service"
	proc "hms/gateway/pkg/docs/service/processing"
)

type Service struct {
	*service.DefaultDocumentService
}

func NewService(docService *service.DefaultDocumentService) *Service {
	return &Service{
		docService,
	}
}

func (s *Service) Set(ctx context.Context, userID, toUserID, reqID string, CID *cid.Cid, accessLevel uint8) error {
	_, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return fmt.Errorf("keystore.Get error: %w userID %s", err, userID)
	}

	toUserPubKey, _, err := s.Infra.Keystore.Get(toUserID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	accessID := sha3.Sum256(append(CID.Bytes()[:], []byte(toUserID)...))

	var keyEncrypted []byte
	{
		docAccessKey, err := s.GetDocAccessKey(ctx, userID, CID)
		if err != nil {
			return fmt.Errorf("Index.GetDocKeyEncrypted error: %w", err)
		}

		keyEncrypted, err = keybox.SealAnonymous(docAccessKey.Bytes(), toUserPubKey)
		if err != nil {
			return fmt.Errorf("keybox.SealAnonymous error: %w", err)
		}
	}

	data, err := s.Infra.Index.SetDocAccess(ctx, &accessID, CID.Bytes(), keyEncrypted, accessLevel, userPrivKey, nil)
	if err != nil {
		return fmt.Errorf("Index.SetGroupAccess error: %w", err)
	}

	txHash, err := s.Infra.Index.SendSingle(ctx, data)
	if err != nil {
		return fmt.Errorf("Index.SendSingle error: %w", err)
	}

	procRequest, err := s.Proc.NewRequest(reqID, userID, "", proc.RequestDocAccessSet)
	if err != nil {
		return fmt.Errorf("Proc.NewRequest error: %w", err)
	}

	procRequest.AddEthereumTx(proc.TxSetDocAccess, txHash, false)

	return nil
}
