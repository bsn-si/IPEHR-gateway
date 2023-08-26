package indexer

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/internal/observability/tracer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/types"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/ehrIndexer"
)

func (i *Index) AddEhrDoc(docType types.DocumentType, docMeta *model.DocumentMeta, privKey *[32]byte) ([]byte, error) {
	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	params := ehrIndexer.IDocsAddEhrDocParams{
		DocType:   uint8(docType),
		Id:        docMeta.Id,
		Version:   docMeta.Version,
		Timestamp: docMeta.Timestamp,
		Attrs:     docMeta.Attrs,
		Signer:    userAddress,
		Deadline:  big.NewInt(time.Now().Add(i.txTimeout).Unix()),
		Signature: make([]byte, signatureLength),
	}

	data, err := i.ehrIndexAbi.Pack("addEhrDoc", params)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack1 error: %w", err)
	}

	params.Signature, err = makeSignature(data, userKey, params.Deadline)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err = i.ehrIndexAbi.Pack("addEhrDoc", params)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack2 error: %w", err)
	}

	return data, nil
}

func (i *Index) GetDocLastByType(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType) (*model.DocumentMeta, error) {
	ctx, span := tracer.Start(ctx, "indexer.GetDocLastByType", trace.WithAttributes(
		attribute.String("ehrUUID", ehrUUID.String()),
		attribute.String("docType", docType.String()),
	))
	defer span.End()

	eID := [32]byte{}
	copy(eID[:], ehrUUID[:])

	docMeta, err := i.ehrIndex.GetLastEhrDocByType(&bind.CallOpts{Context: ctx}, eID, uint8(docType))
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("ehrIndex.GetLastEhrDocByType error: %w ehrUUID %s docType %s", err, ehrUUID.String(), docType.String())
	}

	return (*model.DocumentMeta)(&docMeta), nil
}

func (i *Index) ListDocByType(ctx context.Context, userID, systemID string, docType types.DocumentType) ([]model.DocumentMeta, error) {
	ctx, span := tracer.Start(ctx, "indexer.ListDocByType", trace.WithAttributes(
		attribute.String("userID", userID),
		attribute.String("systemID", systemID),
		attribute.String("docType", docType.String()),
	))
	defer span.End()

	IDHash := sha3.Sum256([]byte(userID + systemID))

	docsMeta, err := i.ehrIndex.GetEhrDocs(&bind.CallOpts{Context: ctx}, IDHash, uint8(docType))
	if err != nil {
		if errors.Is(err, errors.ErrNotFound) {
			return nil, err
		}
		return nil, fmt.Errorf("ehrIndex.GetEhrDocs error: %w IDHash %x docType %s", err, IDHash, docType.String())
	}

	var list []model.DocumentMeta

	for _, dm := range docsMeta {
		list = append(list, model.DocumentMeta(dm))
	}

	return list, nil
}

func (i *Index) GetDocLastByBaseID(ctx context.Context, userID, systemID string, docType types.DocumentType, UIDHash *[32]byte) (*model.DocumentMeta, error) {
	ctx, span := tracer.Start(ctx, "indexer.GetDocLastByBaseID")
	defer span.End()

	IDHash := sha3.Sum256([]byte(userID + systemID))

	docMeta, err := i.ehrIndex.GetDocLastByBaseID(&bind.CallOpts{Context: ctx}, IDHash, uint8(docType), *UIDHash)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return nil, fmt.Errorf("ehrIndex.GetDocLastByBaseID error: %w", errors.ErrNotFound)
		}
		return nil, fmt.Errorf("ehrIndex.GetDocLastByBaseID error: %w userID %s docType %s docBaseUIDHash %x", err, userID, docType.String(), UIDHash)
	}

	return (*model.DocumentMeta)(&docMeta), nil
}

func (i *Index) GetDocByTime(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType, timestamp uint32) (*model.DocumentMeta, error) {
	ctx, span := tracer.Start(ctx, "indexer.GetDocByTime", trace.WithAttributes(
		attribute.String("ehrUUID", ehrUUID.String()),
		attribute.String("docType", docType.String()),
	))
	defer span.End()

	eID := [32]byte{}
	copy(eID[:], ehrUUID[:])

	docMeta, err := i.ehrIndex.GetDocByTime(&bind.CallOpts{Context: ctx}, eID, uint8(docType), timestamp)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return nil, fmt.Errorf("ehrIndex.GetDocByTime error: %w", errors.ErrNotFound)
		}
		return nil, fmt.Errorf("ehrIndex.GetDocByTime error: %w ehrUUID %s docType %s timestamp %d", err, ehrUUID.String(), docType.String(), timestamp)
	}

	return (*model.DocumentMeta)(&docMeta), nil
}

