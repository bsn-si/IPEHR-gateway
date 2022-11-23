package indexer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"

	"hms/gateway/pkg/errors"
)

func (i *Index) UserNew(ctx context.Context, userID string, systemID string, role uint8, pwdHash []byte, privKey *[32]byte, nonce *big.Int) (string, error) {
	i.Lock()
	defer i.Unlock()

	var uID, sID [32]byte

	copy(uID[:], userID)
	copy(sID[:], systemID)

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

	sig := make([]byte, 65)

	data, err := i.abi.Pack("userNew", userAddress, uID, sID, role, pwdHash, userAddress, sig)
	if err != nil {
		return "", fmt.Errorf("abi.Pack error: %w", err)
	}

	sig, err = makeSignature(data, nonce, userKey)
	if err != nil {
		return "", fmt.Errorf("makeSignature error: %w", err)
	}

	//TODO remove userAddr arg, its same as signer
	tx, err := i.ehrIndex.UserNew(i.transactOpts, userAddress, uID, sID, role, pwdHash, userAddress, sig)
	if err != nil {
		switch err.Error() {
		case ExecutionRevertedAEX:
			return "", errors.ErrAlreadyExist
		default:
			return "", fmt.Errorf("ehrIndex.UserNew error: %w", err)
		}
	}

	return tx.Hash().Hex(), nil
}

func (i *Index) GetUserPasswordHash(ctx context.Context, userAddr common.Address) ([]byte, error) {
	user, err := i.ehrIndex.Users(&bind.CallOpts{Context: ctx}, userAddr)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.Users error: %w userAddr %s", err, userAddr.String())
	}

	if user.Id == [32]byte{} {
		return nil, errors.ErrNotFound
	}

	return user.PwdHash, nil
}
