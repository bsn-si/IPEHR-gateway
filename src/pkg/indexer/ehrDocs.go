package indexer

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer/ehrIndexer"
)

func (i *Index) AddEhrDoc(ctx context.Context, docType types.DocumentType, docMeta *model.DocumentMeta, privKey *[32]byte, nonce *big.Int) ([]byte, error) {
	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	if nonce == nil {
		nonce, err = i.userNonce(ctx, &userAddress)
		if err != nil {
			return nil, fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	params := ehrIndexer.DocsAddEhrDocParams{
		DocType:   uint8(docType),
		Id:        docMeta.Id,
		Version:   docMeta.Version,
		Timestamp: docMeta.Timestamp,
		Attrs:     docMeta.Attrs,
		Signer:    userAddress,
		Signature: make([]byte, signatureLength),
	}

	data, err := i.abi.Pack("addEhrDoc", params)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack1 error: %w", err)
	}

	params.Signature, err = makeSignature(data, nonce, userKey)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err = i.abi.Pack("addEhrDoc", params)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack2 error: %w", err)
	}

	return data, nil
}

func (i *Index) GetDocLastByType(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType) (*model.DocumentMeta, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		eID      [32]byte
	)

	copy(eID[:], ehrUUID[:])

	docMeta, err := i.ehrIndex.GetLastEhrDocByType(callOpts, eID, uint8(docType))
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return nil, fmt.Errorf("ehrIndex.GetLastEhrDocByType error: %w", errors.ErrNotFound)
		}
		return nil, fmt.Errorf("ehrIndex.GetLastEhrDocByType error: %w ehrUUID %s docType %s", err, ehrUUID.String(), docType.String())
	}

	return (*model.DocumentMeta)(&docMeta), nil
}

func (i *Index) ListDocByType(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType) ([]model.DocumentMeta, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		eID      [32]byte
	)

	copy(eID[:], ehrUUID[:])

	docsMeta, err := i.ehrIndex.GetEhrDocs(callOpts, eID, uint8(docType))
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return nil, fmt.Errorf("ehrIndex.GetEhrDocs error: %w", errors.ErrNotFound)
		}
		return nil, fmt.Errorf("ehrIndex.GetEhrDocs error: %w ehrUUID %s docType %s", err, ehrUUID.String(), docType.String())
	}

	var list []model.DocumentMeta

	for _, dm := range docsMeta {
		list = append(list, model.DocumentMeta(dm))
	}

	return list, nil
}

func (i *Index) GetDocLastByBaseID(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType, docBaseUIDHash *[32]byte) (*model.DocumentMeta, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		eID      [32]byte
	)

	copy(eID[:], ehrUUID[:])

	docMeta, err := i.ehrIndex.GetDocLastByBaseID(callOpts, eID, uint8(docType), *docBaseUIDHash)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return nil, fmt.Errorf("ehrIndex.GetDocLastByBaseID error: %w", errors.ErrNotFound)
		}
		return nil, fmt.Errorf("ehrIndex.GetDocLastByBaseID error: %w ehrUUID %s docType %s docBaseUIDHash %x", err, ehrUUID.String(), docType.String(), docBaseUIDHash)
	}

	return (*model.DocumentMeta)(&docMeta), nil
}

func (i *Index) GetDocByTime(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType, timestamp uint32) (*model.DocumentMeta, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		eID      [32]byte
	)

	copy(eID[:], ehrUUID[:])

	docMeta, err := i.ehrIndex.GetDocByTime(callOpts, eID, uint8(docType), timestamp)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return nil, fmt.Errorf("ehrIndex.GetDocByTime error: %w", errors.ErrNotFound)
		}
		return nil, fmt.Errorf("ehrIndex.GetDocByTime error: %w ehrUUID %s docType %s timestamp %d", err, ehrUUID.String(), docType.String(), timestamp)
	}

	return (*model.DocumentMeta)(&docMeta), nil
}

