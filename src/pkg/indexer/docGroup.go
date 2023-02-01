package indexer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/chachaPoly"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/crypto/keybox"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/ehrIndexer"
)

func (i *Index) DocGroupCreate(ctx context.Context, gID *uuid.UUID, gIDEncr, gKeyEncr, gNameEncr []byte, userPrivKey *[32]byte, nonce *big.Int) ([]byte, error) {
	gIDHash := Keccak256(gID[:])

	userKey, err := crypto.ToECDSA(userPrivKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	params := ehrIndexer.DocGroupsDocGroupCreateParams{
		GroupIDHash: *gIDHash,
		Attrs: []ehrIndexer.AttributesAttribute{
			{Code: model.AttributeIDEncr, Value: gIDEncr},
			{Code: model.AttributeKeyEncr, Value: gKeyEncr},
			{Code: model.AttributeNameEncr, Value: gNameEncr},
		},
		Signer:    userAddress,
		Signature: make([]byte, signatureLength),
	}

	data, err := i.ehrIndexAbi.Pack("docGroupCreate", params)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack1 error: %w", err)
	}

	if nonce == nil {
		nonce, err = i.Nonce(ctx, i.ehrIndex, &userAddress)
		if err != nil {
			return nil, fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	params.Signature, err = makeSignature(data, nonce, userKey)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err = i.ehrIndexAbi.Pack("docGroupCreate", params)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack2 error: %w", err)
	}

	return data, nil
}

// Returns: []CIDEncr
func (i *Index) DocGroupGetDocs(ctx context.Context, gID *uuid.UUID) ([][]byte, error) {
	groupIDHash := Keccak256(gID[:])

	CIDs, err := i.ehrIndex.DocGroupGetDocs(&bind.CallOpts{Context: ctx}, *groupIDHash)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.DocGroupGetDocs error: %w groupIDHash: %x", err, groupIDHash)
	}

	return CIDs, nil
}

func (i *Index) DocGroupAddDoc(ctx context.Context, gID *uuid.UUID, docCIDHash *[32]byte, docCIDEncr []byte, userPrivKey *[32]byte, nonce *big.Int) ([]byte, error) {
	groupIDHash := Keccak256(gID[:])

	userKey, err := crypto.ToECDSA(userPrivKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	data, err := i.ehrIndexAbi.Pack("docGroupAddDoc", groupIDHash, docCIDHash, docCIDEncr, userAddress, make([]byte, signatureLength))
	if err != nil {
		return nil, fmt.Errorf("abi.Pack1 error: %w", err)
	}

	if nonce == nil {
		nonce, err = i.Nonce(ctx, i.ehrIndex, &userAddress)
		if err != nil {
			return nil, fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	sig, err := makeSignature(data, nonce, userKey)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err = i.ehrIndexAbi.Pack("docGroupAddDoc", groupIDHash, docCIDHash, docCIDEncr, userAddress, sig)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack2 error: %w", err)
	}

	return data, nil
}

func (i *Index) DocGroupGetByID(ctx context.Context, gID *uuid.UUID, userPubKey, userPrivateKey *[32]byte) (*model.DocumentGroup, error) {
	groupIDHash := Keccak256(gID[:])

	attrs, err := i.ehrIndex.DocGroupGetAttrs(&bind.CallOpts{Context: ctx}, *groupIDHash)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.DocGroupGetAttrs error: %w", err)
	}

	if len(attrs) == 0 {
		return nil, errors.ErrNotFound
	}

	keyEncr := model.AttributesEhr(attrs).GetByCode(model.AttributeKeyEncr)
	if len(keyEncr) == 0 {
		return nil, errors.ErrFieldIsEmpty("KeyEncr")
	}

	keyDecr, err := keybox.OpenAnonymous(keyEncr, userPubKey, userPrivateKey)
	if err != nil {
		return nil, fmt.Errorf("docGroup key decryption error: %w", err)
	}

	key, err := chachaPoly.NewKeyFromBytes(keyDecr)
	if err != nil {
		return nil, fmt.Errorf("chachaPoly.NewKeyFromBytes error: %w", err)
	}

	nameEncr := model.AttributesEhr(attrs).GetByCode(model.AttributeNameEncr)
	if len(nameEncr) == 0 {
		return nil, errors.ErrFieldIsEmpty("NameEncr")
	}

	nameDecr, err := key.Decrypt(nameEncr)
	if err != nil {
		return nil, fmt.Errorf("docGroup name decryption error: %w", err)
	}

	return &model.DocumentGroup{
		GroupID:      *gID,
		Name:         string(nameDecr),
		GroupKeyEncr: keyEncr,
		GroupKey:     key,
	}, nil
}
