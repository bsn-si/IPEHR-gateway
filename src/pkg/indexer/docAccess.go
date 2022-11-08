package indexer

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"

	"hms/gateway/pkg/access"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer/ehrIndexer"
)

func (i *Index) DocAccessList(ctx context.Context, userID string) (access.List, error) {
	var uID [32]byte

	copy(uID[:], []byte(userID))

	data, err := abi.Arguments{{Type: Bytes32}, {Type: Uint8}}.Pack(uID, access.Doc)
	if err != nil {
		return nil, fmt.Errorf("args.Pack error: %w", err)
	}

	accessID := crypto.Keccak256Hash(data)

	acl, err := i.ehrIndex.GetUserAccessList(&bind.CallOpts{Context: ctx}, accessID)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return nil, errors.ErrNotFound
		}

		return nil, fmt.Errorf("GetUserAccessList error: %w", err)
	}

	var l access.List

	for _, a := range acl {
		idHash := make([]byte, len(a.IdHash))
		copy(idHash, a.IdHash[:])

		level := []byte{a.Level}

		l = append(l, &access.Item{
			Fields: map[string][]byte{
				"idHash":  idHash,
				"idEncr":  a.IdEncr,
				"keyEncr": a.KeyEncr,
				"level":   level,
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

	accessObj := ehrIndexer.AccessObject{
		IdHash:  idHash,
		IdEncr:  CIDEncr,
		KeyEncr: keyEncr,
		Level:   accessLevel,
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)
	toUserAddress := crypto.PubkeyToAddress(toUserKey.PublicKey)

	if nonce == nil {
		nonce, err = i.userNonce(ctx, &userAddress)
		if err != nil {
			return nil, fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	sig, err := makeSignature(
		userKey,
		abi.Arguments{{Type: String}, {Type: Bytes}, {Type: Access}, {Type: Address}, {Type: Uint256}},
		"setDocAccess", CID, accessObj, toUserAddress, nonce,
	)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err = i.abi.Pack("setDocAccess", CID, accessObj, toUserAddress, nonce, userAddress, sig)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack error: %w", err)
	}

	return data, nil
}
