package indexer

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/internal/observability/tracer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/docs/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/users"
	userModel "github.com/bsn-si/IPEHR-gateway/src/pkg/user/model"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/user/roles"
)

func (i *Index) UserNew(ctx context.Context, userID, systemID string, role uint8, pwdHash, content []byte, userPrivKey *[32]byte) ([]byte, error) {
	ctx, span := tracer.GetTracer().Start(ctx, "user_index.user_new") //nolint
	defer span.End()

	i.Lock()
	defer i.Unlock()

	IDHash := sha3.Sum256([]byte(userID + systemID))

	userKey, err := crypto.ToECDSA(userPrivKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

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

	deadline := big.NewInt(time.Now().Add(i.txTimeout).Unix())

	data, err := i.usersAbi.Pack("userNew", userAddress, IDHash, role, attrs, i.signerAddress, deadline, make([]byte, signatureLength))
	if err != nil {
		return nil, fmt.Errorf("abi.Pack1 error: %w", err)
	}

	signature, err := makeSignature(data, i.signerKey, deadline)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err = i.usersAbi.Pack("userNew", userAddress, IDHash, role, attrs, i.signerAddress, deadline, signature)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack2 error: %w", err)
	}

	return data, nil
}

func (i *Index) GetUserPasswordHash(ctx context.Context, userAddr common.Address) ([]byte, error) {
	ctx, span := tracer.GetTracer().Start(ctx, "user_index.get_user_password_hash") //nolint
	defer span.End()

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
	ctx, span := tracer.GetTracer().Start(ctx, "user_index.get_user") //nolint
	defer span.End()

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
	ctx, span := tracer.GetTracer().Start(ctx, "user_index.get_user_by_code") //nolint
	defer span.End()

	user, err := i.users.GetUserByCode(&bind.CallOpts{Context: ctx}, code)
	if err != nil {
		return nil, fmt.Errorf("users.GetUserByCode error: %w code %d", err, code)
	}

	if user.IDHash == [32]byte{} {
		return nil, errors.ErrNotFound
	}

	return &user, nil
}

func (i *Index) SetAccessWrapper(ctx context.Context, subjectIDHash *[32]byte, accessObj *AccessObject, userPrivKey *[32]byte) ([]byte, error) {
	ctx, span := tracer.GetTracer().Start(ctx, "user_index.set_access_wrapper") //nolint
	defer span.End()

	userKey, err := crypto.ToECDSA(userPrivKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	data, err := abi.Arguments{{Type: Bytes32}, {Type: Uint8}}.Pack(*subjectIDHash, accessObj.Kind)
	if err != nil {
		return nil, fmt.Errorf("args.Pack error: %w", err)
	}

	accessID := Keccak256(data)

	deadline := big.NewInt(time.Now().Add(i.txTimeout).Unix())

	data, err = i.usersAbi.Pack("setAccess", accessID, *accessObj, userAddress, deadline, make([]byte, signatureLength))
	if err != nil {
		return nil, fmt.Errorf("abi.Pack1 error: %w", err)
	}

	signature, err := makeSignature(data, userKey, deadline)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err = i.usersAbi.Pack("setAccess", accessID, *accessObj, userAddress, deadline, signature)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack2 error: %w", err)
	}

	return data, nil
}
