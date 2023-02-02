package indexer

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/users"
	userModel "github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/user/roles"
)

func (i *Index) UserNew(ctx context.Context, userID, systemID string, role uint8, pwdHash, content []byte, userPrivKey *[32]byte, nonce *big.Int) ([]byte, error) {
	i.Lock()
	defer i.Unlock()

	IDHash := sha3.Sum256([]byte(userID + systemID))

	userKey, err := crypto.ToECDSA(userPrivKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	if nonce == nil {
		nonce, err = i.Nonce(ctx, i.users, &i.signerAddress)
		if err != nil {
			return nil, fmt.Errorf("signerNonce error: %w address: %s", err, i.signerAddress.String())
		}
	}

	var attrs []users.AttributesAttribute

	switch roles.Role(role) {
	case roles.Patient:
		attrs = []users.AttributesAttribute{
			{Code: model.AttributePasswordHash, Value: pwdHash},
		}
	case roles.Doctor:
		attrs = []users.AttributesAttribute{
			{Code: model.AttributePasswordHash, Value: pwdHash},
			{Code: model.AttributeContent, Value: content},
		}
	default:
		return nil, errors.ErrFieldIsIncorrect("role")
	}

	data, err := i.usersAbi.Pack("userNew", userAddress, IDHash, role, attrs, i.signerAddress, make([]byte, signatureLength))
	if err != nil {
		return nil, fmt.Errorf("abi.Pack1 error: %w", err)
	}

	signature, err := makeSignature(data, nonce, i.signerKey)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err = i.usersAbi.Pack("userNew", userAddress, IDHash, role, attrs, i.signerAddress, signature)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack2 error: %w", err)
	}

	return data, nil
}

func (i *Index) GetUserPasswordHash(ctx context.Context, userAddr common.Address) ([]byte, error) {
	user, err := i.users.GetUser(&bind.CallOpts{Context: ctx}, userAddr)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.Users error: %w userAddr %s", err, userAddr.String())
	}

	if user.IDHash == [32]byte{} {
		return nil, errors.ErrNotFound
	}

	pwdHash := model.AttributesUsers(user.Attrs).GetByCode(model.AttributePasswordHash)
	if pwdHash == nil {
		return nil, errors.ErrFieldIsEmpty("PasswordHash")
	}

	return pwdHash, nil
}

func (i *Index) GetUser(ctx context.Context, userAddr common.Address) (*userModel.User, error) {
	user, err := i.users.GetUser(&bind.CallOpts{Context: ctx}, userAddr)
	if err != nil {
		return nil, fmt.Errorf("users.GetUser error: %w userAddr %s", err, userAddr.String())
	}

	if user.IDHash == [32]byte{} {
		return nil, errors.ErrNotFound
	}

	return &user, nil
}

func (i *Index) GetUserByCode(ctx context.Context, code uint64) (*userModel.User, error) {
	user, err := i.users.GetUserByCode(&bind.CallOpts{Context: ctx}, code)
	if err != nil {
		return nil, fmt.Errorf("users.GetUserByCode error: %w code %d", err, code)
	}

	if user.IDHash == [32]byte{} {
		return nil, errors.ErrNotFound
	}

	return &user, nil
}
