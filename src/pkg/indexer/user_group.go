package indexer

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/access"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/users"
	userModel "github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
)

func (i *Index) UserGroupCreate(ctx context.Context, groupID *uuid.UUID, idEncr, keyEncr, contentEncr []byte, privKey *[32]byte, nonce *big.Int) ([]byte, error) {
	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	if nonce == nil {
		nonce, err = i.Nonce(ctx, i.users, &userAddress)
		if err != nil {
			return nil, fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	attrs := []users.AttributesAttribute{
		{Code: model.AttributeKeyEncr, Value: keyEncr},         // encrypted by userKey
		{Code: model.AttributeIDEncr, Value: idEncr},           // encrypted by group key
		{Code: model.AttributeContentEncr, Value: contentEncr}, // encrypted by group key
	}

	IDHash := Keccak256(groupID[:])

	data, err := i.usersAbi.Pack("userGroupCreate", IDHash, attrs, userAddress, make([]byte, signatureLength))
	if err != nil {
		return nil, fmt.Errorf("abi.Pack1 error: %w", err)
	}

	signature, err := makeSignature(data, nonce, userKey)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err = i.usersAbi.Pack("userGroupCreate", IDHash, attrs, userAddress, signature)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack2 error: %w", err)
	}

	return data, nil
}

func (i *Index) UserGroupGetByID(ctx context.Context, groupID *uuid.UUID) (*userModel.UserGroup, error) {
	groupIDHash := Keccak256(groupID[:])

	ug, err := i.users.UserGroupGetByID(&bind.CallOpts{Context: ctx}, *groupIDHash)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.UserGroupGetByID error: %w", err)
	}

	if len(ug.Attrs) == 0 {
		return nil, errors.ErrNotFound
	}

	contentEncr := model.AttributesUsers(ug.Attrs).GetByCode(model.AttributeContentEncr)
	if contentEncr == nil {
		return nil, errors.ErrFieldIsEmpty("ContentEncr")
	}

	groupKeyEncr := model.AttributesUsers(ug.Attrs).GetByCode(model.AttributeKeyEncr)
	if groupKeyEncr == nil {
		return nil, errors.ErrFieldIsEmpty("KeyEncr")
	}

	userGroup := &userModel.UserGroup{
		GroupID:      groupID,
		ContentEncr:  contentEncr,
		GroupKeyEncr: groupKeyEncr,
		Members:      []string{},
	}

	for _, m := range ug.Members {
		userGroup.MembersEncr = append(userGroup.MembersEncr, m.UserIDEncr)
	}

	return userGroup, nil
}

func (i *Index) UserGroupAddUser(ctx context.Context, addUserID, addSystemID string, level access.Level, groupID *uuid.UUID, addingUserIDEncr, groupKeyEncr []byte, privKey *[32]byte, nonce *big.Int) (string, error) {
	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return "", fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	if nonce == nil {
		nonce, err = i.Nonce(ctx, i.users, &userAddress)
		if err != nil {
			return "", fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	groupIDHash := Keccak256(groupID[:])

	params := users.IUsersGroupAddUserParams{
		GroupIDHash: *groupIDHash,
		UserIDHash:  sha3.Sum256([]byte(addUserID + addSystemID)),
		Level:       level,
		UserIDEncr:  addingUserIDEncr,
		KeyEncr:     groupKeyEncr,
		Signer:      userAddress,
		Signature:   make([]byte, signatureLength),
	}

	data, err := i.usersAbi.Pack("groupAddUser", params)
	if err != nil {
		return "", fmt.Errorf("abi.Pack error: %w", err)
	}

	params.Signature, err = makeSignature(data, nonce, userKey)
	if err != nil {
		return "", fmt.Errorf("makeSignature error: %w", err)
	}

	tx, err := i.users.GroupAddUser(i.transactOpts, params)
	if err != nil {
		if strings.Contains(err.Error(), "DNY") {
			return "", errors.ErrAccessDenied
		} else if strings.Contains(err.Error(), "AEX") {
			return "", errors.ErrAlreadyExist
		}

		return "", fmt.Errorf("ehrIndex.UserGroupCreate error: %w", err)
	}

	return tx.Hash().Hex(), nil
}

func (i *Index) UserGroupRemoveUser(ctx context.Context, removeUserID, removeSystemID string, groupID *uuid.UUID, privKey *[32]byte, nonce *big.Int) (string, error) {
	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return "", fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	if nonce == nil {
		nonce, err = i.Nonce(ctx, i.users, &userAddress)
		if err != nil {
			return "", fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	groupIDHash := Keccak256(groupID[:])
	removeUserIDHash := sha3.Sum256([]byte(removeUserID + removeSystemID))

	data, err := i.usersAbi.Pack("groupRemoveUser", groupIDHash, removeUserIDHash, userAddress, make([]byte, signatureLength))
	if err != nil {
		return "", fmt.Errorf("abi.Pack error: %w", err)
	}

	signature, err := makeSignature(data, nonce, userKey)
	if err != nil {
		return "", fmt.Errorf("makeSignature error: %w", err)
	}

	tx, err := i.users.GroupRemoveUser(i.transactOpts, *groupIDHash, removeUserIDHash, userAddress, signature)
	if err != nil {
		if strings.Contains(err.Error(), "DNY") {
			return "", errors.ErrAccessDenied
		} else if strings.Contains(err.Error(), "NFD") {
			return "", errors.ErrNotFound
		}

		return "", fmt.Errorf("users.GroupRemoveUser error: %w", err)
	}

	return tx.Hash().Hex(), nil
}
