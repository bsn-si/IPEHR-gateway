package indexer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/access"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/accessStore"
)

type AccessObject = accessStore.IAccessStoreAccess

func (i *Index) GetUserAccess(ctx context.Context, userIDHash *[32]byte, kind access.Kind, accessID []byte) ([]byte, access.Level, error) {
	accessIDHash := Keccak256(accessID)

	acc, err := i.accessStore.UserAccess(&bind.CallOpts{Context: ctx}, *userIDHash, kind, *accessIDHash)
	if err != nil {
		return nil, 0, fmt.Errorf("ehrIndex.UserAccess error: %w", err)
	}

	if acc.Level == access.NoAccess {
		return nil, 0, errors.ErrAccessDenied
	}

	return acc.KeyEncr, acc.Level, nil
}

func (i *Index) SetAccess(ctx context.Context, subjectIDHash *[32]byte, accessObj *AccessObject, userPrivKey *[32]byte, nonce *big.Int) (string, error) {
	userKey, err := crypto.ToECDSA(userPrivKey[:])
	if err != nil {
		return "", fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	data, err := abi.Arguments{{Type: Bytes32}, {Type: Uint8}}.Pack(*subjectIDHash, accessObj.Kind)
	if err != nil {
		return "", fmt.Errorf("args.Pack error: %w", err)
	}

	accessID := Keccak256(data)

	data, err = i.accessStoreAbi.Pack("setAccess", accessID, *accessObj, userAddress, make([]byte, signatureLength))
	if err != nil {
		return "", fmt.Errorf("abi.Pack1 error: %w", err)
	}

	if nonce == nil {
		nonce, err = i.Nonce(ctx, i.accessStore, &userAddress)
		if err != nil {
			return "", fmt.Errorf("accessNonce error: %w address: %s", err, userAddress.String())
		}
	}

	signature, err := makeSignature(data, nonce, userKey)
	if err != nil {
		return "", fmt.Errorf("makeSignature error: %w", err)
	}

	tx, err := i.accessStore.SetAccess(i.transactOpts, *accessID, *accessObj, userAddress, signature)
	if err != nil {
		return "", fmt.Errorf("accessStore.SetAccess error: %w", err)
	}

	return tx.Hash().String(), nil
}

func (i *Index) GetAccessKey(ctx context.Context, userIDHash *[32]byte, kind access.Kind, accessID []byte, userPubKey, userPrivKey *[32]byte) (*chachaPoly.Key, error) {
	keyEncr, level, err := i.GetUserAccess(ctx, userIDHash, kind, accessID)
	if err != nil {
		return nil, fmt.Errorf("Index.UserGroupGetByID error: %w", err)
	}

	if level == access.NoAccess {
		return nil, errors.ErrAccessDenied
	}

	keyDecr, err := keybox.OpenAnonymous(keyEncr, userPubKey, userPrivKey)
	if err != nil {
		return nil, fmt.Errorf("keybox.Open error: %w", err)
	}

	key, err := chachaPoly.NewKeyFromBytes(keyDecr)
	if err != nil {
		return nil, fmt.Errorf("chachaPoly.NewKeyFromBytes error: %w", err)
	}

	return key, nil
}
