package indexer

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/access"
	"hms/gateway/pkg/errors"
)

func (i *Index) GetUserAccess(ctx context.Context, userID, systemID string, kind access.Kind, accessID []byte) ([]byte, access.Level, error) {
	userIDHash := sha3.Sum256([]byte(userID + systemID))
	accessIDHash := sha3.Sum256(accessID)

	acc, err := i.accessStore.UserAccess(&bind.CallOpts{Context: ctx}, userIDHash, kind, accessIDHash)
	if err != nil {
		return nil, 0, fmt.Errorf("ehrIndex.UserAccess error: %w", err)
	}

	if acc.Level == access.NoAccess {
		return nil, 0, errors.ErrAccessDenied
	}

	return acc.KeyEncr, acc.Level, nil
}
