package indexer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"

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
	nonce *big.Int
	kind  MulticallKind
}

func (i *Index) MultiCallEhrNew(ctx context.Context, pk *[32]byte) (*MultiCallTx, error) {
	userKey, err := crypto.ToECDSA(pk[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	address := crypto.PubkeyToAddress(userKey.PublicKey)

	nonce, err := i.ehrIndex.Nonces(&bind.CallOpts{Context: ctx}, address)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.Nonces error: %w address: %s", err, address.String())
	}

	nonce.Add(nonce, big.NewInt(1))

	return &MultiCallTx{
		index: i,
		nonce: nonce,
		kind:  MulticallEhr,
	}, nil
}

func (i *Index) MultiCallUsersNew(ctx context.Context, pk *[32]byte) (*MultiCallTx, error) {
	userKey, err := crypto.ToECDSA(pk[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	address := crypto.PubkeyToAddress(userKey.PublicKey)

	nonce, err := i.users.Nonces(&bind.CallOpts{Context: ctx}, address)
	if err != nil {
		return nil, fmt.Errorf("users.Nonces error: %w address: %s", err, address.String())
	}

	nonce.Add(nonce, big.NewInt(1))

	return &MultiCallTx{index: i, nonce: nonce, kind: MulticallUsers}, nil
}

func (m *MultiCallTx) Add(kind uint8, packed []byte) {
	m.kinds = append(m.kinds, kind)
	m.data = append(m.data, packed)
	m.nonce.Add(m.nonce, big.NewInt(1))
}

func (m *MultiCallTx) GetTxKinds() []uint8 {
	return m.kinds
}

func (m *MultiCallTx) Nonce() *big.Int {
	return m.nonce
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
