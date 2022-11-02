package indexer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/crypto"

	"hms/gateway/pkg/indexer/ehrIndexer"
)

func (i *Index) SetDocAccess(ctx context.Context, accessID *[32]byte, CID, keyEncrypted []byte, accessLevel uint8, privateKey *[32]byte, nonce *big.Int) ([]byte, error) {
	userKey, err := crypto.ToECDSA(privateKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	accessObj := ehrIndexer.EhrAccessAccess{
		Level:        accessLevel,
		KeyEncrypted: keyEncrypted,
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	if nonce == nil {
		nonce, err = i.userNonce(ctx, &userAddress)
		if err != nil {
			return nil, fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	sig, err := makeSignature(
		userKey,
		abi.Arguments{{Type: String}, {Type: Bytes32}, {Type: Bytes}, {Type: Access}, {Type: Uint256}},
		"setDocAccess", *accessID, CID, accessObj, nonce,
	)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err := i.abi.Pack("setDocAccess", *accessID, CID, accessObj, nonce, userAddress, sig)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack error: %w", err)
	}

	return data, nil
}
