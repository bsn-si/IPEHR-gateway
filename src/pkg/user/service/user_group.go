package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/internal/observability/tracer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/access"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/compressor"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	proc "github.com/bsn-si/IPEHR-gateway/src/pkg/docs/service/processing"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
)

func (s *Service) GroupCreate(ctx context.Context, req proc.RequestInterface, userID, systemID, groupName, groupDescription string) (*model.UserGroup, error) {
	ctx, span := tracer.GetTracer().Start(ctx, "user_service.group_create") //nolint
	defer span.End()

	userPubKey, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return nil, fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	multiCallTx := s.Infra.Index.MultiCallUsersNew()

	userGroup, err := s.groupCreate(ctx, groupName, groupDescription, userPubKey, userPrivKey)
	if err != nil {
		return nil, fmt.Errorf("groupCreatePack error: %w", err)
	}

	multiCallTx.Add(uint8(proc.TxUserGroupCreate), userGroup.Packed)

	packed, err := s.setGroupAccess(ctx, userGroup, userID, systemID, access.Owner, userPrivKey)
	if err != nil {
		return nil, fmt.Errorf("setGroupAccess error: %w", err)
	}

	multiCallTx.Add(uint8(proc.TxSetUserGroupAccess), packed)

	txHash, err := multiCallTx.Commit()
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return nil, errors.ErrNotFound
		} else if strings.Contains(err.Error(), "AEX") {
			return nil, errors.ErrAlreadyExist
		}

		return nil, fmt.Errorf("multiCallTx.Commit error: %w", err)
	}

	req.AddEthereumTx(processing.TxUserGroupCreate, txHash)

	return userGroup, nil
}

func (s *Service) GroupGetByID(ctx context.Context, userID, systemID string, groupID *uuid.UUID, groupKey *chachaPoly.Key) (*model.UserGroup, error) {
	ctx, span := tracer.GetTracer().Start(ctx, "user_service.group_get_by_id") //nolint
	defer span.End()

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

	copy(userGroupResult.Key[:], groupKey.Bytes())

	return &userGroupResult, nil
}

func (s *Service) GroupAddUser(ctx context.Context, userID, systemID, addUserID, addSystemID, reqID string, level access.Level, groupID *uuid.UUID) error {
	ctx, span := tracer.GetTracer().Start(ctx, "user_service.group_add_user") //nolint
	defer span.End()

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

	txHash, err := s.Infra.Index.UserGroupAddUser(ctx, addUserID, addSystemID, level, groupID, userIDEncr, groupKeyEncr, userPrivKey)
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
	ctx, span := tracer.GetTracer().Start(ctx, "user_service.group_remove_user") //nolint
	defer span.End()

	_, userPrivKey, err := s.Infra.Keystore.Get(userID)
	if err != nil {
		return fmt.Errorf("Keystore.Get error: %w userID %s", err, userID)
	}

	multiCallTx := s.Infra.Index.MultiCallUsersNew()

	packed, err := s.Infra.Index.UserGroupRemoveUser(removeUserID, removeSystemID, groupID, userPrivKey)
	if err != nil {
		return fmt.Errorf("Index.UserGroupRemoveUser error: %w", err)
	}

	multiCallTx.Add(uint8(proc.TxUserGroupRemoveUser), packed)

	packed, err = s.setGroupAccess(ctx, &model.UserGroup{GroupID: groupID}, removeUserID, removeSystemID, access.NoAccess, userPrivKey)
	if err != nil {
		return fmt.Errorf("setGroupAccess error: %w", err)
	}

	multiCallTx.Add(uint8(proc.TxSetUserGroupAccess), packed)

	txHash, err := multiCallTx.Commit()
	if err != nil {
		return fmt.Errorf("multiCallTx.Commit error: %w", err)
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

func (s *Service) groupCreate(ctx context.Context, name, description string, userPubKey, userPrivKey *[32]byte) (*model.UserGroup, error) {
	ctx, span := tracer.GetTracer().Start(ctx, "user_service.group_create") //nolint
	defer span.End()

	groupID := uuid.New()

	userGroup := &model.UserGroup{
		GroupID:     &groupID,
		Name:        name,
		Description: description,
	}

	key := chachaPoly.GenerateKey()

	var err error

	userGroup.IDEncr, err = key.Encrypt(groupID[:])
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

	userGroup.ContentEncr, err = key.Encrypt(contentCompresed)
	if err != nil {
		return nil, fmt.Errorf("key.Encrypt content error: %w", err)
	}

	userGroup.KeyEncr, err = keybox.SealAnonymous(key.Bytes(), userPubKey)
	if err != nil {
		return nil, fmt.Errorf("keybox.Seal error: %w", err)
	}

	userGroup.Packed, err = s.Infra.Index.UserGroupCreate(ctx, &groupID, userGroup.IDEncr, userGroup.KeyEncr, userGroup.ContentEncr, userPrivKey)
	if err != nil {
		return nil, fmt.Errorf("Index.GroupCreate error: %w", err)
	}

	return userGroup, nil
}

func (s *Service) setGroupAccess(ctx context.Context, userGroup *model.UserGroup, userID, systemID string, accessLevel access.Level, userPrivKey *[32]byte) ([]byte, error) {
	ctx, span := tracer.GetTracer().Start(ctx, "user_service.set_group_access") //nolint
	defer span.End()

	userIDHash := sha3.Sum256([]byte(userID + systemID))
	groupIDHash := indexer.Keccak256(userGroup.GroupID[:])

	accessObj := indexer.AccessObject{
		Kind:    access.UserGroup,
		IdHash:  *groupIDHash,
		IdEncr:  userGroup.IDEncr,
		KeyEncr: userGroup.KeyEncr,
		Level:   accessLevel,
	}

	packed, err := s.Infra.Index.SetAccessWrapper(ctx, &userIDHash, &accessObj, userPrivKey)
	if err != nil {
		return nil, fmt.Errorf("Index.SetAccess user to composition error: %w", err)
	}

	return packed, nil
}
