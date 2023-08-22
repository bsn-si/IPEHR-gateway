package indexer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/bsn-si/IPEHR-gateway/src/internal/observability/tracer"
	"github.com/google/uuid"
)

func (i *Index) DataUpdate(ctx context.Context, groupID, dataID, ehrID *uuid.UUID, data []byte) (string, uint64, error) {
	ctx, span := tracer.Start(ctx, "indexer.DataUpdate") //nolint
	defer span.End()

	var gID, dID, eID [32]byte

	copy(gID[:], groupID[:])
	copy(dID[:], dataID[:])
	copy(eID[:], ehrID[:])

	deadline := big.NewInt(time.Now().Add(i.txTimeout).Unix())

	packed, err := i.dataStoreAbi.Pack("dataUpdate", gID, dID, eID, data, i.signerAddress, deadline, make([]byte, signatureLength))
	if err != nil {
		return "", 0, fmt.Errorf("abi.Pack1 error: %w", err)
	}

	signature, err := makeSignature(packed, i.signerKey, deadline)
	if err != nil {
		return "", 0, fmt.Errorf("makeSignature error: %w", err)
	}

	tx, err := i.dataStore.DataUpdate(i.noncer.GetNewOpts(i.transactOpts), gID, dID, eID, data, i.signerAddress, deadline, signature)
	if err != nil {
		return "", 0, fmt.Errorf("dataStore.DataUpdate error: %w", err)
	}

	return tx.Hash().String(), tx.Nonce(), nil
}
