package indexer

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ipfs/go-cid"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer/ehrIndexer"
	"hms/gateway/pkg/storage"
)

const ExecutionRevertedNFD = "execution reverted: NFD"

type Index struct {
	sync.RWMutex
	id           *[32]byte
	name         string
	cache        map[string][]byte
	storage      storage.Storager
	client       *ethclient.Client
	ehrIndex     *ehrIndexer.EhrIndexer
	transactOpts *bind.TransactOpts
	abi          *abi.ABI
}

type MultiCallTx struct {
	index *Index
	kinds []uint8
	data  [][]byte
}

func (i *Index) MultiCallTxNew() *MultiCallTx {
	return &MultiCallTx{index: i}
}

func (m *MultiCallTx) Add(kind uint8, packed []byte) {
	m.kinds = append(m.kinds, kind)
	m.data = append(m.data, packed)
}

func (m *MultiCallTx) GetTxKinds() []uint8 {
	return m.kinds
}

func (m *MultiCallTx) Commit() (string, error) {
	if len(m.data) == 0 {
		return "", fmt.Errorf("%w MultiCallTx data is empty", errors.ErrCustom)
	}

	tx, err := m.index.ehrIndex.Multicall(m.index.transactOpts, m.data)
	if err != nil {
		return "", fmt.Errorf("ehrIndex.Multicall error: %w", err)
	}

	return tx.Hash().Hex(), nil
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

	ehrIndex, err := ehrIndexer.NewEhrIndexer(address, client) // shoulbe interface
	if err != nil {
		log.Fatal(err)
	}

	bcAbi, _ := ehrIndexer.EhrIndexerMetaData.GetAbi()

	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal(err)
	}

	return &Index{
		client:       client,
		ehrIndex:     ehrIndex,
		transactOpts: transactOpts,
		abi:          bcAbi,
	}
}

func (i *Index) pack(name string, args ...interface{}) (result []byte, err error) {
	result, err = i.abi.Pack(name, args...)
	if err != nil {
		return nil, fmt.Errorf("add setEhrUser error: %w", err)
	}
	return
}

func (i *Index) SetEhrUser(userID string, ehrUUID *uuid.UUID) (packed []byte, err error) {
	var uID, eID [32]byte

	copy(uID[:], []byte(userID))
	copy(eID[:], ehrUUID[:])

	packed, err = i.pack("setEhrUser", uID, eID)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.SetEhrUser error: %w", err)
	}

	return packed, err
}

func (i *Index) GetEhrUUIDByUserID(ctx context.Context, userID string) (*uuid.UUID, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		uID      [32]byte
	)

	copy(uID[:], []byte(userID))

	ehrUUIDRaw, err := i.ehrIndex.EhrUsers(callOpts, uID)
	if err != nil {
		return nil, fmt.Errorf("EhrUsers get error: %w userID %s", err, userID)
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

func (i *Index) AddEhrDoc(ehrUUID *uuid.UUID, docMeta *model.DocumentMeta) (packed []byte, err error) {
	var eID [32]byte

	copy(eID[:], ehrUUID[:])

	packed, err = i.pack("addEhrDoc",
		eID,
		(ehrIndexer.EhrIndexerDocumentMeta)(*docMeta),
	)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.AddEhrDoc error: %w", err)
	}

	return
}

func (i *Index) GetDocLastByType(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType) (*model.DocumentMeta, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		eID      [32]byte
	)

	copy(eID[:], ehrUUID[:])

	docMeta, err := i.ehrIndex.GetLastEhrDocByType(callOpts, eID, uint8(docType))
	if err != nil {
		if err.Error() == ExecutionRevertedNFD {
			return nil, fmt.Errorf("ehrIndex.GetLastEhrDocByType error: %w", errors.ErrNotFound)
		}
		return nil, fmt.Errorf("ehrIndex.GetLastEhrDocByType error: %w ehrUUID %s docType %s", err, ehrUUID.String(), docType.String())
	}

	return (*model.DocumentMeta)(&docMeta), nil
}

