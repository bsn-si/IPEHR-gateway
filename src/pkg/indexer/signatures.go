package indexer

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
)

const signatureLength = 65

func makeSignature(data []byte, pk *ecdsa.PrivateKey, deadline *big.Int) ([]byte, error) {
	data = data[:len(data)-(signatureLength+32)]

	deadlineBytes, _ := abi.Arguments{{Type: Uint256}}.Pack(deadline)

	prefixedHash := crypto.Keccak256Hash(
		[]byte("\x19Ethereum Signed Message:\n32"),
		crypto.Keccak256(data),
		deadlineBytes,
	)

	sig, err := crypto.Sign(prefixedHash.Bytes(), pk)
	if err != nil {
		return nil, fmt.Errorf("crypto.Sign error: %w", err)
	}

	// https://ethereum.stackexchange.com/questions/78929/whats-the-magic-numbers-meaning-of-27-or-28-in-vrs-use-to-ecrover-the-sender
	sig[signatureLength-1] += 27

	return sig, nil
}