func (i *Index) GetDocByVersion(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType, docBaseUIDHash, version *[32]byte) (*model.DocumentMeta, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		eID      [32]byte
	)

	copy(eID[:], ehrUUID[:])

	docMeta, err := i.ehrIndex.GetDocByVersion(callOpts, eID, uint8(docType), *docBaseUIDHash, *version)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("ehrIndex.GetDocByVersion error: %w ehrUUID %s docType %s docBaseUIDHash %x version %s", err, ehrUUID.String(), docType.String(), docBaseUIDHash, version)
	}

	return (*model.DocumentMeta)(&docMeta), nil
}

func (i *Index) SetEhrSubject(ctx context.Context, ehrUUID *uuid.UUID, subjectID, subjectNamespace string, privKey *[32]byte, nonce *big.Int) ([]byte, error) {
	var eID [32]byte

	copy(eID[:], ehrUUID[:])

	subjectKey := sha3.Sum256([]byte(subjectID + subjectNamespace))

	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	if nonce == nil {
		nonce, err = i.userNonce(ctx, &userAddress)
		if err != nil {
			return nil, fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	sig := make([]byte, 65)

	data, err := i.abi.Pack("setEhrSubject", subjectKey, eID, userAddress, sig)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack1 error: %w", err)
	}

	sig, err = makeSignature(data, nonce, userKey)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err = i.abi.Pack("setEhrSubject", subjectKey, eID, userAddress, sig)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack2 error: %w", err)
	}

	return data, nil
}

func (i *Index) GetEhrUUIDBySubject(ctx context.Context, subjectID, subjectNamespace string) (*uuid.UUID, error) {
	subjectKey := sha3.Sum256([]byte(subjectID + subjectNamespace))

	callOpts := &bind.CallOpts{
		Context: ctx,
	}

	ehrUUIDRaw, err := i.ehrIndex.EhrSubject(callOpts, subjectKey)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.EhrSubjec error: %w", err)
	}

	ehrUUID, err := uuid.FromBytes(ehrUUIDRaw[:16])
	if err != nil {
		return nil, fmt.Errorf("ehrUUID FromBytes error: %w ehrUUIDRaw %x", err, ehrUUIDRaw)
	}

	return &ehrUUID, nil
}

func (i *Index) DeleteDoc(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType, docBaseUIDHash, version, privKey *[32]byte, nonce *big.Int) (string, error) {
	var eID [32]byte

	copy(eID[:], ehrUUID[:])

	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return "", fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	i.Lock()
	defer i.Unlock()

	if nonce == nil {
		nonce, err = i.userNonce(ctx, &userAddress)
		if err != nil {
			return "", fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	sig := make([]byte, 65)

	data, err := i.abi.Pack("deleteDoc", eID, uint8(docType), *docBaseUIDHash, version, userAddress, sig)
	if err != nil {
		return "", fmt.Errorf("abi.Pack error: %w", err)
	}

	sig, err = makeSignature(data, nonce, userKey)
	if err != nil {
		return "", fmt.Errorf("makeSignature error: %w", err)
	}

	tx, err := i.ehrIndex.DeleteDoc(i.transactOpts, eID, uint8(docType), *docBaseUIDHash, *version, userAddress, sig)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return "", errors.ErrNotFound
		} else if strings.Contains(err.Error(), "ADL") {
			return "", errors.ErrAlreadyDeleted
		}
		return "", fmt.Errorf("ehrIndex.DeleteDoc error: %w ehrUUID %s docType %s", err, ehrUUID.String(), docType.String())
	}

	return tx.Hash().Hex(), nil
}

func (i *Index) GetEhrIDByUserID(ctx context.Context, userID string) (*uuid.UUID, error) {
	var uID [32]byte

	copy(uID[:], userID)

	ehrUUIDRaw, err := i.ehrIndex.EhrUsers(&bind.CallOpts{Context: ctx}, uID)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.EhrUsers error: %w", err)
	}

	if ehrUUIDRaw == [32]byte{} {
		return nil, errors.ErrNotFound
	}

	ehrUUID, err := uuid.FromBytes(ehrUUIDRaw[:16])
	if err != nil {
		return nil, fmt.Errorf("ehrUUID FromBytes error: %w ehrUUIDRaw %x", err, ehrUUIDRaw)
	}

	return &ehrUUID, nil
}