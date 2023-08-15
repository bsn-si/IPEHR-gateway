package indexer

import (
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

type NoncHolder struct {
	nonce uint64
	mx    *sync.Mutex
}

func NewNoncHolder(nonce uint64) *NoncHolder {
	return &NoncHolder{
		nonce: nonce,
		mx:    &sync.Mutex{},
	}
}

func (n *NoncHolder) GetNewOpts(opts *bind.TransactOpts) *bind.TransactOpts {
	n.mx.Lock()
	nonce := n.nonce
	n.nonce++
	n.mx.Unlock()

	opts = &bind.TransactOpts{
		From:   opts.From,
		Nonce:  big.NewInt(int64(nonce)),
		Signer: opts.Signer,

		Value:     opts.Value,
		GasPrice:  opts.GasPrice,
		GasFeeCap: opts.GasFeeCap,
		GasTipCap: opts.GasTipCap,
		GasLimit:  opts.GasLimit,

		Context: opts.Context,

		NoSend: opts.NoSend,
	}

	return opts
}
