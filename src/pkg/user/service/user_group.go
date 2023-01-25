package service

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/access"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/compressor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
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

func (s *Service) GroupGetByID(ctx context.Context, userID, systemID string, groupID *uuid.UUID, groupKey *chachaPoly.Key) (*model.UserGroup, error) {
	var err error

	if groupKey == nil {
		userIDHash := sha3.Sum256([]byte(userID + systemID))

		userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
		if err != nil {
			return nil, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
		}

		groupKey, err = s.Infra.Index.GetAccessKey(ctx, &userIDHash, access.UserGroup, groupID[:], userPubKey, userPrivKey)
		if err != nil {
			if errors.Is(err, errors.ErrAccessDenied) {
				return nil, err
			}

			return nil, fmt.Errorf("getAccessKey error: %w", err)
		}
	}

	userGroup, err := s.Infra.Index.UserGroupGetByID(ctx, groupID)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("Index.UserGroupGetByID error: %w", err)
	}

	contentCompresed, err := groupKey.Decrypt(userGroup.ContentEncr)
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

	userGroupResult.Members = []string{}

	for i, mEncr := range userGroup.MembersEncr {
		uID, err := groupKey.Decrypt(mEncr)
		if err != nil {
			return nil, fmt.Errorf("UserGroup member %d ID decrypt error: %w", i, err)
		}

		userGroupResult.Members = append(userGroupResult.Members, string(uID))
	}

	copy(userGroupResult.GroupKey[:], groupKey.Bytes())

	return &userGroupResult, nil
}

func (s *Service) GroupAddUser(ctx context.Context, userID, systemID, addUserID, addSystemID, reqID string, level access.Level, groupID *uuid.UUID) error {
	var auID [32]byte

	copy(auID[:], addUserID)

	userIDHash := sha3.Sum256([]byte(userID + systemID))

	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	groupKey, err := s.Infra.Index.GetAccessKey(ctx, &userIDHash, access.UserGroup, groupID[:], userPubKey, userPrivKey)
	if err != nil {
		if errors.Is(err, errors.ErrAccessDenied) {
			return err
		}

		return fmt.Errorf("getAccessKey error: %w", err)
	}

	userIDEncr, err := groupKey.Encrypt(auID[:])
	if err != nil {
		return fmt.Errorf("key.Encrypt addUserID error: %w", err)
	}

	addUserPubKey, _, err := s.Infra.Keystore.Get(addUserID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, addUserID)
	}

	groupKeyEncr, err := keybox.SealAnonymous(groupKey.Bytes(), addUserPubKey)
	if err != nil {
		return fmt.Errorf("keybox.Seal error: %w", err)
	}

	txHash, err := s.Infra.Index.UserGroupAddUser(ctx, addUserID, addSystemID, level, groupID, userIDEncr, groupKeyEncr, userPrivKey, nil)
	if err != nil {
		if errors.Is(err, errors.ErrAccessDenied) {
			return err
		} else if errors.Is(err, errors.ErrAlreadyExist) {
			return err
		}

		return fmt.Errorf("Index.UserGroupAddUser error: %w", err)
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

func (s *Service) GroupRemoveUser(ctx context.Context, userID, systemID, removeUserID, removeSystemID, reqID string, groupID *uuid.UUID) error {
	_, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	txHash, err := s.Infra.Index.UserGroupRemoveUser(ctx, removeUserID, removeSystemID, groupID, userPrivKey, nil)
	if err != nil {
		if errors.Is(err, errors.ErrAccessDenied) {
			return err
		} else if errors.Is(err, errors.ErrNotFound) {
			return err
		}

		return fmt.Errorf("Index.UserGroupRemoveUser error: %w", err)
	}

	procRequest, err := s.Proc.NewRequest(reqID, userID, "", processing.RequestUserGroupRemoveUser)
	if err != nil {
		return fmt.Errorf("Proc.NewRequest error: %w", err)
	}

	procRequest.AddEthereumTx(processing.TxUserGroupRemoveUser, txHash)

	if err := procRequest.Commit(); err != nil {
		return fmt.Errorf("Remove user from group procRequest commit error: %w", err)
	}

	return nil
}

func (s *Service) GroupGetList(ctx context.Context, userID, systemID string) ([]*model.UserGroup, error) {
	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("keystore.Get error: %w userID %s", err, userID)
	}

	IDHash := sha3.Sum256([]byte(userID + systemID))

	acl, err := s.Infra.Index.GetAccessList(ctx, &IDHash, access.UserGroup)
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}

		return nil, fmt.Errorf("Index.GetAccessList error: %w userID: %s", err, userID)
	}

	var userGroupList []*model.UserGroup

	for i, a := range acl {
		err := access.ExtractWithUserKey(a, userPubKey, userPrivKey)
		if err != nil {
			if errors.Is(err, errors.ErrAccessDenied) {
				continue
			}

			return nil, fmt.Errorf("index: %d access.Extract error: %w", i, err)
		}

		groupUUID, err := uuid.FromBytes(a.ID)
		if err != nil {
			return nil, fmt.Errorf("groupID %d uuid.ParseBytes error: %w idDecr: %x", i, err, a.ID)
		}

		userGroup, err := s.GroupGetByID(ctx, userID, systemID, &groupUUID, a.Key)
		if err != nil {
			return nil, fmt.Errorf("GroupGetByID error: %w groupUUID: %s", err, groupUUID)
		}

		userGroupList = append(userGroupList, userGroup)
	}

	return userGroupList, nil
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

	//keyEncr, err := keybox.Seal(key.Bytes(), userPubKey, userPrivKey)
	keyEncr, err := keybox.SealAnonymous(key.Bytes(), userPubKey)
	if err != nil {
		return nil, nil, fmt.Errorf("keybox.Seal error: %w", err)
	}

	packed, err := s.Infra.Index.UserGroupCreate(ctx, &groupID, idEncr, keyEncr, contentEncr, userPrivKey, nonce)
	if err != nil {
		return nil, nil, fmt.Errorf("Index.GroupCreate error: %w", err)
	}

	return packed, userGroup.GroupID, nil
}
