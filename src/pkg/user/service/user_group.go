package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"

	"hms/gateway/pkg/access"
	"hms/gateway/pkg/compressor"
	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/user/model"
)

func (s *Service) GroupCreate(ctx context.Context, userID, systemID, reqID, name, description string) (*uuid.UUID, error) {
	groupID := uuid.New()

	userGroup := &model.UserGroup{
		GroupID:     &groupID,
		Name:        name,
		Description: description,
	}

	key := chachaPoly.GenerateKey()

	idEncr, err := key.Encrypt(groupID[:])
	if err != nil {
		return nil, fmt.Errorf("key.Encrypt groupID error: %w", err)
	}

	content, err := msgpack.Marshal(userGroup)
	if err != nil {
		return nil, fmt.Errorf("msgpack.Marshal error: %w", err)
	}

	contentCompresed, err := compressor.New(compressor.BestCompression).Compress(content)
	if err != nil {
		return nil, fmt.Errorf("UserGroup content compression error: %w", err)
	}

	contentEncr, err := key.Encrypt(contentCompresed)
	if err != nil {
		return nil, fmt.Errorf("key.Encrypt content error: %w", err)
	}

	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	keyEncr, err := keybox.Seal(key.Bytes(), userPubKey, userPrivKey)
	if err != nil {
		return nil, fmt.Errorf("keybox.SealAnonymous error: %w", err)
	}

	procRequest, err := s.Proc.NewRequest(reqID, userID, "", processing.RequestUserGroupCreate)
	if err != nil {
		return nil, fmt.Errorf("Proc.NewRequest error: %w", err)
	}

	txHash, err := s.Infra.Index.UserGroupCreate(ctx, &groupID, idEncr, keyEncr, contentEncr, userPrivKey, nil)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		} else if errors.Is(err, errors.ErrAlreadyExist) {
			return nil, err
		}

		return nil, fmt.Errorf("Index.GroupCreate error: %w", err)
	}

	procRequest.AddEthereumTx(processing.TxUserGroupCreate, txHash)

	if err := procRequest.Commit(); err != nil {
		return nil, fmt.Errorf("EHR create procRequest commit error: %w", err)
	}

	return userGroup.GroupID, nil
}

func (s *Service) GroupGetByID(ctx context.Context, userID string, groupID *uuid.UUID) (*model.UserGroup, error) {
	key, err := s.getAccessKey(ctx, userID, access.UserGroup, groupID[:])
	if err != nil {
		if errors.Is(err, errors.ErrAccessDenied) {
			return nil, err
		}

		return nil, fmt.Errorf("getAccessKey error: %w", err)
	}

	userGroup, err := s.Infra.Index.UserGroupGetByID(ctx, userID, groupID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("Index.UserGroupGetByID error: %w", err)
	}

	contentCompresed, err := key.Decrypt(userGroup.ContentEncr)
	if err != nil {
		return nil, fmt.Errorf("UserGroup Content decrypt error: %w", err)
	}

	content, err := compressor.New(compressor.BestCompression).Decompress(contentCompresed)
	if err != nil {
		return nil, fmt.Errorf("UserGroup content decompression error: %w", err)
	}

	err = msgpack.Unmarshal(content, userGroup)
	if err != nil {
		return nil, fmt.Errorf("UserGroup Content unmarshal error: %w", err)
	}

	for i, mEncr := range userGroup.MembersEncr {
		uID, err := key.Decrypt(mEncr)
		if err != nil {
			return nil, fmt.Errorf("UserGroup member %d ID decrypt error: %w", i, err)
		}

		userGroup.Members = append(userGroup.Members, string(uID))
	}

	return userGroup, nil
}

func (s *Service) getAccessKey(ctx context.Context, userID string, kind access.Kind, accessID []byte) (*chachaPoly.Key, error) {
	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	keyEncr, level, err := s.Infra.Index.GetUserAccess(ctx, userID, kind, accessID)
	if err != nil {
		return nil, fmt.Errorf("Index.UserGroupGetByID error: %w", err)
	}

	if level == access.NoAccess {
		return nil, errors.ErrAccessDenied
	}

	keyDecr, err := keybox.Open(keyEncr, userPubKey, userPrivKey)
	if err != nil {
		return nil, fmt.Errorf("keybox.Open error: %w", err)
	}

	key, err := chachaPoly.NewKeyFromBytes(keyDecr)
	if err != nil {
		return nil, fmt.Errorf("chachaPoly.NewKeyFromBytes error: %w", err)
	}

	return key, nil
}
