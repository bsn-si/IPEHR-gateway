package indexer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"sync"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"golang.org/x/crypto/sha3"

	"github.com/bsn-si/IPEHR-gateway/src/pkg/access"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/errors"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/accessStore"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/dataStore"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/ehrIndexer"
	"github.com/bsn-si/IPEHR-gateway/src/pkg/indexer/users"
)

type Index struct {
	sync.RWMutex
	client        *ethclient.Client
	transactOpts  *bind.TransactOpts
	signerKey     *ecdsa.PrivateKey
	signerAddress common.Address

	ehrIndex    *ehrIndexer.EhrIndexer
	accessStore *accessStore.AccessStore
	users       *users.Users
	dataStore   *dataStore.DataStore

	ehrIndexAbi  *abi.ABI
	usersAbi     *abi.ABI
	dataStoreAbi *abi.ABI
}

const (
	ExecutionRevertedNFD = "execution reverted: NFD"
	ExecutionRevertedDNY = "execution reverted: DNY"
	ExecutionRevertedAEX = "execution reverted: AEX"
)

var (
	String, _  = abi.NewType("string", "", nil)
	Bytes32, _ = abi.NewType("bytes32", "", nil)
	Bytes, _   = abi.NewType("bytes", "", nil)
	Uint8, _   = abi.NewType("uint8", "", nil)
	Uint256, _ = abi.NewType("uint256", "", nil)
	Address, _ = abi.NewType("address", "", nil)
	Access, _  = abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "IdHash", Type: "bytes32"},
		{Name: "IdEncr", Type: "bytes"},
		{Name: "KeyEncr", Type: "bytes"},
		{Name: "Level", Type: "uint8"},
	})
	DocMeta, _ = abi.NewType("tuple", "", []abi.ArgumentMarshaling{
		{Name: "docType", Type: "uint8"},
		{Name: "status", Type: "uint8"},
		{Name: "CID", Type: "bytes"},
		{Name: "dealCID", Type: "bytes"},
		{Name: "minerAddress", Type: "bytes"},
		{Name: "docUIDEncrypted", Type: "bytes"},
		{Name: "docBaseUIDHash", Type: "bytes32"},
		{Name: "version", Type: "bytes32"},
		{Name: "isLast", Type: "bool"},
		{Name: "timestamp", Type: "uint32"},
	})
)

func New(ehrIndexAddr, accessStoreAddr, usersAddr, dataStoreAddr, keyPath string, client *ethclient.Client, gasTipCap int64) *Index {
	ctx := context.Background()

	key, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimSpace(string(key)))
	if err != nil {
		log.Fatal(err)
	}

	signerAddress := crypto.PubkeyToAddress(privateKey.PublicKey)

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal(err)
	}

	switch {
	case !common.IsHexAddress(accessStoreAddr):
		log.Fatal("ehrIndex contract address is incorrect")
	case !common.IsHexAddress(accessStoreAddr):
		log.Fatal("accessStore contract address is incorrect")
	case !common.IsHexAddress(usersAddr):
		log.Fatal("users contract address is incorrect")
	case !common.IsHexAddress(dataStoreAddr):
		log.Fatal("dataStore contract address is incorrect")
	}

	ehrIndex, err := ehrIndexer.NewEhrIndexer(common.HexToAddress(ehrIndexAddr), client)
	if err != nil {
		log.Fatal(err)
	}

	accessStore, err := accessStore.NewAccessStore(common.HexToAddress(accessStoreAddr), client)
	if err != nil {
		log.Fatal(err)
	}

	_users, err := users.NewUsers(common.HexToAddress(usersAddr), client)
	if err != nil {
		log.Fatal(err)
	}

	_dataStore, err := dataStore.NewDataStore(common.HexToAddress(dataStoreAddr), client)
	if err != nil {
		log.Fatal(err)
	}

	ehrIndexAbi, _ := ehrIndexer.EhrIndexerMetaData.GetAbi()
	usersAbi, _ := users.UsersMetaData.GetAbi()
	dataStoreAbi, _ := dataStore.DataStoreMetaData.GetAbi()

	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	if gasTipCap > 0 {
		transactOpts.GasTipCap = big.NewInt(gasTipCap)
	}

	return &Index{
		client:        client,
		transactOpts:  transactOpts,
		signerKey:     privateKey,
		signerAddress: signerAddress,

		ehrIndex:    ehrIndex,
		accessStore: accessStore,
		users:       _users,
		dataStore:   _dataStore,

		ehrIndexAbi:  ehrIndexAbi,
		usersAbi:     usersAbi,
		dataStoreAbi: dataStoreAbi,
	}
}

