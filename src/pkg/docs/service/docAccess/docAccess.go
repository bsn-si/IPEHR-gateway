package docAccess

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/ipfs/go-cid"

	"hms/gateway/pkg/access"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/service"
	proc "hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer"
)

type Service struct {
	*service.DefaultDocumentService
}

func NewService(docService *service.DefaultDocumentService) *Service {
	return &Service{
		docService,
	}
}

func (s *Service) List(ctx context.Context, userID string) (access.List, error) {
	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("keystore.Get error: %w userID %s", err, userID)
	}

	acl, err := s.Infra.Index.DocAccessList(ctx, userID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("Index.DocAccessList error: %w userID: %s", err, userID)
	}

	for _, a := range acl {
		idHash, ok := a.Fields["idHash"]
		if !ok {
			return nil, fmt.Errorf("%w ACL filed 'idHash' required", errors.ErrIncorrectFormat)
		}

		idEncr, ok := a.Fields["idEncr"]
		if !ok {
			return nil, fmt.Errorf("%w ACL filed 'idEncr' required", errors.ErrIncorrectFormat)
		}

		_, ok = a.Fields["level"]
		if !ok {
			return nil, fmt.Errorf("%w ACL filed 'level' required", errors.ErrIncorrectFormat)
		}

		if len(a.Fields["level"]) == 0 {
			return nil, fmt.Errorf("%w ACL filed 'level' is empty", errors.ErrIncorrectFormat)
		}

		idDecr, err := keybox.OpenAnonymous(idEncr, userPubKey, userPrivKey)
		if err != nil {
			return nil, fmt.Errorf("keybox.Open error: %w idHash: %x", err, idHash)
		}

		idDecrHash := indexer.Keccak256(idDecr)

		if !bytes.Equal(idHash, idDecrHash) {
			return nil, fmt.Errorf("%w: mismatch idHash: %x keccak256(idDecr): %x", errors.ErrCustom, idHash, idDecrHash)
		}

		CID, err := cid.Parse(idDecr)
		if err != nil {
			return nil, fmt.Errorf("cid.Parse error: %w idDecr: %x", err, idDecr)
		}

		level := access.LevelToString(a.Fields["level"][0])

		//TODO doc description

		a.Fields = map[string][]byte{
			"CID":   []byte(CID.String()),
			"level": []byte(level),
		}
	}

	return acl, nil
}

func (s *Service) Set(ctx context.Context, userID, toUserID, reqID string, CID *cid.Cid, accessLevel uint8) error {
	_, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return fmt.Errorf("keystore.Get error: %w userID %s", err, userID)
	}

	toUserPubKey, toUserPrivKey, err := s.Infra.Keystore.Get(toUserID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	var keyEncr, CIDEncr []byte
	{
		docAccessKey, err := s.GetDocAccessKey(ctx, userID, CID)
		if err != nil {
			return fmt.Errorf("Index.GetDocKeyEncrypted error: %w", err)
		}

		keyEncr, err = keybox.SealAnonymous(docAccessKey.Bytes(), toUserPubKey)
		if err != nil {
			return fmt.Errorf("keybox.SealAnonymous error: %w", err)
		}

		CIDEncr, err = keybox.SealAnonymous(CID.Bytes(), toUserPubKey)
		if err != nil {
			return fmt.Errorf("keybox.SealAnonymous error: %w", err)
		}
	}

	data, err := s.Infra.Index.DocAccessSet(ctx, CID.Bytes(), CIDEncr, keyEncr, accessLevel, userPrivKey, toUserPrivKey, nil)
	if err != nil {
		return fmt.Errorf("Index.DocAccessSet error: %w", err)
	}

	txHash, err := s.Infra.Index.SendSingle(ctx, data)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return errors.ErrNotFound
		}

		return fmt.Errorf("Index.SendSingle error: %w", err)
	}

	procRequest, err := s.Proc.NewRequest(reqID, userID, "", proc.RequestDocAccessSet)
	if err != nil {
		return fmt.Errorf("Proc.NewRequest error: %w", err)
	}

	procRequest.AddEthereumTx(proc.TxSetDocAccess, txHash, false)

	return nil
}
