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

	"hms/gateway/pkg/access"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer/ehrIndexer"
	userModel "hms/gateway/pkg/user/model"
)

func (i *Index) UserGroupCreate(ctx context.Context, groupID *uuid.UUID, idEncr, keyEncr, contentEncr []byte, privKey *[32]byte, nonce *big.Int) (string, error) {
	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return "", fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	if nonce == nil {
		nonce, err = i.userNonce(ctx, &userAddress)
		if err != nil {
			return "", fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	params := ehrIndexer.UsersUserGroupCreateParams{
		GroupIdHash: sha3.Sum256(groupID[:]),
		Attrs: []ehrIndexer.AttributesAttribute{
			{Code: model.AttributeKeyEncr, Value: keyEncr},         // encrypted by userKey
			{Code: model.AttributeIDEncr, Value: idEncr},           // encrypted by group key
			{Code: model.AttributeContentEncr, Value: contentEncr}, // encrypted by group key
		},
		Signer:    userAddress,
		Signature: make([]byte, signatureLength),
	}

	data, err := i.abi.Pack("userGroupCreate", params)
	if err != nil {
		return "", fmt.Errorf("abi.Pack error: %w", err)
	}

	params.Signature, err = makeSignature(data, nonce, userKey)
	if err != nil {
		return "", fmt.Errorf("makeSignature error: %w", err)
	}

	tx, err := i.ehrIndex.UserGroupCreate(i.transactOpts, params)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return "", errors.ErrNotFound
		} else if strings.Contains(err.Error(), "AEX") {
			return "", errors.ErrAlreadyExist
		}

		return "", fmt.Errorf("ehrIndex.UserGroupCreate error: %w", err)
	}

	return tx.Hash().Hex(), nil
}

func (i *Index) UserGroupGetByID(ctx context.Context, userID string, groupID *uuid.UUID) (*userModel.UserGroup, error) {
	groupIDHash := sha3.Sum256(groupID[:])

	ug, err := i.ehrIndex.UserGroupGetByID(&bind.CallOpts{Context: ctx}, groupIDHash)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.UserGroupGetByID error: %w", err)
	}

	if len(ug.Attrs) == 0 {
		return nil, errors.ErrNotFound
	}

	contentEncr := model.Attributes(ug.Attrs).GetByCode(model.AttributeContentEncr)
	if contentEncr == nil {
		return nil, errors.ErrFieldIsEmpty("ContentEncr")
	}

	userGroup := &userModel.UserGroup{
		GroupID:     groupID,
		ContentEncr: contentEncr,
		Members:     []string{},
	}

	for _, m := range ug.Members {
		userGroup.MembersEncr = append(userGroup.MembersEncr, m.UserIDEncr)
	}

	return userGroup, nil
}

func (i *Index) UserGroupAddUser(ctx context.Context, userID string, level access.Level, groupID *uuid.UUID, userIDEncr, groupKeyEncr []byte, privKey *[32]byte, nonce *big.Int) (string, error) {
	var uID [32]byte

	copy(uID[:], userID)

	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return "", fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	if nonce == nil {
		nonce, err = i.userNonce(ctx, &userAddress)
		if err != nil {
			return "", fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	params := ehrIndexer.UsersGroupAddUserParams{
		GroupIDHash: sha3.Sum256(groupID[:]),
		UserIDHash:  sha3.Sum256(uID[:]),
		Level:       level,
		UserIDEncr:  userIDEncr,
		KeyEncr:     groupKeyEncr,
		Signer:      userAddress,
		Signature:   make([]byte, signatureLength),
	}

	data, err := i.abi.Pack("groupAddUser", params)
	if err != nil {
		return "", fmt.Errorf("abi.Pack error: %w", err)
	}

	params.Signature, err = makeSignature(data, nonce, userKey)
	if err != nil {
		return "", fmt.Errorf("makeSignature error: %w", err)
	}

	tx, err := i.ehrIndex.GroupAddUser(i.transactOpts, params)
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
