package indexer

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
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
	"github.com/google/uuid"
	"github.com/vmihailenco/msgpack/v5"
	"golang.org/x/crypto/sha3"

	"hms/gateway/pkg/access"
	"hms/gateway/pkg/docs/model"
	"hms/gateway/pkg/docs/types"
	"hms/gateway/pkg/errors"
	"hms/gateway/pkg/indexer/ehrIndexer"
	log "hms/gateway/pkg/log"
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
	abi          *abi.ABI
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

func New(contractAddr, keyPath string, client *ethclient.Client, gasTipCap int64) *Index {
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

	if gasTipCap > 0 {
		transactOpts.GasTipCap = big.NewInt(gasTipCap)
	}

	return &Index{
		client:       client,
		ehrIndex:     ehrIndex,
		transactOpts: transactOpts,
		abi:          bcAbi,
	}
}

func (i *Index) pack(name string, args ...interface{}) ([]byte, error) {
	result, err := i.abi.Pack(name, args...)
	if err != nil {
		return nil, fmt.Errorf("abi.Pack error: %w", err)
	}
	return result, nil
}

func (i *Index) SetEhrUser(ctx context.Context, userID string, ehrUUID *uuid.UUID, privKey *[32]byte, nonce *big.Int) ([]byte, error) {
	var uID, eID [32]byte

	copy(uID[:], userID)
	copy(eID[:], ehrUUID[:])

	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	if nonce == nil {
		nonce, err = i.userNonce(ctx, &userAddress)
		if err != nil {
			return nil, fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	sig, err := makeSignature(userKey, nonce, "setEhrUser", uID, eID)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err := i.pack("setEhrUser", uID, eID, userAddress, sig)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.SetEhrUser error: %w", err)
	}

	return data, err
}

func (i *Index) GetEhrUUIDByUserID(ctx context.Context, userID string) (*uuid.UUID, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		uID      [32]byte
	)

	copy(uID[:], userID)

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

func (i *Index) AddEhrDoc(ctx context.Context, ehrUUID *uuid.UUID, docMeta *model.DocumentMeta, keyEncrypted, CIDEncr []byte, privKey *[32]byte, nonce *big.Int) ([]byte, error) {
	var eID [32]byte

	copy(eID[:], ehrUUID[:])

	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	if nonce == nil {
		nonce, err = i.userNonce(ctx, &userAddress)
		if err != nil {
			return nil, fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	sig, err := makeSignature(userKey, nonce, "addEhrDoc", eID, *docMeta, keyEncrypted, CIDEncr)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	params := ehrIndexer.DocsAddEhrDocParams{
		EhrId:     eID,
		DocMeta:   (ehrIndexer.DocsDocumentMeta)(*docMeta),
		KeyEncr:   keyEncrypted,
		CIDEncr:   CIDEncr,
		Signer:    userAddress,
		Signature: sig,
	}

	data, err := i.pack("addEhrDoc", params)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.AddEhrDoc error: %w", err)
	}

	/*
		log.Printf("signature: %x", sig)
		log.Printf("userAddress: %x", userAddress.Bytes())
		log.Printf("data: %x", data)
	*/

	return data, nil
}

func (i *Index) GetDocLastByType(ctx context.Context, ehrUUID *uuid.UUID, docType types.DocumentType) (*model.DocumentMeta, error) {
	var (
		callOpts = &bind.CallOpts{Context: ctx}
		eID      [32]byte
	)

	copy(eID[:], ehrUUID[:])

	docMeta, err := i.ehrIndex.GetLastEhrDocByType(callOpts, eID, uint8(docType))
	if err != nil {
		if strings.Contains(err.Error(), "NFD") {
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
		if strings.Contains(err.Error(), "NFD") {
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
		if strings.Contains(err.Error(), "NFD") {
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
		if strings.Contains(err.Error(), "NFD") {
			return nil, errors.ErrNotFound
		}
		return nil, fmt.Errorf("ehrIndex.GetDocByVersion error: %w ehrUUID %s docType %s docBaseUIDHash %x version %s", err, ehrUUID.String(), docType.String(), docBaseUIDHash, version)
	}

	return (*model.DocumentMeta)(&docMeta), nil
}

func (i *Index) GetDocKeyEncrypted(ctx context.Context, userID string, CID []byte) ([]byte, error) {
	var uID [32]byte

	copy(uID[:], userID)

	data, err := abi.Arguments{{Type: Bytes32}, {Type: Uint8}}.Pack(uID, access.Doc)
	if err != nil {
		return nil, fmt.Errorf("args.Pack error: %w", err)
	}

	accessID := crypto.Keccak256Hash(data)

	data, err = abi.Arguments{{Type: Bytes}}.Pack(CID)
	if err != nil {
		return nil, fmt.Errorf("args.Pack error: %w", err)
	}

	CIDHash := crypto.Keccak256Hash(data)

	callOpts := &bind.CallOpts{
		Context: ctx,
	}

	accessObj, err := i.ehrIndex.GetAccessByIdHash(callOpts, accessID, CIDHash)
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

func makeSignature(pk *ecdsa.PrivateKey, nonce *big.Int, values ...interface{}) ([]byte, error) {
	var args abi.Arguments

	for i, a := range values {
		switch a.(type) {
		case string:
			args = append(args, abi.Argument{Type: String})
		case [32]byte:
			args = append(args, abi.Argument{Type: Bytes32})
		case []byte:
			args = append(args, abi.Argument{Type: Bytes})
		case uint8:
			args = append(args, abi.Argument{Type: Uint8})
		case *big.Int:
			args = append(args, abi.Argument{Type: Uint256})
		case common.Address:
			args = append(args, abi.Argument{Type: Address})
		case ehrIndexer.AccessObject:
			args = append(args, abi.Argument{Type: Access})
		case model.DocumentMeta, ehrIndexer.DocsDocumentMeta:
			args = append(args, abi.Argument{Type: DocMeta})
		default:
			return nil, fmt.Errorf("%w: makeSignature unknown %d argument type: %v", errors.ErrIncorrectFormat, i, a)
		}
	}

	data, err := args.Pack(values...)
	if err != nil {
		return nil, fmt.Errorf("args.Pack error: %w args: %v values: %v", err, args, values)
	}

	hash := crypto.Keccak256Hash(data)

	nonceBytes, _ := abi.Arguments{{Type: Uint256}}.Pack(nonce)

	prefixedHash := crypto.Keccak256Hash(
		[]byte("\x19Ethereum Signed Message:\n32"),
		hash.Bytes(),
		nonceBytes,
	)

	sig, err := crypto.Sign(prefixedHash.Bytes(), pk)
	if err != nil {
		return nil, fmt.Errorf("crypto.Sign error: %w", err)
	}

	sig[64] += 27

	return sig, nil
}

func (i *Index) SetEhrSubject(ctx context.Context, ehrUUID *uuid.UUID, subjectID, subjectNamespace string, privKey *[32]byte, nonce *big.Int) ([]byte, error) {
	var eID [32]byte

	copy(eID[:], ehrUUID[:])

	subjectKey := sha3.Sum256([]byte(subjectID + subjectNamespace))

	userKey, err := crypto.ToECDSA(privKey[:])
	if err != nil {
		return nil, fmt.Errorf("crypto.ToECDSA error: %w", err)
	}

	userAddress := crypto.PubkeyToAddress(userKey.PublicKey)

	if nonce == nil {
		nonce, err = i.userNonce(ctx, &userAddress)
		if err != nil {
			return nil, fmt.Errorf("userNonce error: %w address: %s", err, userAddress.String())
		}
	}

	sig, err := makeSignature(userKey, nonce, "setEhrSubject", subjectKey, eID)
	if err != nil {
		return nil, fmt.Errorf("makeSignature error: %w", err)
	}

	data, err := i.pack("setEhrSubject", subjectKey, eID, userAddress, sig)
	if err != nil {
		return nil, fmt.Errorf("ehrIndex.SetSubject error: %w", err)
	}

	return data, nil
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
		if strings.Contains(err.Error(), "NFD") {
			return "", errors.ErrNotFound
		} else if strings.Contains(err.Error(), "ADL") {
			return "", errors.ErrAlreadyDeleted
		}
		return "", fmt.Errorf("ehrIndex.DeleteDoc error: %w ehrUUID %s docType %s", err, ehrUUID.String(), docType.String())
	}

	return tx.Hash().Hex(), nil
}

func (i *Index) SetAllowed(ctx context.Context, address string) (string, error) {
	i.Lock()
	defer i.Unlock()

	tx, err := i.ehrIndex.SetAllowed(i.transactOpts, common.HexToAddress(address), true)
	if err != nil {
		return "", fmt.Errorf("ehrIndex.SetAllowed error: %w", err)
	}

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

func Keccak256(data []byte) []byte {
	data, _ = abi.Arguments{{Type: Bytes}}.Pack(data)

	return crypto.Keccak256Hash(data).Bytes()
}
