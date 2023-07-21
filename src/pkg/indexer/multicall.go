package indexer

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/core/types"

	"github.com/bsn-si/IPEHR-gateway/src/internal/observability/tracer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
)

type MulticallKind uint8

const (
	MulticallEhr MulticallKind = iota
	MulticallUsers
)

type MultiCallTx struct {
	index *Index
	kinds []uint8
	data  [][]byte
	kind  MulticallKind
}

func (i *Index) MultiCallEhrNew() *MultiCallTx {
	return &MultiCallTx{index: i, kind: MulticallEhr}
}

func (i *Index) MultiCallUsersNew() *MultiCallTx {
	return &MultiCallTx{index: i, kind: MulticallUsers}
}

func (m *MultiCallTx) Add(kind uint8, packed []byte) {
	m.kinds = append(m.kinds, kind)
	m.data = append(m.data, packed)
}

func (m *MultiCallTx) GetTxKinds() []uint8 {
	return m.kinds
}

func (m *MultiCallTx) Commit() (string, error) {
	if len(m.data) == 0 {
		return "", fmt.Errorf("%w MultiCallTx data is empty", errors.ErrCustom)
	}

	var (
		tx  *types.Transaction
		err error
	)

	switch m.kind {
	case MulticallEhr:
		tx, err = m.index.ehrIndex.Multicall(m.index.transactOpts, m.data)
	case MulticallUsers:
		tx, err = m.index.users.Multicall(m.index.transactOpts, m.data)
	default:
		return "", fmt.Errorf("%w: unknown kind %d", err, m.kind)
	}

	if err != nil {
		return "", fmt.Errorf("Multicall error: %w", err)
	}

	return tx.Hash().Hex(), nil
}

func (i *Index) SendSingle(ctx context.Context, data []byte, kind MulticallKind) (string, error) {
	ctx, span := tracer.Start(ctx, "indexer.SendSingle") //nolint
	defer span.End()

	var (
		tx  *types.Transaction
		err error
	)

	switch kind {
	case MulticallEhr:
		tx, err = i.ehrIndex.Multicall(i.transactOpts, [][]byte{data})
	case MulticallUsers:
		tx, err = i.users.Multicall(i.transactOpts, [][]byte{data})
	}

	if err != nil {
		return "", fmt.Errorf("Multicall error: %w", err)
	}

	return tx.Hash().String(), nil
}
