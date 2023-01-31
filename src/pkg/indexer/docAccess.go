package indexer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/access"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/accessStore"
)

func (i *Index) GetAccessList(ctx context.Context, IDHash *[32]byte, kind access.Kind) (access.List, error) {
	data, err := abi.Arguments{{Type: Bytes32}, {Type: Uint8}}.Pack(*IDHash, kind)
	if err != nil {
		return nil, fmt.Errorf("args.Pack error: %w", err)
	}

	accessID := crypto.Keccak256Hash(data)

	acl, err := i.accessStore.GetAccess(&bind.CallOpts{Context: ctx}, accessID)
	if err != nil {
		if len(acl) == 0 {
			return nil, errors.ErrNotFound
		}

		return nil, fmt.Errorf("GetUserAccessList error: %w", err)
	}

	var l access.List

	for _, a := range acl {
		IDHash := make([]byte, len(a.IdHash))
		copy(IDHash, a.IdHash[:])

		l = append(l, &access.Item{
			Fields: map[string][]byte{
				"idHash":  IDHash,
				"idEncr":  a.IdEncr,
				"keyEncr": a.KeyEncr,
				"level":   {a.Level},
			},
		})
	}

	return l, nil
}

func (i *Index) DocAccessSet(ctx context.Context, CID, CIDEncr, keyEncr []byte, accessLevel uint8, userPrivKey, toUserPrivKey *[32]byte, nonce *big.Int) ([]byte, error) {
	userKey, err := crypto.ToECDSA(userPrivKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	toUserKey, err := crypto.ToECDSA(toUserPrivKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	data, err := abi.Arguments{{Type: Bytes}}.Pack(CID)
	if err != nil {
		return nil, fmt.Errorf("args.Pack error: %w", err)
	}

	idHash := crypto.Keccak256Hash(data)

	accessObj := accessStore.IAccessStoreAccess{
		IdHash:  idHash,
		IdEncr:  CIDEncr,
		KeyEncr: keyEncr,
		Level:   accessLevel,
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)
	toUserAddress := crypto.PubkeyToAddress(toUserKey.PublicKey)

	if nonce == nil {
		nonce, err = i.Nonce(ctx, i.ehrIndex, &userAddress)
		if err != nil {
			return nil, fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	sig := make([]byte, signatureLength)

	data, err = i.ehrIndexAbi.Pack("setDocAccess", CID, accessObj, toUserAddress, userAddress, sig)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack1 error: %w", err)
	}

	sig, err = makeSignature(data, nonce, userKey)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err = i.ehrIndexAbi.Pack("setDocAccess", CID, accessObj, toUserAddress, userAddress, sig)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack2 error: %w", err)
	}

	return data, nil
}
