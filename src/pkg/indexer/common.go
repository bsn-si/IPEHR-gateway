package indexer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bsn-si/IPEHR-gateway/src/internal/observability/tracer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

const (
	TxTimeout = 15 * time.Minute
)

func Keccak256(data []byte) *[32]byte {
	var b [32]byte

	copy(b[:], crypto.Keccak256(data))

	return &b
}

func (i *Index) TxWait(ctx context.Context, hash string) (uint64, error) {
	ctx, span := tracer.Start(ctx, "indexer.tx_wait") //nolint
	defer span.End()

	h := common.HexToHash(hash)

	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case <-ticker.C:
			receipt, err := i.client.TransactionReceipt(ctx, h)

			switch {
			case err != nil && !errors.Is(err, ethereum.NotFound):
				return 0, err
			case err == nil:
				return receipt.Status, nil
			default:
			}
		case <-ctx.Done():
			return 0, errors.ErrTimeout
		}
	}
}

func (i *Index) GetTxStatus(ctx context.Context, hash string) (uint64, error) {
	ctx, span := tracer.Start(ctx, "indexer.GetTxStatus") //nolint
	defer span.End()

	r, err := i.client.TransactionReceipt(ctx, common.HexToHash(hash))
	if err != nil {
		return 0, err
	}

	return r.Status, nil
}

func (i *Index) GetTxReceipt(ctx context.Context, hash string) (*types.Receipt, error) {
	ctx, span := tracer.Start(ctx, "indexer.GetTxReceipt") //nolint
	defer span.End()

	h := common.HexToHash(hash)

	receipt, err := i.client.TransactionReceipt(ctx, h)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return nil, errors.ErrIsNotExist
		}

		return nil, fmt.Errorf("GetTxReceipt error: %w hash %s", err, hash)
	}

	return receipt, nil
}

func (i *Index) GetNewOpts(opts *bind.TransactOpts) *bind.TransactOpts {
	i.Lock()
	defer i.Unlock()

	opts = &bind.TransactOpts{
		From:   opts.From,
		Nonce:  big.NewInt(opts.Nonce.Int64()),
		Signer: opts.Signer,

		Value:     opts.Value,
		GasPrice:  opts.GasPrice,
		GasFeeCap: opts.GasFeeCap,
		GasTipCap: opts.GasTipCap,
		GasLimit:  opts.GasLimit,

		Context: opts.Context,

		NoSend: opts.NoSend,
	}

	i.transactOpts.Nonce.Add(opts.Nonce, big.NewInt(1))

	return opts
}
