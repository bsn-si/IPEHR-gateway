package indexer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
)

func (i *Index) DataUpdate(ctx context.Context, groupID, dataID, ehrID *uuid.UUID, data []byte) (string, error) {
	var gID, dID, eID [32]byte

	copy(gID[:], groupID[:])
	copy(dID[:], dataID[:])
	copy(eID[:], ehrID[:])

	deadline := big.NewInt(time.Now().Add(i.txTimeout).Unix())

	packed, err := i.dataStoreAbi.Pack("dataUpdate", gID, dID, eID, data, i.signerAddress, deadline, make([]byte, signatureLength))
	if err != nil {
		return "", fmt.Errorf("abi.Pack1 error: %w", err)
	}

	signature, err := makeSignature(packed, i.signerKey, deadline)
	if err != nil {
		return "", fmt.Errorf("makeSignature error: %w", err)
	}

	tx, err := i.dataStore.DataUpdate(i.transactOpts, gID, dID, eID, data, i.signerAddress, deadline, signature)
	if err != nil {
		return "", fmt.Errorf("dataStore.DataUpdate error: %w", err)
	}

	return tx.Hash().String(), nil
}
