package indexer

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/access"
	"hms/gateway/pkg/errors"
)

func (i *Index) GetUserAccess(ctx context.Context, userID string, kind access.Kind, accessID []byte) ([]byte, access.Level, error) {
	var uID [32]byte

	copy(uID[:], userID)

	IDHash := sha3.Sum256(accessID)

	acc, err := i.ehrIndex.UserAccess(&bind.CallOpts{Context: ctx}, uID, kind, IDHash)
	if err != nil {
		return nil, 0, fmt.Errorf("ehrIndex.UserAccess error: %w", err)
	}

	if acc.Level == access.NoAccess {
		return nil, 0, errors.ErrAccessDenied
	}

	return acc.KeyEncr, acc.Level, nil
}
