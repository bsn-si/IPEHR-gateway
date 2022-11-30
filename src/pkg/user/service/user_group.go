package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"

	"hms/gateway/pkg/compressor"
	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/user/model"
)

func (s *Service) GroupCreate(ctx context.Context, userID, systemID, reqID, name, description string) error {
	groupID := uuid.New()

	userGroup := &model.UserGroup{
		GroupID:     &groupID,
		Name:        name,
		Description: description,
	}

	key := chachaPoly.GenerateKey()

	idEncr, err := key.Encrypt(groupID[:])
	if err != nil {
		return fmt.Errorf("key.Encrypt groupID error: %w", err)
	}

	content, err := msgpack.Marshal(userGroup)
	if err != nil {
		return fmt.Errorf("msgpack.Marshal error: %w", err)
	}

	contentCompresed, err := compressor.New(compressor.BestCompression).Compress(content)
	if err != nil {
		return fmt.Errorf("Query Compress error: %w", err)
	}

	contentEncr, err := key.Encrypt(contentCompresed)
	if err != nil {
		return fmt.Errorf("key.Encrypt content error: %w", err)
	}

	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	keyEncr, err := keybox.Seal(key.Bytes(), userPubKey, userPrivKey)
	if err != nil {
		return fmt.Errorf("keybox.SealAnonymous error: %w", err)
	}

	procRequest, err := s.Proc.NewRequest(reqID, userID, "", processing.RequestUserGroupCreate)
	if err != nil {
		return fmt.Errorf("Proc.NewRequest error: %w", err)
	}

	txHash, err := s.Infra.Index.GroupCreate(ctx, &groupID, idEncr, keyEncr, contentEncr, userPrivKey, nil)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return err
		} else if errors.Is(err, errors.ErrAlreadyExist) {
			return err
		}

		return fmt.Errorf("Index.GroupCreate error: %w", err)
	}

	procRequest.AddEthereumTx(processing.TxUserGroupCreate, txHash)

	if err := procRequest.Commit(); err != nil {
		return fmt.Errorf("EHR create procRequest commit error: %w", err)
	}

	return nil
}
