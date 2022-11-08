package indexer

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"hms/gateway/pkg/errors"
)

func (i *Index) UserNew(ctx context.Context, userID string, systemID string, role uint8, pwdHash []byte, privKey *[32]byte, nonce *big.Int) (string, error) {
	i.Lock()
	defer i.Unlock()

	var uID, sID [32]byte

	copy(uID[:], []byte(userID)[:])
	copy(sID[:], []byte(systemID)[:])

	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return "", fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	if nonce == nil {
		nonce, err = i.userNonce(ctx, &userAddress)
		if err != nil {
			return "", fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	sig, err := makeSignature(
		userKey,
		abi.Arguments{{Type: String}, {Type: Address}, {Type: Bytes32}, {Type: Bytes32}, {Type: Uint256}, {Type: Bytes}, {Type: Uint256}},
		"userAdd", userAddress, uID, sID, big.NewInt(int64(role)), pwdHash, nonce,
	)
	if err != nil {
		return "", fmt.Errorf("makeSignature error: %w", err)
	}

	//TODO remove userAddr arg, its same as signer
	tx, err := i.ehrIndex.UserNew(i.transactOpts, userAddress, uID, sID, role, pwdHash, nonce, userAddress, sig)
	if err != nil {
		switch err.Error() {
		case ExecutionRevertedAEX:
			return "", errors.ErrAlreadyExist
		default:
			return "", fmt.Errorf("ehrIndex.UserAdd error: %w", err)
		}
	}

	return tx.Hash().Hex(), nil
}

func (i *Index) GetUserPasswordHash(ctx context.Context, userAddr common.Address) ([]byte, error) {
	userPasswordHash, err := i.ehrIndex.GetUserPasswordHash(&bind.CallOpts{Context: ctx}, userAddr)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("ehrIndex.GetUserPasswordHash error: %w userAddr %s", err, userAddr.String())
	}

	return userPasswordHash, nil
}
