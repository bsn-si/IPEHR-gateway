package fakeData

import (
	"math/rand"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

func Cid() *cid.Cid {
	p := make([]byte, 256)

	// nolint
	_, err := rand.Read(p)
	if err != nil {
		panic(err)
	}

	h, err := multihash.Sum(p, multihash.SHA3, 4)
	if err != nil {
		panic(err)
	}

	CID := cid.NewCidV1(7, h)

	return &CID
}
