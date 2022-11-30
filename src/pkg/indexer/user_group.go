package indexer

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer/ehrIndexer"
)

func (i *Index) GroupCreate(ctx context.Context, groupID *uuid.UUID, idEncr, keyEncr, contentEncr []byte, privKey *[32]byte, nonce *big.Int) (string, error) {
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
		Signature: make([]byte, 65),
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
