package indexer

import (
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer/ehrIndexer"
	userModel "hms/gateway/pkg/user/model"
	"hms/gateway/pkg/user/roles"
)

func (i *Index) UserNew(ctx context.Context, userID, systemID string, role uint8, pwdHash, content []byte, userPrivKey *[32]byte, nonce *big.Int) (string, error) {
	i.Lock()
	defer i.Unlock()

	IDHash := sha3.Sum256([]byte(userID + systemID))

	userKey, err := crypto.ToECDSA(userPrivKey[:])
	if err != nil {
		return "", fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	if nonce == nil {
		nonce, err = i.userNonce(ctx, &i.signerAddress)
		if err != nil {
			return "", fmt.Errorf("signerNonce error: %w address: %s", err, i.signerAddress.String())
		}
	}

	var attrs []ehrIndexer.AttributesAttribute

	switch roles.Role(role) {
	case roles.Patient:
		attrs = []ehrIndexer.AttributesAttribute{
			{Code: model.AttributePasswordHash, Value: pwdHash},
		}
	case roles.Doctor:
		attrs = []ehrIndexer.AttributesAttribute{
			{Code: model.AttributePasswordHash, Value: pwdHash},
			{Code: model.AttributeContent, Value: content},
		}
	default:
		return "", errors.ErrFieldIsIncorrect("role")
	}

	sig := make([]byte, signatureLength)

	data, err := i.abi.Pack("userNew", userAddress, IDHash, role, attrs, i.signerAddress, sig)
	if err != nil {
		return "", fmt.Errorf("abi.Pack error: %w", err)
	}

	sig, err = makeSignature(data, nonce, i.signerKey)
	if err != nil {
		return "", fmt.Errorf("makeSignature error: %w", err)
	}

	tx, err := i.ehrIndex.UserNew(i.transactOpts, userAddress, IDHash, role, attrs, i.signerAddress, sig)
	if err != nil {
		if strings.Contains(err.Error(), "AEX") {
			return "", errors.ErrAlreadyExist
		}

		return "", fmt.Errorf("ehrIndex.UserNew error: %w", err)
	}

	return tx.Hash().Hex(), nil
}

func (i *Index) GetUserPasswordHash(ctx context.Context, userAddr common.Address) ([]byte, error) {
	user, err := i.ehrIndex.GetUser(&bind.CallOpts{Context: ctx}, userAddr)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.Users error: %w userAddr %s", err, userAddr.String())
	}

	if user.IDHash == [32]byte{} {
		return nil, errors.ErrNotFound
	}

	pwdHash := model.Attributes(user.Attrs).GetByCode(model.AttributePasswordHash)
	if pwdHash == nil {
		return nil, errors.ErrFieldIsEmpty("PasswordHash")
	}

	return pwdHash, nil
}

func (i *Index) GetUser(ctx context.Context, userAddr common.Address) (*userModel.User, error) {
	user, err := i.ehrIndex.GetUser(&bind.CallOpts{Context: ctx}, userAddr)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.Users error: %w userAddr %s", err, userAddr.String())
	}

	if user.IDHash == [32]byte{} {
		return nil, errors.ErrNotFound
	}

	return &user, nil
}
