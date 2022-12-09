package service

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"

	"hms/gateway/pkg/access"
	"hms/gateway/pkg/compressor"
	"hms/gateway/pkg/crypto/chachaPoly"
	"hms/gateway/pkg/crypto/keybox"
	"hms/gateway/pkg/docs/service/processing"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer"
	"hms/gateway/pkg/user/model"
)

func (s *Service) GroupCreate(ctx context.Context, userID, name, description string) (string, *uuid.UUID, error) {
	packed, groupUUID, err := s.groupCreatePack(ctx, userID, name, description, nil)
	if err != nil {
		return "", nil, fmt.Errorf("groupCreatePack error: %w", err)
	}

	txHash, err := s.Infra.Index.SendSingle(ctx, packed, indexer.MulticallUsers)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return "", nil, errors.ErrNotFound
		} else if strings.Contains(err.Error(), "AEX") {
			return "", nil, errors.ErrAlreadyExist
		}

		return "", nil, fmt.Errorf("Index.SendSingle error: %w", err)
	}

	return txHash, groupUUID, nil
}

func (s *Service) GroupGetByID(ctx context.Context, userID, systemID string, groupID *uuid.UUID) (*model.UserGroup, error) {
	key, err := s.getAccessKey(ctx, userID, systemID, access.UserGroup, groupID[:])
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

	var userGroupResult model.UserGroup

	err = msgpack.Unmarshal(content, &userGroupResult)
	if err != nil {
		return nil, fmt.Errorf("UserGroup Content unmarshal error: %w", err)
	}

	for i, mEncr := range userGroup.MembersEncr {
		uID, err := key.Decrypt(mEncr)
		if err != nil {
			return nil, fmt.Errorf("UserGroup member %d ID decrypt error: %w", i, err)
		}

		userGroupResult.Members = append(userGroupResult.Members, string(uID))
	}

	return &userGroupResult, nil
}

func (s *Service) GroupAddUser(ctx context.Context, userID, systemID, addingUserID, reqID string, level access.Level, groupID *uuid.UUID) error {
	var auID [32]byte

	copy(auID[:], addingUserID)

	groupKey, err := s.getAccessKey(ctx, userID, systemID, access.UserGroup, groupID[:])
	if err != nil {
		if errors.Is(err, errors.ErrAccessDenied) {
			return err
		}

		return fmt.Errorf("getAccessKey error: %w", err)
	}

	userIDEncr, err := groupKey.Encrypt(auID[:])
	if err != nil {
		return fmt.Errorf("key.Encrypt addingUserID error: %w", err)
	}

	_, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	addingUserPubKey, addingUserPrivKey, err := s.Infra.Keystore.Get(addingUserID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, addingUserID)
	}

	groupKeyEncr, err := keybox.Seal(groupKey.Bytes(), addingUserPubKey, addingUserPrivKey)
	if err != nil {
		return fmt.Errorf("keybox.Seal error: %w", err)
	}

	txHash, err := s.Infra.Index.UserGroupAddUser(ctx, userID, level, groupID, userIDEncr, groupKeyEncr, userPrivKey, nil)
	if err != nil {
		if errors.Is(err, errors.ErrAccessDenied) {
			return err
		} else if errors.Is(err, errors.ErrAlreadyExist) {
			return err
		}

		return fmt.Errorf("Index.GroupCreate error: %w", err)
	}

	procRequest, err := s.Proc.NewRequest(reqID, userID, "", processing.RequestUserGroupAddUser)
	if err != nil {
		return fmt.Errorf("Proc.NewRequest error: %w", err)
	}

	procRequest.AddEthereumTx(processing.TxUserGroupAddUser, txHash)

	if err := procRequest.Commit(); err != nil {
		return fmt.Errorf("Add user to group procRequest commit error: %w", err)
	}

	return nil
}

func (s *Service) groupCreatePack(ctx context.Context, userID, name, description string, nonce *big.Int) ([]byte, *uuid.UUID, error) {
	groupID := uuid.New()

	userGroup := &model.UserGroup{
		GroupID:     &groupID,
		Name:        name,
		Description: description,
	}

	key := chachaPoly.GenerateKey()

	idEncr, err := key.Encrypt(groupID[:])
	if err != nil {
		return nil, nil, fmt.Errorf("key.Encrypt groupID error: %w", err)
	}

	content, err := msgpack.Marshal(userGroup)
	if err != nil {
		return nil, nil, fmt.Errorf("msgpack.Marshal error: %w", err)
	}

	contentCompresed, err := compressor.New(compressor.BestCompression).Compress(content)
	if err != nil {
		return nil, nil, fmt.Errorf("UserGroup content compression error: %w", err)
	}

	contentEncr, err := key.Encrypt(contentCompresed)
	if err != nil {
		return nil, nil, fmt.Errorf("key.Encrypt content error: %w", err)
	}

	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	keyEncr, err := keybox.Seal(key.Bytes(), userPubKey, userPrivKey)
	if err != nil {
		return nil, nil, fmt.Errorf("keybox.Seal error: %w", err)
	}

	packed, err := s.Infra.Index.UserGroupCreate(ctx, &groupID, idEncr, keyEncr, contentEncr, userPrivKey, nonce)
	if err != nil {
		return nil, nil, fmt.Errorf("Index.GroupCreate error: %w", err)
	}

	return packed, userGroup.GroupID, nil
}

func (s *Service) getAccessKey(ctx context.Context, userID, systemID string, kind access.Kind, accessID []byte) (*chachaPoly.Key, error) {
	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	keyEncr, level, err := s.Infra.Index.GetUserAccess(ctx, userID, systemID, kind, accessID)
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