func (i *Index) SetEhrUser(ctx context.Context, userID, systemID string, ehrUUID *uuid.UUID, privKey *[32]byte, nonce *big.Int) ([]byte, error) {
	var eID [32]byte

	copy(eID[:], ehrUUID[:])

	IDHash := sha3.Sum256([]byte(userID + systemID))

	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	if nonce == nil {
		nonce, err = i.Nonce(ctx, i.users, &userAddress)
		if err != nil {
			return nil, fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	sig := make([]byte, signatureLength)

	data, err := i.ehrIndexAbi.Pack("setEhrUser", IDHash, eID, userAddress, sig)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack1 error: %w", err)
	}

	sig, err = makeSignature(data, nonce, userKey)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err = i.ehrIndexAbi.Pack("setEhrUser", IDHash, eID, userAddress, sig)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack2 error: %w", err)
	}

	return data, err
}

func (i *Index) GetEhrUUIDByUserID(ctx context.Context, userID, systemID string) (*uuid.UUID, error) {
	IDHash := sha3.Sum256([]byte(userID + systemID))

	ehrUUIDRaw, err := i.ehrIndex.GetEhrUser(&bind.CallOpts{Context: ctx}, IDHash)
	if err != nil {
		return nil, fmt.Errorf("EhrUsers get error: %w userID %s systemID %s", err, userID, systemID)
	}

	if ehrUUIDRaw == [32]byte{} {
		return nil, errors.ErrIsNotExist
	}

	ehrUUID, err := uuid.FromBytes(ehrUUIDRaw[:16])
	if err != nil {
		return nil, fmt.Errorf("EhrUsers parse UUID error: %w userID %s ehrUUIDRaw %x", err, userID, ehrUUIDRaw)
	}

	return &ehrUUID, nil
}

func (i *Index) GetDocKeyEncrypted(ctx context.Context, userID, systemID string, CID []byte) ([]byte, error) {
	IDHash := sha3.Sum256([]byte(userID + systemID))

	data, err := abi.Arguments{{Type: Bytes32}, {Type: Uint8}}.Pack(IDHash, access.Doc)
	if err != nil {
		return nil, fmt.Errorf("args.Pack error: %w", err)
	}

	accessID := crypto.Keccak256Hash(data)
	CIDHash := crypto.Keccak256Hash(CID)

	callOpts := &bind.CallOpts{
		Context: ctx,
	}

	accessObj, err := i.accessStore.GetAccessByIdHash(callOpts, accessID, CIDHash)
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
			return nil, errors.ErrNotFound
		}

		return nil, fmt.Errorf("ehrIndex.DocAccess error: %w", err)
	}

	return accessObj.KeyEncr, nil
}

/*
func (i *Index) SetGroupAccess(ctx context.Context, key *[32]byte, value []byte, accessLevel uint8, privKey *[32]byte, nonce *big.Int) (string, error) {
	i.Lock()
	defer i.Unlock()

	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return "", fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	access := ehrIndexer.EhrAccessAccess{
		Level:        accessLevel,
		KeyEncrypted: value,
	}

	if nonce == nil {
		nonce, err = i.userNonce(ctx, &userAddress)
		if err != nil {
			return "", fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	sig, err := makeSignature(
		userKey,
		abi.Arguments{{Type: String}, {Type: Bytes32}, {Type: Access}, {Type: Uint256}},
		"setGroupAccess", *key, access, nonce,
	)
	if err != nil {
		return "", fmt.Errorf("makeSignature error: %w", err)
	}

	tx, err := i.ehrIndex.SetGroupAccess(i.transactOpts, *key, access, nonce, userAddress, sig)
	if err != nil {
		return "", fmt.Errorf("ehrIndex.SetGroupAccess error: %w", err)
	}

	return tx.Hash().Hex(), nil
}

func (i *Index) GetGroupAccess(ctx context.Context, userID string, groupUUID *uuid.UUID) ([]byte, error) {
	groupAccessIndexKey := sha3.Sum256(append([]byte(userID), groupUUID[:]...))

	access, err := i.ehrIndex.AccessStore(&bind.CallOpts{Context: ctx}, groupAccessIndexKey)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.GroupAccess error: %w", err)
	}

	if len(access.KeyEncrypted) == 0 {
		return nil, errors.ErrIsNotExist
	}

	return access.KeyEncrypted, nil
}
*/

