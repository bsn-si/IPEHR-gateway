package indexer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func (i *Index) GetNonce(ctx context.Context, address *common.Address) (*big.Int, error) {
	nonce, err := i.ehrIndex.Nonces(&bind.CallOpts{Context: ctx}, *address)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.Nonces error: %w address: %s", err, address.String())
	}

	return nonce, nil
}