func (i *Index) GetDocByVersion(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType, docBaseUIDHash, version *[32]byte) (*model.DocumentMeta, error) {
	ctx, span := tracer.Start(ctx, "indexer.GetDocByVersion", trace.WithAttributes(
		attribute.String("ehrUUID", ehrUUID.String()),
	))
	defer span.End()

	eID := [32]byte{}
	copy(eID[:], ehrUUID[:])

	docMeta, err := i.ehrIndex.GetDocByVersion(&bind.CallOpts{Context: ctx}, eID, uint8(docType), *docBaseUIDHash, *version)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("ehrIndex.GetDocByVersion error: %w ehrUUID %s docType %s docBaseUIDHash %x version %s", err, ehrUUID.String(), docType.String(), docBaseUIDHash, version)
	}

	return (*model.DocumentMeta)(&docMeta), nil
}

func (i *Index) SetEhrSubject(ctx context.Context, ehrUUID *uuid.UUID, subjectID, subjectNamespace string, privKey *[32]byte) ([]byte, error) {
	_, span := tracer.Start(ctx, "indexer.SetEhrSubject", trace.WithAttributes(
		attribute.String("ehrUUID", ehrUUID.String()),
		attribute.String("subjectID", subjectID),
	))
	defer span.End()

	eID := [32]byte{}
	copy(eID[:], ehrUUID[:])

	subjectKey := sha3.Sum256([]byte(subjectID + subjectNamespace))

	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	deadline := big.NewInt(time.Now().Add(i.txTimeout).Unix())

	sig := make([]byte, signatureLength)

	data, err := i.ehrIndexAbi.Pack("setEhrSubject", subjectKey, eID, userAddress, deadline, sig)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack1 error: %w", err)
	}

	sig, err = makeSignature(data, userKey, deadline)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err = i.ehrIndexAbi.Pack("setEhrSubject", subjectKey, eID, userAddress, deadline, sig)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack2 error: %w", err)
	}

	return data, nil
}

func (i *Index) GetEhrUUIDBySubject(ctx context.Context, subjectID, subjectNamespace string) (*uuid.UUID, error) {
	_, span := tracer.Start(ctx, "indexer.GetEhrUUIDBySubject", trace.WithAttributes(
		attribute.String("subjectID", subjectID),
		attribute.String("subjectNamespace", subjectNamespace),
	))
	defer span.End()

	subjectKey := sha3.Sum256([]byte(subjectID + subjectNamespace))

	ehrUUIDRaw, err := i.ehrIndex.EhrSubject(&bind.CallOpts{Context: ctx}, subjectKey)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.EhrSubjec error: %w", err)
	}

	ehrUUID, err := uuid.FromBytes(ehrUUIDRaw[:16])
	if err != nil {
		return nil, fmt.Errorf("ehrUUID FromBytes error: %w ehrUUIDRaw %x", err, ehrUUIDRaw)
	}

	return &ehrUUID, nil
}

func (i *Index) DeleteDoc(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType, docBaseUIDHash, version, privKey *[32]byte) (string, uint64, error) {
	_, span := tracer.Start(ctx, "indexer.DeleteDoc", trace.WithAttributes(
		attribute.String("ehrUUID", ehrUUID.String()),
		attribute.String("docType", docType.String()),
	))
	defer span.End()

	var eID [32]byte

	copy(eID[:], ehrUUID[:])

	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return "", 0, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	i.Lock()
	defer i.Unlock()

	sig := make([]byte, signatureLength)

	deadline := big.NewInt(time.Now().Add(i.txTimeout).Unix())

	data, err := i.ehrIndexAbi.Pack("deleteDoc", eID, uint8(docType), *docBaseUIDHash, version, userAddress, deadline, sig)
	if err != nil {
		return "", 0, fmt.Errorf("abi.Pack error: %w", err)
	}

	sig, err = makeSignature(data, userKey, deadline)
	if err != nil {
		return "", 0, fmt.Errorf("makeSignature error: %w", err)
	}

	tx, err := i.ehrIndex.DeleteDoc(i.GetNewOpts(i.transactOpts), eID, uint8(docType), *docBaseUIDHash, *version, userAddress, deadline, sig)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return "", 0, errors.ErrNotFound
		} else if strings.Contains(err.Error(), "ADL") {
			return "", 0, errors.ErrAlreadyDeleted
		}
		return "", 0, fmt.Errorf("ehrIndex.DeleteDoc error: %w ehrUUID %s docType %s", err, ehrUUID.String(), docType.String())
	}

	return tx.Hash().Hex(), tx.Nonce(), nil
}
