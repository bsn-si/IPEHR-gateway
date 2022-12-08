package indexer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func (i *Index) usersNonce(ctx context.Context, address *common.Address) (*big.Int, error) {
	nonce, err := i.users.Nonces(&bind.CallOpts{Context: ctx}, *address)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.Nonces error: %w address: %s", err, address.String())
	}

	if nonce == nil {
		return big.NewInt(1), nil
	}

	nonce.Add(nonce, big.NewInt(1))

	return nonce, nil
}

func (i *Index) ehrNonce(ctx context.Context, address *common.Address) (*big.Int, error) {
	nonce, err := i.ehrIndex.Nonces(&bind.CallOpts{Context: ctx}, *address)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.Nonces error: %w address: %s", err, address.String())
	}

	if nonce == nil {
		return big.NewInt(1), nil
	}

	nonce.Add(nonce, big.NewInt(1))

	return nonce, nil
}
