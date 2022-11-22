package indexer

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"
)

const signatureLength = 65

func makeSignature(data []byte, nonce *big.Int, pk *ecdsa.PrivateKey) ([]byte, error) {
	data = data[:len(data)-97]

	paddingLength := 32 - (len(data) % 32)
	data = append(data, make([]byte, paddingLength)...)

	hash := crypto.Keccak256Hash(data)

	nonceBytes, _ := abi.Arguments{{Type: Uint256}}.Pack(nonce)

	prefixedHash := crypto.Keccak256Hash(
		[]byte("\x19Ethereum Signed Message:\n32"),
		hash.Bytes(),
		nonceBytes,
	)

	sig, err := crypto.Sign(prefixedHash.Bytes(), pk)
	if err != nil {
		return nil, fmt.Errorf("crypto.Sign error: %w", err)
	}

	sig[64] += 27

	return sig, nil
}