func (i *Index) GetDocLastByBaseID(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType, docBaseUIDHash *[32]byte) (*model.DocumentMeta, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		eID      [32]byte
	)

	copy(eID[:], ehrUUID[:])

	docMeta, err := i.ehrIndex.GetDocLastByBaseID(callOpts, eID, uint8(docType), *docBaseUIDHash)
	if err != nil {
		if err.Error() == ExecutionRevertedNFD {
			return nil, fmt.Errorf("ehrIndex.GetDocLastByBaseID error: %w", errors.ErrNotFound)
		}
		return nil, fmt.Errorf("ehrIndex.GetDocLastByBaseID error: %w ehrUUID %s docType %s docBaseUIDHash %x", err, ehrUUID.String(), docType.String(), docBaseUIDHash)
	}

	return (*model.DocumentMeta)(&docMeta), nil
}

func (i *Index) GetDocByTime(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType, timestamp uint32) (*model.DocumentMeta, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		eID      [32]byte
	)

	copy(eID[:], ehrUUID[:])

	docMeta, err := i.ehrIndex.GetDocByTime(callOpts, eID, uint8(docType), timestamp)
	if err != nil {
		if err.Error() == ExecutionRevertedNFD {
			return nil, fmt.Errorf("ehrIndex.GetDocByTime error: %w", errors.ErrNotFound)
		}
		return nil, fmt.Errorf("ehrIndex.GetDocByTime error: %w ehrUUID %s docType %s timestamp %d", err, ehrUUID.String(), docType.String(), timestamp)
	}

	return (*model.DocumentMeta)(&docMeta), nil
}

func (i *Index) GetDocByVersion(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType, docBaseUIDHash *[32]byte, version *[32]byte) (*model.DocumentMeta, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		eID      [32]byte
	)

	copy(eID[:], ehrUUID[:])

	docMeta, err := i.ehrIndex.GetDocByVersion(callOpts, eID, uint8(docType), *docBaseUIDHash, *version)
	if err != nil {
		if err.Error() == ExecutionRevertedNFD {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("ehrIndex.GetDocByVersion error: %w ehrUUID %s docType %s docBaseUIDHash %x version %s", err, ehrUUID.String(), docType.String(), docBaseUIDHash, version)
	}

	return (*model.DocumentMeta)(&docMeta), nil
}

func (i *Index) SetDocKeyEncrypted(key *[32]byte, value []byte) (packed []byte, err error) {
	packed, err = i.pack("setDocAccess", *key, value)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.SetDocAccess error: %w", err)
	}

	return
}

func (i *Index) GetDocKeyEncrypted(ctx context.Context, userID string, CID *cid.Cid) ([]byte, error) {
	docAccessIndexKey := sha3.Sum256(append(CID.Bytes()[:], []byte(userID)...))

	callOpts := &bind.CallOpts{
		Context: ctx,
	}

	docAccessValue, err := i.ehrIndex.DocAccess(callOpts, docAccessIndexKey)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.DocAccess error: %w", err)
	}

	return docAccessValue, nil
}

func (i *Index) SetGroupAccess(ctx context.Context, key *[32]byte, value []byte) (string, error) {
	i.Lock()
	defer i.Unlock()

	tx, err := i.ehrIndex.SetGroupAccess(i.transactOpts, *key, value)
	if err != nil {
		return "", fmt.Errorf("ehrIndex.SetGroupAccess error: %w", err)
	}

	return tx.Hash().Hex(), nil
}

func (i *Index) GetGroupAccess(ctx context.Context, userID string, groupUUID *uuid.UUID) ([]byte, error) {
	groupAccessIndexKey := sha3.Sum256(append([]byte(userID), groupUUID[:]...))

	callOpts := &bind.CallOpts{
		Context: ctx,
	}

	groupAccessValue, err := i.ehrIndex.GroupAccess(callOpts, groupAccessIndexKey)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.GroupAccess error: %w", err)
	}

	if len(groupAccessValue) == 0 {
		return nil, errors.ErrIsNotExist
	}

	return groupAccessValue, nil
}

