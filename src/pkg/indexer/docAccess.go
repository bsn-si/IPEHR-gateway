package indexer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/bsn-si/IPEHR-gateway/src/internal/observability/tracer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/access"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/accessStore"
)

func (i *Index) GetAccessList(ctx context.Context, IDHash *[32]byte, kind access.Kind) (access.List, error) {
	ctx, span := tracer.Start(ctx, "indexer.GetAccessList") //nolint
	defer span.End()

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

func (i *Index) DocAccessSet(ctx context.Context, CID, CIDEncr, keyEncr []byte, accessLevel uint8, userPrivKey, toUserPrivKey *[32]byte) ([]byte, error) {
	ctx, span := tracer.Start(ctx, "indexer.DocAccessSet") //nolint
	defer span.End()

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

	sig := make([]byte, signatureLength)

	data, err = i.ehrIndexAbi.Pack("setDocAccess", CID, accessObj, toUserAddress, userAddress, sig)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack1 error: %w", err)
	}

	deadline := big.NewInt(time.Now().Add(i.txTimeout).Unix())

	sig, err = makeSignature(data, userKey, deadline)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err = i.ehrIndexAbi.Pack("setDocAccess", CID, accessObj, toUserAddress, userAddress, deadline, sig)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack2 error: %w", err)
	}

	return data, nil
}
