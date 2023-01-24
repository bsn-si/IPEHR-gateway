package indexer

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/access"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/accessStore"
)

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

func (i *Index) SetAccess(ctx context.Context, IDHash, objectID *[32]byte, IDEncr, keyEncr []byte, kind access.Kind, level access.Level) (string, error) {
	data, err := abi.Arguments{{Type: Bytes32}, {Type: Uint8}}.Pack(*objectID, kind)
	if err != nil {
		return "", fmt.Errorf("args.Pack error: %w", err)
	}

	accessID := Keccak256(data)

	access := accessStore.IAccessStoreAccess{
		IdHash:  *IDHash,
		IdEncr:  IDEncr,
		KeyEncr: keyEncr,
		Level:   level,
	}

	tx, err := i.accessStore.SetAccess(i.transactOpts, *accessID, access)
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