func (i *Index) UserAdd(requestID string, userAddr common.Address, userID string, role uint8, pwdHash []byte) (string, error) {
	var uID [32]byte

	copy(uID[:], userID[:])

	i.Lock()
	defer i.Unlock()

	tx, err := i.ehrIndex.UserAdd(i.transactOpts, userAddr, uID, role, pwdHash)
	if err != nil {
		if err.Error() == ExecutionRevertedNFD {
			return "", errors.ErrNotFound
		} else if err.Error() == "execution reverted: ADL" {
			return "", errors.ErrAlreadyDeleted
		}
		return "", fmt.Errorf("ehrIndex.UserAdd error: %w", err)
	}

	log.Printf("%s UserAdd tx %s nonce %d", requestID, tx.Hash().Hex(), tx.Nonce())

	return tx.Hash().Hex(), nil
}

func (i *Index) GetUserPasswordHash(ctx context.Context, userAddr common.Address) ([]byte, error) {
	callOpts := &bind.CallOpts{
		Context: ctx,
	}

	userPasswordHash, err := i.ehrIndex.GetUserPasswordHash(callOpts, userAddr)
	if err != nil {
		if err.Error() == ExecutionRevertedNFD {
			return nil, fmt.Errorf("ehrIndex.GetUserPasswordHash error: %w", errors.ErrNotFound)
		}
		return nil, fmt.Errorf("ehrIndex.GetUserPasswordHash error: %w userAddr %s", err, userAddr.String())
	}

	return userPasswordHash, nil
}

func (i *Index) SetSubject(ehrUUID *uuid.UUID, subjectID, subjectNamespace string) (packed []byte, err error) {
	var eID [32]byte

	copy(eID[:], ehrUUID[:])

	subjectKey := sha3.Sum256([]byte(subjectID + subjectNamespace))

	packed, err = i.pack("setEhrSubject", subjectKey, eID)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.SetSubject error: %w", err)
	}

	return
}

func (i *Index) GetEhrUUIDBySubject(ctx context.Context, subjectID, subjectNamespace string) (*uuid.UUID, error) {
	subjectKey := sha3.Sum256([]byte(subjectID + subjectNamespace))

	callOpts := &bind.CallOpts{
		Context: ctx,
	}

	ehrUUIDRaw, err := i.ehrIndex.EhrSubject(callOpts, subjectKey)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.EhrSubjec error: %w", err)
	}

	ehrUUID, err := uuid.FromBytes(ehrUUIDRaw[:16])
	if err != nil {
		return nil, fmt.Errorf("ehrUUID FromBytes error: %w ehrUUIDRaw %x", err, ehrUUIDRaw)
	}

	return &ehrUUID, nil
}

func (i *Index) DeleteDoc(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType, docBaseUIDHash *[32]byte, version *[32]byte) (string, error) {
	var eID [32]byte

	copy(eID[:], ehrUUID[:])

	i.Lock()
	defer i.Unlock()

	tx, err := i.ehrIndex.DeleteDoc(i.transactOpts, eID, uint8(docType), *docBaseUIDHash, *version)
	if err != nil {
		if err.Error() == ExecutionRevertedNFD {
			return "", errors.ErrNotFound
		} else if err.Error() == "execution reverted: ADL" {
			return "", errors.ErrAlreadyDeleted
		}
		return "", fmt.Errorf("ehrIndex.DeleteDoc error: %w ehrUUID %s docType %s", err, ehrUUID.String(), docType.String())
	}

	log.Printf("%s DeleteDoc tx %s nonce %d", ctx.(*gin.Context).GetString("reqId"), tx.Hash().Hex(), tx.Nonce())

	return tx.Hash().Hex(), nil
}

func (i *Index) SetAllowed(ctx context.Context, address string) (string, error) {
	i.Lock()
	defer i.Unlock()

	tx, err := i.ehrIndex.SetAllowed(i.transactOpts, common.HexToAddress(address), true)
	if err != nil {
		return "", fmt.Errorf("ehrIndex.SetAllowed error: %w", err)
	}

	log.Printf("%s SetAllowed tx %s nonce %d", ctx.(*gin.Context).GetString("reqId"), tx.Hash().Hex(), tx.Nonce())

	return tx.Hash().Hex(), nil
}

func (i *Index) TxWait(ctx context.Context, hash string) (uint64, error) {
	h := common.HexToHash(hash)

	ticker := time.NewTicker(5 * time.Second)

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
