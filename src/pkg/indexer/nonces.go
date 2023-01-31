package indexer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type Noncer interface {
	Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error)
}

func (i *Index) Nonce(ctx context.Context, noncer Noncer, address *common.Address) (*big.Int, error) {
	nonce, err := noncer.Nonces(&bind.CallOpts{Context: ctx}, *address)
	if err != nil {
		return nil, fmt.Errorf("Index.Nonces error: %w address: %s", err, address.String())
	}

	if nonce == nil {
		return big.NewInt(1), nil
	}

	nonce.Add(nonce, big.NewInt(1))

	return nonce, nil
}
