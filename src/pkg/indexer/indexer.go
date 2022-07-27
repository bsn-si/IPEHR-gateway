package indexer

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer/ehrIndexer"
	"hms/gateway/pkg/storage"
)

type Index struct {
	sync.RWMutex
	id           *[32]byte
	name         string
	cache        map[string][]byte
	storage      storage.Storager
	client       *ethclient.Client
	ehrIndex     *ehrIndexer.EhrIndexer
	transactOpts *bind.TransactOpts
}

func New(contractAddr, endpoint, keyPath string) *Index {
	ctx := context.Background()

	key, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimSpace(string(key)))
	if err != nil {
		log.Fatal(err)
	}

	client, err := ethclient.Dial(endpoint)
	if err != nil {
		log.Fatal(err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal(err)
	}

	address := common.HexToAddress(contractAddr)

	ehrIndex, err := ehrIndexer.NewEhrIndexer(address, client)
	if err != nil {
		log.Fatal(err)
	}

	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	return &Index{
		client:       client,
		ehrIndex:     ehrIndex,
		transactOpts: transactOpts,
	}
}

func (i *Index) SetEhrUser(userID string, ehrUUID *uuid.UUID) (string, error) {
	tx, err := i.ehrIndex.SetEhrUser(
		i.transactOpts,
		new(big.Int).SetBytes([]byte(userID)),
		new(big.Int).SetBytes(ehrUUID[:]),
	)
	if err != nil {
		return "", fmt.Errorf("ehrIndex.SetEhrUser error: %w", err)
	}

	log.Printf("SetEhrUser tx %s nonce %d", tx.Hash().Hex(), tx.Nonce())

	return tx.Hash().Hex(), nil
}

func (i *Index) GetEhrUUIDByUserID(ctx context.Context, userID string) (*uuid.UUID, error) {
	callOpts := &bind.CallOpts{
		Context: ctx,
	}

	ehrUUIDRaw, err := i.ehrIndex.EhrUsers(callOpts, new(big.Int).SetBytes([]byte(userID)))
	if err != nil {
		return nil, fmt.Errorf("EhrUsers get error: %w userID %s", err, userID)
	}

	if len(ehrUUIDRaw.Bits()) == 0 {
		return nil, errors.ErrIsNotExist
	}

	ehrUUID, err := uuid.Parse(hex.EncodeToString(ehrUUIDRaw.Bytes()))
	if err != nil {
		return nil, fmt.Errorf("EhrUsers parse UUID error: %w userID %s ehrUUIDRaw %x", err, userID, ehrUUIDRaw)
	}

	return &ehrUUID, nil
}

func (i *Index) TxWait(ctx context.Context, hash string) (uint64, error) {
	h := common.HexToHash(hash)

	ticker := time.NewTicker(15 * time.Second)

	for {
		select {
		case <-ticker.C:
			receipt, err := i.client.TransactionReceipt(ctx, h)

			switch {
			case err != nil && !errors.Is(err, ethereum.NotFound):
				return 0, err
			case err == nil:
				return receipt.Status, nil
			default:
			}
		case <-ctx.Done():
			return 0, errors.ErrTimeout
		}
	}
}

func (i *Index) GetTxStatus(ctx context.Context, hash string) (uint64, error) {
	h := common.HexToHash(hash)

	receipt, err := i.client.TransactionReceipt(ctx, h)
	if err != nil {
		if errors.Is(err, ethereum.NotFound) {
			return 0, errors.ErrIsNotExist
		}

		return 0, fmt.Errorf("GetTxStatus error: %w hash %s", err, hash)
	}

	return receipt.Status, nil
}

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
