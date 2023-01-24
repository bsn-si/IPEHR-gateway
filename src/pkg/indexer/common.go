package indexer

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

func Keccak256(data []byte) *[32]byte {
	var b [32]byte

	copy(b[:], crypto.Keccak256(data))

	return &b
}

func (i *Index) TxWait(ctx context.Context, hash string) (uint64, error) {
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
	h := common.HexToHash(hash)

	receipt, err := i.client.TransactionReceipt(ctx, h)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return 0, errors.ErrIsNotExist
		}

		return 0, fmt.Errorf("GetTxStatus error: %w hash %s", err, hash)
	}

	return receipt.Status, nil
}
