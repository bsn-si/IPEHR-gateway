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

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/types"
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

func New(contractAddr, keyPath string, client *ethclient.Client) *Index {
	ctx := context.Background()

	key, err := os.ReadFile(keyPath)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimSpace(string(key)))
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

func (i *Index) AddEhrDoc(ehrUUID *uuid.UUID, docMeta *model.DocumentMeta) (string, error) {
	tx, err := i.ehrIndex.AddEhrDoc(
		i.transactOpts,
		new(big.Int).SetBytes(ehrUUID[:]),
		ehrIndexer.EhrIndexerDocumentMeta{
			DocType:        uint8(docMeta.TypeCode),
			Status:         uint8(docMeta.Status),
			StorageId:      new(big.Int).SetBytes(docMeta.CID[:]),
			DocIdEncrypted: docMeta.DocUIDEncrypted,
			Timestamp:      uint32(docMeta.Timestamp),
		},
	)
	if err != nil {
		return "", fmt.Errorf("ehrIndex.SetEhrDoc error: %w", err)
	}

	log.Printf("SetEhrDoc tx %s nonce %d", tx.Hash().Hex(), tx.Nonce())

	return tx.Hash().Hex(), nil
}

func (i *Index) GetDocLastByType(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType) (*model.DocumentMeta, error) {
	// TODO
	return &model.DocumentMeta{}, nil
}

func (i *Index) GetDocLastByBaseID(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType, docBaseUIDHash *[32]byte) (*model.DocumentMeta, error) {
	// TODO
	return &model.DocumentMeta{}, nil
}

func (i *Index) GetDocByTime(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType, timestamp uint32) (*model.DocumentMeta, error) {
	// TODO
	return &model.DocumentMeta{}, nil
}

func (i *Index) GetDocByVersion(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType, docBaseUIDHash *[32]byte, version string) (*model.DocumentMeta, error) {
	// TODO
	return &model.DocumentMeta{}, nil
}

func (i *Index) SetDocKeyEncrypted(key *[32]byte, value []byte) (string, error) {
	tx, err := i.ehrIndex.SetDocAccess(
		i.transactOpts,
		new(big.Int).SetBytes(key[:]),
		value,
	)
	if err != nil {
		return "", fmt.Errorf("ehrIndex.SetDocAccess error: %w", err)
	}

	log.Printf("SetDocAccess tx %s nonce %d", tx.Hash().Hex(), tx.Nonce())

	return tx.Hash().Hex(), nil
}

func (i *Index) GetDocKeyEncrypted(ctx context.Context, userID string, cidBytes *[32]byte) ([]byte, error) {
	docAccessIndexKey := sha3.Sum256(append(cidBytes[:], []byte(userID)...))

	callOpts := &bind.CallOpts{
		Context: ctx,
	}

	docAccessValue, err := i.ehrIndex.DocAccess(callOpts, new(big.Int).SetBytes(docAccessIndexKey[:]))
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.DocAccess error: %w", err)
	}

	return docAccessValue, nil
}

func (i *Index) SetGroupAccess(ctx context.Context, key *[32]byte, value []byte) (string, error) {
	// TODO переименовать в контракте DataAccess -> GroupAccess
	tx, err := i.ehrIndex.SetDataAccess(
		i.transactOpts,
		new(big.Int).SetBytes(key[:]),
		value,
	)
	if err != nil {
		return "", fmt.Errorf("ehrIndex.SetDataAccess error: %w", err)
	}

	log.Printf("SetDataAccess tx %s nonce %d", tx.Hash().Hex(), tx.Nonce())

	return tx.Hash().Hex(), nil
}

func (i *Index) GetGroupAccess(ctx context.Context, userID string, groupUUID *uuid.UUID) ([]byte, error) {
	dataAccessIndexKey := sha3.Sum256(append([]byte(userID), groupUUID[:]...))

	callOpts := &bind.CallOpts{
		Context: ctx,
	}

	dataAccessValue, err := i.ehrIndex.DataAccess(callOpts, new(big.Int).SetBytes(dataAccessIndexKey[:]))
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.DataAccess error: %w", err)
	}

	log.Printf("dataAccessValue: %x", dataAccessValue)

	if len(dataAccessValue) == 0 {
		return nil, errors.ErrIsNotExist
	}

	return dataAccessValue, nil
}

func (i *Index) SetSubject(ctx context.Context, ehrUUID *uuid.UUID, subjectID, subjectNamespace string) (string, error) {
	subjectKey := sha3.Sum256([]byte(subjectID + subjectNamespace))

	tx, err := i.ehrIndex.SetEhrSubject(
		i.transactOpts,
		new(big.Int).SetBytes(subjectKey[:]),
		new(big.Int).SetBytes(ehrUUID[:]),
	)
	if err != nil {
		return "", fmt.Errorf("ehrIndex.SetSubject error: %w", err)
	}

	log.Printf("SetSubject tx %s nonce %d", tx.Hash().Hex(), tx.Nonce())

	return tx.Hash().Hex(), nil
}

func (i *Index) GetEhrUUIDBySubject(ctx context.Context, subjectID, subjectNamespace string) (*uuid.UUID, error) {
	subjectKey := sha3.Sum256([]byte(subjectID + subjectNamespace))

	callOpts := &bind.CallOpts{
		Context: ctx,
	}

	ehrUUIDRaw, err := i.ehrIndex.EhrSubject(
		callOpts,
		new(big.Int).SetBytes(subjectKey[:]),
	)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.EhrSubjec error: %w", err)
	}

	ehrUUID, err := uuid.ParseBytes(ehrUUIDRaw.Bytes())
	if err != nil {
		return nil, fmt.Errorf("ehrUUID ParseBytes error: %w ehrUUIDRaw %x", err, ehrUUIDRaw.Bytes())
	}

	return &ehrUUID, nil
}

func (i *Index) DeleteDoc(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType, docBaseUIDHash *[32]byte, version string) (string, error) {
	// TODO
	return "", nil
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