func (i *Index) SetAllowed(ctx context.Context, address string) (string, error) {
	i.Lock()
	defer i.Unlock()

	tx, err := i.ehrIndex.SetAllowed(i.transactOpts, common.HexToAddress(address), true)
	if err != nil {
		return "", fmt.Errorf("ehrIndex.SetAllowed error: %w", err)
	}

	return tx.Hash().Hex(), nil
}

/*
func Init(name string) *Index {
	if name == "" {
		log.Fatal("name is empty")
	}

	id := sha3.Sum256([]byte(name))

	stor := storage.Storage()

	data, err := stor.Get(&id)
	if err != nil && !errors.Is(err, errors.ErrIsNotExist) {
		log.Fatal(err)
	}

	var cache map[string][]byte
	if errors.Is(err, errors.ErrIsNotExist) {
		cache = make(map[string][]byte)
	} else {
		err = msgpack.Unmarshal(data, &cache)
		if err != nil {
			log.Fatal(err)
		}
	}

	return &Index{
		id:      &id,
		name:    name,
		cache:   cache,
		storage: stor,
	}
}

func (i *Index) Add(itemID string, item interface{}) (err error) {
	i.Lock()
	defer func() {
		if err != nil {
			delete(i.cache, itemID)
		}
		i.Unlock()
	}()

	if _, ok := i.cache[itemID]; ok {
		return errors.ErrAlreadyExist
	}

	data, err := msgpack.Marshal(item)
	if err != nil {
		return fmt.Errorf("item marshal error: %w", err)
	}

	i.cache[itemID] = data

	data, err = msgpack.Marshal(i.cache)
	if err != nil {
		return fmt.Errorf("cache marshal error: %w", err)
	}

	if err = i.storage.ReplaceWithID(i.id, data); err != nil {
		return fmt.Errorf("storage.ReplaceWithID error: %w", err)
	}

	return nil
}

func (i *Index) Replace(itemID string, item interface{}) (err error) {
	i.Lock()
	defer func() {
		if err != nil {
			delete(i.cache, itemID)
		}
		i.Unlock()
	}()

	data, err := msgpack.Marshal(item)
	if err != nil {
		return fmt.Errorf("item marshal error: %w", err)
	}

	i.cache[itemID] = data

	data, err = msgpack.Marshal(i.cache)
	if err != nil {
		return err
	}

	err = i.storage.ReplaceWithID(i.id, data)
	if err != nil {
		return fmt.Errorf("storage.ReplaceWithID error: %w", err)
	}

	return nil
}

func (i *Index) GetByID(itemID string, dst interface{}) error {
	i.RLock()
	item, ok := i.cache[itemID]
	i.RUnlock()

	if !ok {
		return errors.ErrIsNotExist
	}

	if err := msgpack.Unmarshal(item, dst); err != nil {
		return fmt.Errorf("item unmarshal error: %w", err)
	}

	return nil
}

func (i *Index) Delete(itemID string) error {
	i.Lock()
	defer i.Unlock()

	item, ok := i.cache[itemID]
	if !ok {
		return errors.ErrIsNotExist
	}

	delete(i.cache, itemID)

	data, err := msgpack.Marshal(i.cache)
	if err != nil {
		i.cache[itemID] = item
		return err
	}

	err = i.storage.ReplaceWithID(i.id, data)
	if err != nil {
		i.cache[itemID] = item
		return err
	}

	return nil
}
*/
