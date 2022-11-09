// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ehrIndexer

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// AccessObject is an auto generated low-level Go binding around an user-defined struct.
type AccessObject struct {
	IdHash  [32]byte
	IdEncr  []byte
	KeyEncr []byte
	Level   uint8
}

// EhrDocsAddEhrDocParams is an auto generated low-level Go binding around an user-defined struct.
type EhrDocsAddEhrDocParams struct {
	EhrId     [32]byte
	DocMeta   EhrDocsDocumentMeta
	KeyEncr   []byte
	CIDEncr   []byte
	Nonce     *big.Int
	Signer    common.Address
	Signature []byte
}

// EhrDocsDocumentMeta is an auto generated low-level Go binding around an user-defined struct.
type EhrDocsDocumentMeta struct {
	DocType         uint8
	Status          uint8
	CID             []byte
	DealCID         []byte
	MinerAddress    []byte
	DocUIDEncrypted []byte
	DocBaseUIDHash  [32]byte
	Version         [32]byte
	IsLast          bool
	Timestamp       uint32
}

// UsersGroupAddUserParams is an auto generated low-level Go binding around an user-defined struct.
type UsersGroupAddUserParams struct {
	GroupIdHash    [32]byte
	AddingUserAddr common.Address
	Level          uint8
	IdEncr         []byte
	KeyEncr        []byte
	Nonce          *big.Int
	Signer         common.Address
	Signature      []byte
}

// UsersKeyValue is an auto generated low-level Go binding around an user-defined struct.
type UsersKeyValue struct {
	Key   [32]byte
	Value []byte
}

// UsersUserGroupCreateParams is an auto generated low-level Go binding around an user-defined struct.
type UsersUserGroupCreateParams struct {
	GroupIdHash  [32]byte
	GroupIdEncr  []byte
	GroupKeyEncr []byte
	Params       []UsersKeyValue
	Nonce        *big.Int
	Signer       common.Address
	Signature    []byte
}

// EhrIndexerMetaData contains all meta data concerning the EhrIndexer contract.
var EhrIndexerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"enumEhrDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"enumEhrDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"CID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"dealCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"minerAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"docUIDEncrypted\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structEhrDocs.DocumentMeta\",\"name\":\"docMeta\",\"type\":\"tuple\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"CIDEncr\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structEhrDocs.AddEhrDocParams\",\"name\":\"p\",\"type\":\"tuple\"}],\"name\":\"addEhrDoc\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedChange\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dataSearch\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeType\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"nodeID\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumEhrDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"}],\"name\":\"deleteDoc\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"ehrSubject\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"ehrUsers\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"objectIdHash\",\"type\":\"bytes32\"}],\"name\":\"getAccessByIdHash\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumAccess.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structAccess.Object\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrID\",\"type\":\"bytes32\"},{\"internalType\":\"enumEhrDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"name\":\"getDocByTime\",\"outputs\":[{\"components\":[{\"internalType\":\"enumEhrDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"enumEhrDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"CID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"dealCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"minerAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"docUIDEncrypted\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structEhrDocs.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumEhrDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"}],\"name\":\"getDocByVersion\",\"outputs\":[{\"components\":[{\"internalType\":\"enumEhrDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"enumEhrDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"CID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"dealCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"minerAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"docUIDEncrypted\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structEhrDocs.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumEhrDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"}],\"name\":\"getDocLastByBaseID\",\"outputs\":[{\"components\":[{\"internalType\":\"enumEhrDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"enumEhrDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"CID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"dealCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"minerAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"docUIDEncrypted\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structEhrDocs.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumEhrDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"}],\"name\":\"getEhrDocs\",\"outputs\":[{\"components\":[{\"internalType\":\"enumEhrDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"enumEhrDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"CID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"dealCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"minerAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"docUIDEncrypted\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structEhrDocs.DocumentMeta[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumEhrDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"}],\"name\":\"getLastEhrDocByType\",\"outputs\":[{\"components\":[{\"internalType\":\"enumEhrDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"enumEhrDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"CID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"dealCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"minerAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"docUIDEncrypted\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structEhrDocs.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userIdHash\",\"type\":\"bytes32\"}],\"name\":\"getUserAccessList\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumAccess.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structAccess.Object[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAddr\",\"type\":\"address\"}],\"name\":\"getUserPasswordHash\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"groupIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"addingUserAddr\",\"type\":\"address\"},{\"internalType\":\"enumAccess.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUsers.GroupAddUserParams\",\"name\":\"p\",\"type\":\"tuple\"}],\"name\":\"groupAddUser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"groupIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"removingUserAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"groupRemoveUser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAllowed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes\",\"name\":\"CID\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumAccess.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structAccess.Object\",\"name\":\"accessObj\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"userAddr\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"setDocAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"subjectKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"setEhrSubject\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"setEhrUser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"groupIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"groupIdEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"groupKeyEncr\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structUsers.KeyValue[]\",\"name\":\"params\",\"type\":\"tuple[]\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUsers.UserGroupCreateParams\",\"name\":\"p\",\"type\":\"tuple\"}],\"name\":\"userGroupCreate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAddr\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"systemID\",\"type\":\"bytes32\"},{\"internalType\":\"enumUsers.Role\",\"name\":\"role\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"pwdHash\",\"type\":\"bytes\"},{\"internalType\":\"uint256\",\"name\":\"nonce\",\"type\":\"uint256\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"userNew\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"users\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"systemID\",\"type\":\"bytes32\"},{\"internalType\":\"enumUsers.Role\",\"name\":\"role\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"pwdHash\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// EhrIndexerABI is the input ABI used to generate the binding from.
// Deprecated: Use EhrIndexerMetaData.ABI instead.
var EhrIndexerABI = EhrIndexerMetaData.ABI

// EhrIndexer is an auto generated Go binding around an Ethereum contract.
type EhrIndexer struct {
	EhrIndexerCaller     // Read-only binding to the contract
	EhrIndexerTransactor // Write-only binding to the contract
	EhrIndexerFilterer   // Log filterer for contract events
}

// EhrIndexerCaller is an auto generated read-only Go binding around an Ethereum contract.
type EhrIndexerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EhrIndexerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EhrIndexerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EhrIndexerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EhrIndexerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EhrIndexerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EhrIndexerSession struct {
	Contract     *EhrIndexer       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EhrIndexerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EhrIndexerCallerSession struct {
	Contract *EhrIndexerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// EhrIndexerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EhrIndexerTransactorSession struct {
	Contract     *EhrIndexerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// EhrIndexerRaw is an auto generated low-level Go binding around an Ethereum contract.
type EhrIndexerRaw struct {
	Contract *EhrIndexer // Generic contract binding to access the raw methods on
}

// EhrIndexerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EhrIndexerCallerRaw struct {
	Contract *EhrIndexerCaller // Generic read-only contract binding to access the raw methods on
}

// EhrIndexerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EhrIndexerTransactorRaw struct {
	Contract *EhrIndexerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEhrIndexer creates a new instance of EhrIndexer, bound to a specific deployed contract.
func NewEhrIndexer(address common.Address, backend bind.ContractBackend) (*EhrIndexer, error) {
	contract, err := bindEhrIndexer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &EhrIndexer{EhrIndexerCaller: EhrIndexerCaller{contract: contract}, EhrIndexerTransactor: EhrIndexerTransactor{contract: contract}, EhrIndexerFilterer: EhrIndexerFilterer{contract: contract}}, nil
}

// NewEhrIndexerCaller creates a new read-only instance of EhrIndexer, bound to a specific deployed contract.
func NewEhrIndexerCaller(address common.Address, caller bind.ContractCaller) (*EhrIndexerCaller, error) {
	contract, err := bindEhrIndexer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EhrIndexerCaller{contract: contract}, nil
}

// NewEhrIndexerTransactor creates a new write-only instance of EhrIndexer, bound to a specific deployed contract.
func NewEhrIndexerTransactor(address common.Address, transactor bind.ContractTransactor) (*EhrIndexerTransactor, error) {
	contract, err := bindEhrIndexer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EhrIndexerTransactor{contract: contract}, nil
}

// NewEhrIndexerFilterer creates a new log filterer instance of EhrIndexer, bound to a specific deployed contract.
func NewEhrIndexerFilterer(address common.Address, filterer bind.ContractFilterer) (*EhrIndexerFilterer, error) {
	contract, err := bindEhrIndexer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EhrIndexerFilterer{contract: contract}, nil
}

// bindEhrIndexer binds a generic wrapper to an already deployed contract.
func bindEhrIndexer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EhrIndexerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EhrIndexer *EhrIndexerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EhrIndexer.Contract.EhrIndexerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EhrIndexer *EhrIndexerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EhrIndexer.Contract.EhrIndexerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EhrIndexer *EhrIndexerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EhrIndexer.Contract.EhrIndexerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_EhrIndexer *EhrIndexerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _EhrIndexer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_EhrIndexer *EhrIndexerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EhrIndexer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_EhrIndexer *EhrIndexerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _EhrIndexer.Contract.contract.Transact(opts, method, params...)
}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_EhrIndexer *EhrIndexerCaller) AllowedChange(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "allowedChange", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_EhrIndexer *EhrIndexerSession) AllowedChange(arg0 common.Address) (bool, error) {
	return _EhrIndexer.Contract.AllowedChange(&_EhrIndexer.CallOpts, arg0)
}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_EhrIndexer *EhrIndexerCallerSession) AllowedChange(arg0 common.Address) (bool, error) {
	return _EhrIndexer.Contract.AllowedChange(&_EhrIndexer.CallOpts, arg0)
}

// DataSearch is a free data retrieval call binding the contract method 0x4be82179.
//
// Solidity: function dataSearch() view returns(bytes32 nodeType, bytes32 nodeID)
func (_EhrIndexer *EhrIndexerCaller) DataSearch(opts *bind.CallOpts) (struct {
	NodeType [32]byte
	NodeID   [32]byte
}, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "dataSearch")

	outstruct := new(struct {
		NodeType [32]byte
		NodeID   [32]byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.NodeType = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.NodeID = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)

	return *outstruct, err

}

// DataSearch is a free data retrieval call binding the contract method 0x4be82179.
//
// Solidity: function dataSearch() view returns(bytes32 nodeType, bytes32 nodeID)
func (_EhrIndexer *EhrIndexerSession) DataSearch() (struct {
	NodeType [32]byte
	NodeID   [32]byte
}, error) {
	return _EhrIndexer.Contract.DataSearch(&_EhrIndexer.CallOpts)
}

// DataSearch is a free data retrieval call binding the contract method 0x4be82179.
//
// Solidity: function dataSearch() view returns(bytes32 nodeType, bytes32 nodeID)
func (_EhrIndexer *EhrIndexerCallerSession) DataSearch() (struct {
	NodeType [32]byte
	NodeID   [32]byte
}, error) {
	return _EhrIndexer.Contract.DataSearch(&_EhrIndexer.CallOpts)
}

// EhrSubject is a free data retrieval call binding the contract method 0xfe1b5580.
//
// Solidity: function ehrSubject(bytes32 ) view returns(bytes32)
func (_EhrIndexer *EhrIndexerCaller) EhrSubject(opts *bind.CallOpts, arg0 [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "ehrSubject", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EhrSubject is a free data retrieval call binding the contract method 0xfe1b5580.
//
// Solidity: function ehrSubject(bytes32 ) view returns(bytes32)
func (_EhrIndexer *EhrIndexerSession) EhrSubject(arg0 [32]byte) ([32]byte, error) {
	return _EhrIndexer.Contract.EhrSubject(&_EhrIndexer.CallOpts, arg0)
}

// EhrSubject is a free data retrieval call binding the contract method 0xfe1b5580.
//
// Solidity: function ehrSubject(bytes32 ) view returns(bytes32)
func (_EhrIndexer *EhrIndexerCallerSession) EhrSubject(arg0 [32]byte) ([32]byte, error) {
	return _EhrIndexer.Contract.EhrSubject(&_EhrIndexer.CallOpts, arg0)
}

// EhrUsers is a free data retrieval call binding the contract method 0x5113c226.
//
// Solidity: function ehrUsers(bytes32 ) view returns(bytes32)
func (_EhrIndexer *EhrIndexerCaller) EhrUsers(opts *bind.CallOpts, arg0 [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "ehrUsers", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EhrUsers is a free data retrieval call binding the contract method 0x5113c226.
//
// Solidity: function ehrUsers(bytes32 ) view returns(bytes32)
func (_EhrIndexer *EhrIndexerSession) EhrUsers(arg0 [32]byte) ([32]byte, error) {
	return _EhrIndexer.Contract.EhrUsers(&_EhrIndexer.CallOpts, arg0)
}

// EhrUsers is a free data retrieval call binding the contract method 0x5113c226.
//
// Solidity: function ehrUsers(bytes32 ) view returns(bytes32)
func (_EhrIndexer *EhrIndexerCallerSession) EhrUsers(arg0 [32]byte) ([32]byte, error) {
	return _EhrIndexer.Contract.EhrUsers(&_EhrIndexer.CallOpts, arg0)
}

// GetAccessByIdHash is a free data retrieval call binding the contract method 0x9ae2da76.
//
// Solidity: function getAccessByIdHash(bytes32 userIdHash, bytes32 objectIdHash) view returns((bytes32,bytes,bytes,uint8))
func (_EhrIndexer *EhrIndexerCaller) GetAccessByIdHash(opts *bind.CallOpts, userIdHash [32]byte, objectIdHash [32]byte) (AccessObject, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getAccessByIdHash", userIdHash, objectIdHash)

	if err != nil {
		return *new(AccessObject), err
	}

	out0 := *abi.ConvertType(out[0], new(AccessObject)).(*AccessObject)

	return out0, err

}

// GetAccessByIdHash is a free data retrieval call binding the contract method 0x9ae2da76.
//
// Solidity: function getAccessByIdHash(bytes32 userIdHash, bytes32 objectIdHash) view returns((bytes32,bytes,bytes,uint8))
func (_EhrIndexer *EhrIndexerSession) GetAccessByIdHash(userIdHash [32]byte, objectIdHash [32]byte) (AccessObject, error) {
	return _EhrIndexer.Contract.GetAccessByIdHash(&_EhrIndexer.CallOpts, userIdHash, objectIdHash)
}

// GetAccessByIdHash is a free data retrieval call binding the contract method 0x9ae2da76.
//
// Solidity: function getAccessByIdHash(bytes32 userIdHash, bytes32 objectIdHash) view returns((bytes32,bytes,bytes,uint8))
func (_EhrIndexer *EhrIndexerCallerSession) GetAccessByIdHash(userIdHash [32]byte, objectIdHash [32]byte) (AccessObject, error) {
	return _EhrIndexer.Contract.GetAccessByIdHash(&_EhrIndexer.CallOpts, userIdHash, objectIdHash)
}

// GetDocByTime is a free data retrieval call binding the contract method 0x4f722b16.
//
// Solidity: function getDocByTime(bytes32 ehrID, uint8 docType, uint32 timestamp) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerCaller) GetDocByTime(opts *bind.CallOpts, ehrID [32]byte, docType uint8, timestamp uint32) (EhrDocsDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getDocByTime", ehrID, docType, timestamp)

	if err != nil {
		return *new(EhrDocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(EhrDocsDocumentMeta)).(*EhrDocsDocumentMeta)

	return out0, err

}

// GetDocByTime is a free data retrieval call binding the contract method 0x4f722b16.
//
// Solidity: function getDocByTime(bytes32 ehrID, uint8 docType, uint32 timestamp) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerSession) GetDocByTime(ehrID [32]byte, docType uint8, timestamp uint32) (EhrDocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocByTime(&_EhrIndexer.CallOpts, ehrID, docType, timestamp)
}

// GetDocByTime is a free data retrieval call binding the contract method 0x4f722b16.
//
// Solidity: function getDocByTime(bytes32 ehrID, uint8 docType, uint32 timestamp) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerCallerSession) GetDocByTime(ehrID [32]byte, docType uint8, timestamp uint32) (EhrDocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocByTime(&_EhrIndexer.CallOpts, ehrID, docType, timestamp)
}

// GetDocByVersion is a free data retrieval call binding the contract method 0x179fcacf.
//
// Solidity: function getDocByVersion(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerCaller) GetDocByVersion(opts *bind.CallOpts, ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte) (EhrDocsDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getDocByVersion", ehrId, docType, docBaseUIDHash, version)

	if err != nil {
		return *new(EhrDocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(EhrDocsDocumentMeta)).(*EhrDocsDocumentMeta)

	return out0, err

}

// GetDocByVersion is a free data retrieval call binding the contract method 0x179fcacf.
//
// Solidity: function getDocByVersion(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerSession) GetDocByVersion(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte) (EhrDocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocByVersion(&_EhrIndexer.CallOpts, ehrId, docType, docBaseUIDHash, version)
}

// GetDocByVersion is a free data retrieval call binding the contract method 0x179fcacf.
//
// Solidity: function getDocByVersion(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerCallerSession) GetDocByVersion(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte) (EhrDocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocByVersion(&_EhrIndexer.CallOpts, ehrId, docType, docBaseUIDHash, version)
}

// GetDocLastByBaseID is a free data retrieval call binding the contract method 0x949c8e77.
//
// Solidity: function getDocLastByBaseID(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerCaller) GetDocLastByBaseID(opts *bind.CallOpts, ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte) (EhrDocsDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getDocLastByBaseID", ehrId, docType, docBaseUIDHash)

	if err != nil {
		return *new(EhrDocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(EhrDocsDocumentMeta)).(*EhrDocsDocumentMeta)

	return out0, err

}

// GetDocLastByBaseID is a free data retrieval call binding the contract method 0x949c8e77.
//
// Solidity: function getDocLastByBaseID(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerSession) GetDocLastByBaseID(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte) (EhrDocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocLastByBaseID(&_EhrIndexer.CallOpts, ehrId, docType, docBaseUIDHash)
}

// GetDocLastByBaseID is a free data retrieval call binding the contract method 0x949c8e77.
//
// Solidity: function getDocLastByBaseID(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerCallerSession) GetDocLastByBaseID(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte) (EhrDocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocLastByBaseID(&_EhrIndexer.CallOpts, ehrId, docType, docBaseUIDHash)
}

// GetEhrDocs is a free data retrieval call binding the contract method 0xeec6b331.
//
// Solidity: function getEhrDocs(bytes32 ehrId, uint8 docType) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32)[])
func (_EhrIndexer *EhrIndexerCaller) GetEhrDocs(opts *bind.CallOpts, ehrId [32]byte, docType uint8) ([]EhrDocsDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getEhrDocs", ehrId, docType)

	if err != nil {
		return *new([]EhrDocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new([]EhrDocsDocumentMeta)).(*[]EhrDocsDocumentMeta)

	return out0, err

}

// GetEhrDocs is a free data retrieval call binding the contract method 0xeec6b331.
//
// Solidity: function getEhrDocs(bytes32 ehrId, uint8 docType) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32)[])
func (_EhrIndexer *EhrIndexerSession) GetEhrDocs(ehrId [32]byte, docType uint8) ([]EhrDocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetEhrDocs(&_EhrIndexer.CallOpts, ehrId, docType)
}

// GetEhrDocs is a free data retrieval call binding the contract method 0xeec6b331.
//
// Solidity: function getEhrDocs(bytes32 ehrId, uint8 docType) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32)[])
func (_EhrIndexer *EhrIndexerCallerSession) GetEhrDocs(ehrId [32]byte, docType uint8) ([]EhrDocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetEhrDocs(&_EhrIndexer.CallOpts, ehrId, docType)
}

// GetLastEhrDocByType is a free data retrieval call binding the contract method 0x15dbcf1a.
//
// Solidity: function getLastEhrDocByType(bytes32 ehrId, uint8 docType) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerCaller) GetLastEhrDocByType(opts *bind.CallOpts, ehrId [32]byte, docType uint8) (EhrDocsDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getLastEhrDocByType", ehrId, docType)

	if err != nil {
		return *new(EhrDocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(EhrDocsDocumentMeta)).(*EhrDocsDocumentMeta)

	return out0, err

}

// GetLastEhrDocByType is a free data retrieval call binding the contract method 0x15dbcf1a.
//
// Solidity: function getLastEhrDocByType(bytes32 ehrId, uint8 docType) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerSession) GetLastEhrDocByType(ehrId [32]byte, docType uint8) (EhrDocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetLastEhrDocByType(&_EhrIndexer.CallOpts, ehrId, docType)
}

// GetLastEhrDocByType is a free data retrieval call binding the contract method 0x15dbcf1a.
//
// Solidity: function getLastEhrDocByType(bytes32 ehrId, uint8 docType) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerCallerSession) GetLastEhrDocByType(ehrId [32]byte, docType uint8) (EhrDocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetLastEhrDocByType(&_EhrIndexer.CallOpts, ehrId, docType)
}

// GetUserAccessList is a free data retrieval call binding the contract method 0xbb059b5c.
//
// Solidity: function getUserAccessList(bytes32 userIdHash) view returns((bytes32,bytes,bytes,uint8)[])
func (_EhrIndexer *EhrIndexerCaller) GetUserAccessList(opts *bind.CallOpts, userIdHash [32]byte) ([]AccessObject, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getUserAccessList", userIdHash)

	if err != nil {
		return *new([]AccessObject), err
	}

	out0 := *abi.ConvertType(out[0], new([]AccessObject)).(*[]AccessObject)

	return out0, err

}

// GetUserAccessList is a free data retrieval call binding the contract method 0xbb059b5c.
//
// Solidity: function getUserAccessList(bytes32 userIdHash) view returns((bytes32,bytes,bytes,uint8)[])
func (_EhrIndexer *EhrIndexerSession) GetUserAccessList(userIdHash [32]byte) ([]AccessObject, error) {
	return _EhrIndexer.Contract.GetUserAccessList(&_EhrIndexer.CallOpts, userIdHash)
}

// GetUserAccessList is a free data retrieval call binding the contract method 0xbb059b5c.
//
// Solidity: function getUserAccessList(bytes32 userIdHash) view returns((bytes32,bytes,bytes,uint8)[])
func (_EhrIndexer *EhrIndexerCallerSession) GetUserAccessList(userIdHash [32]byte) ([]AccessObject, error) {
	return _EhrIndexer.Contract.GetUserAccessList(&_EhrIndexer.CallOpts, userIdHash)
}

// GetUserPasswordHash is a free data retrieval call binding the contract method 0x67835ade.
//
// Solidity: function getUserPasswordHash(address userAddr) view returns(bytes)
func (_EhrIndexer *EhrIndexerCaller) GetUserPasswordHash(opts *bind.CallOpts, userAddr common.Address) ([]byte, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getUserPasswordHash", userAddr)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GetUserPasswordHash is a free data retrieval call binding the contract method 0x67835ade.
//
// Solidity: function getUserPasswordHash(address userAddr) view returns(bytes)
func (_EhrIndexer *EhrIndexerSession) GetUserPasswordHash(userAddr common.Address) ([]byte, error) {
	return _EhrIndexer.Contract.GetUserPasswordHash(&_EhrIndexer.CallOpts, userAddr)
}

// GetUserPasswordHash is a free data retrieval call binding the contract method 0x67835ade.
//
// Solidity: function getUserPasswordHash(address userAddr) view returns(bytes)
func (_EhrIndexer *EhrIndexerCallerSession) GetUserPasswordHash(userAddr common.Address) ([]byte, error) {
	return _EhrIndexer.Contract.GetUserPasswordHash(&_EhrIndexer.CallOpts, userAddr)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_EhrIndexer *EhrIndexerCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "nonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_EhrIndexer *EhrIndexerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _EhrIndexer.Contract.Nonces(&_EhrIndexer.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_EhrIndexer *EhrIndexerCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _EhrIndexer.Contract.Nonces(&_EhrIndexer.CallOpts, arg0)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_EhrIndexer *EhrIndexerCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_EhrIndexer *EhrIndexerSession) Owner() (common.Address, error) {
	return _EhrIndexer.Contract.Owner(&_EhrIndexer.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_EhrIndexer *EhrIndexerCallerSession) Owner() (common.Address, error) {
	return _EhrIndexer.Contract.Owner(&_EhrIndexer.CallOpts)
}

// Users is a free data retrieval call binding the contract method 0xa87430ba.
//
// Solidity: function users(address ) view returns(bytes32 id, bytes32 systemID, uint8 role, bytes pwdHash)
func (_EhrIndexer *EhrIndexerCaller) Users(opts *bind.CallOpts, arg0 common.Address) (struct {
	Id       [32]byte
	SystemID [32]byte
	Role     uint8
	PwdHash  []byte
}, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "users", arg0)

	outstruct := new(struct {
		Id       [32]byte
		SystemID [32]byte
		Role     uint8
		PwdHash  []byte
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Id = *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)
	outstruct.SystemID = *abi.ConvertType(out[1], new([32]byte)).(*[32]byte)
	outstruct.Role = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.PwdHash = *abi.ConvertType(out[3], new([]byte)).(*[]byte)

	return *outstruct, err

}

// Users is a free data retrieval call binding the contract method 0xa87430ba.
//
// Solidity: function users(address ) view returns(bytes32 id, bytes32 systemID, uint8 role, bytes pwdHash)
func (_EhrIndexer *EhrIndexerSession) Users(arg0 common.Address) (struct {
	Id       [32]byte
	SystemID [32]byte
	Role     uint8
	PwdHash  []byte
}, error) {
	return _EhrIndexer.Contract.Users(&_EhrIndexer.CallOpts, arg0)
}

// Users is a free data retrieval call binding the contract method 0xa87430ba.
//
// Solidity: function users(address ) view returns(bytes32 id, bytes32 systemID, uint8 role, bytes pwdHash)
func (_EhrIndexer *EhrIndexerCallerSession) Users(arg0 common.Address) (struct {
	Id       [32]byte
	SystemID [32]byte
	Role     uint8
	PwdHash  []byte
}, error) {
	return _EhrIndexer.Contract.Users(&_EhrIndexer.CallOpts, arg0)
}

// AddEhrDoc is a paid mutator transaction binding the contract method 0xc123bffb.
//
// Solidity: function addEhrDoc((bytes32,(uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32),bytes,bytes,uint256,address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerTransactor) AddEhrDoc(opts *bind.TransactOpts, p EhrDocsAddEhrDocParams) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "addEhrDoc", p)
}

// AddEhrDoc is a paid mutator transaction binding the contract method 0xc123bffb.
//
// Solidity: function addEhrDoc((bytes32,(uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32),bytes,bytes,uint256,address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerSession) AddEhrDoc(p EhrDocsAddEhrDocParams) (*types.Transaction, error) {
	return _EhrIndexer.Contract.AddEhrDoc(&_EhrIndexer.TransactOpts, p)
}

// AddEhrDoc is a paid mutator transaction binding the contract method 0xc123bffb.
//
// Solidity: function addEhrDoc((bytes32,(uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32),bytes,bytes,uint256,address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) AddEhrDoc(p EhrDocsAddEhrDocParams) (*types.Transaction, error) {
	return _EhrIndexer.Contract.AddEhrDoc(&_EhrIndexer.TransactOpts, p)
}

// DeleteDoc is a paid mutator transaction binding the contract method 0xe9a330cc.
//
// Solidity: function deleteDoc(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version) returns()
func (_EhrIndexer *EhrIndexerTransactor) DeleteDoc(opts *bind.TransactOpts, ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "deleteDoc", ehrId, docType, docBaseUIDHash, version)
}

// DeleteDoc is a paid mutator transaction binding the contract method 0xe9a330cc.
//
// Solidity: function deleteDoc(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version) returns()
func (_EhrIndexer *EhrIndexerSession) DeleteDoc(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.DeleteDoc(&_EhrIndexer.TransactOpts, ehrId, docType, docBaseUIDHash, version)
}

// DeleteDoc is a paid mutator transaction binding the contract method 0xe9a330cc.
//
// Solidity: function deleteDoc(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) DeleteDoc(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.DeleteDoc(&_EhrIndexer.TransactOpts, ehrId, docType, docBaseUIDHash, version)
}

// GroupAddUser is a paid mutator transaction binding the contract method 0xc8d37fd6.
//
// Solidity: function groupAddUser((bytes32,address,uint8,bytes,bytes,uint256,address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerTransactor) GroupAddUser(opts *bind.TransactOpts, p UsersGroupAddUserParams) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "groupAddUser", p)
}

// GroupAddUser is a paid mutator transaction binding the contract method 0xc8d37fd6.
//
// Solidity: function groupAddUser((bytes32,address,uint8,bytes,bytes,uint256,address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerSession) GroupAddUser(p UsersGroupAddUserParams) (*types.Transaction, error) {
	return _EhrIndexer.Contract.GroupAddUser(&_EhrIndexer.TransactOpts, p)
}

// GroupAddUser is a paid mutator transaction binding the contract method 0xc8d37fd6.
//
// Solidity: function groupAddUser((bytes32,address,uint8,bytes,bytes,uint256,address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) GroupAddUser(p UsersGroupAddUserParams) (*types.Transaction, error) {
	return _EhrIndexer.Contract.GroupAddUser(&_EhrIndexer.TransactOpts, p)
}

// GroupRemoveUser is a paid mutator transaction binding the contract method 0x7ca1ed44.
//
// Solidity: function groupRemoveUser(bytes32 groupIdHash, address removingUserAddr, uint256 nonce, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactor) GroupRemoveUser(opts *bind.TransactOpts, groupIdHash [32]byte, removingUserAddr common.Address, nonce *big.Int, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "groupRemoveUser", groupIdHash, removingUserAddr, nonce, signer, signature)
}

// GroupRemoveUser is a paid mutator transaction binding the contract method 0x7ca1ed44.
//
// Solidity: function groupRemoveUser(bytes32 groupIdHash, address removingUserAddr, uint256 nonce, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerSession) GroupRemoveUser(groupIdHash [32]byte, removingUserAddr common.Address, nonce *big.Int, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.GroupRemoveUser(&_EhrIndexer.TransactOpts, groupIdHash, removingUserAddr, nonce, signer, signature)
}

// GroupRemoveUser is a paid mutator transaction binding the contract method 0x7ca1ed44.
//
// Solidity: function groupRemoveUser(bytes32 groupIdHash, address removingUserAddr, uint256 nonce, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) GroupRemoveUser(groupIdHash [32]byte, removingUserAddr common.Address, nonce *big.Int, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.GroupRemoveUser(&_EhrIndexer.TransactOpts, groupIdHash, removingUserAddr, nonce, signer, signature)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_EhrIndexer *EhrIndexerTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_EhrIndexer *EhrIndexerSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.Multicall(&_EhrIndexer.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_EhrIndexer *EhrIndexerTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.Multicall(&_EhrIndexer.TransactOpts, data)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_EhrIndexer *EhrIndexerTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_EhrIndexer *EhrIndexerSession) RenounceOwnership() (*types.Transaction, error) {
	return _EhrIndexer.Contract.RenounceOwnership(&_EhrIndexer.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_EhrIndexer *EhrIndexerTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _EhrIndexer.Contract.RenounceOwnership(&_EhrIndexer.TransactOpts)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_EhrIndexer *EhrIndexerTransactor) SetAllowed(opts *bind.TransactOpts, addr common.Address, allowed bool) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setAllowed", addr, allowed)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_EhrIndexer *EhrIndexerSession) SetAllowed(addr common.Address, allowed bool) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetAllowed(&_EhrIndexer.TransactOpts, addr, allowed)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) SetAllowed(addr common.Address, allowed bool) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetAllowed(&_EhrIndexer.TransactOpts, addr, allowed)
}

// SetDocAccess is a paid mutator transaction binding the contract method 0x4b492ad2.
//
// Solidity: function setDocAccess(bytes CID, (bytes32,bytes,bytes,uint8) accessObj, address userAddr, uint256 nonce, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactor) SetDocAccess(opts *bind.TransactOpts, CID []byte, accessObj AccessObject, userAddr common.Address, nonce *big.Int, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setDocAccess", CID, accessObj, userAddr, nonce, signer, signature)
}

// SetDocAccess is a paid mutator transaction binding the contract method 0x4b492ad2.
//
// Solidity: function setDocAccess(bytes CID, (bytes32,bytes,bytes,uint8) accessObj, address userAddr, uint256 nonce, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerSession) SetDocAccess(CID []byte, accessObj AccessObject, userAddr common.Address, nonce *big.Int, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetDocAccess(&_EhrIndexer.TransactOpts, CID, accessObj, userAddr, nonce, signer, signature)
}

// SetDocAccess is a paid mutator transaction binding the contract method 0x4b492ad2.
//
// Solidity: function setDocAccess(bytes CID, (bytes32,bytes,bytes,uint8) accessObj, address userAddr, uint256 nonce, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) SetDocAccess(CID []byte, accessObj AccessObject, userAddr common.Address, nonce *big.Int, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetDocAccess(&_EhrIndexer.TransactOpts, CID, accessObj, userAddr, nonce, signer, signature)
}

// SetEhrSubject is a paid mutator transaction binding the contract method 0xd7cb67f2.
//
// Solidity: function setEhrSubject(bytes32 subjectKey, bytes32 ehrId, uint256 nonce, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactor) SetEhrSubject(opts *bind.TransactOpts, subjectKey [32]byte, ehrId [32]byte, nonce *big.Int, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setEhrSubject", subjectKey, ehrId, nonce, signer, signature)
}

// SetEhrSubject is a paid mutator transaction binding the contract method 0xd7cb67f2.
//
// Solidity: function setEhrSubject(bytes32 subjectKey, bytes32 ehrId, uint256 nonce, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerSession) SetEhrSubject(subjectKey [32]byte, ehrId [32]byte, nonce *big.Int, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrSubject(&_EhrIndexer.TransactOpts, subjectKey, ehrId, nonce, signer, signature)
}

// SetEhrSubject is a paid mutator transaction binding the contract method 0xd7cb67f2.
//
// Solidity: function setEhrSubject(bytes32 subjectKey, bytes32 ehrId, uint256 nonce, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) SetEhrSubject(subjectKey [32]byte, ehrId [32]byte, nonce *big.Int, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrSubject(&_EhrIndexer.TransactOpts, subjectKey, ehrId, nonce, signer, signature)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0xeebfc00a.
//
// Solidity: function setEhrUser(bytes32 userId, bytes32 ehrId, uint256 nonce, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactor) SetEhrUser(opts *bind.TransactOpts, userId [32]byte, ehrId [32]byte, nonce *big.Int, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setEhrUser", userId, ehrId, nonce, signer, signature)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0xeebfc00a.
//
// Solidity: function setEhrUser(bytes32 userId, bytes32 ehrId, uint256 nonce, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerSession) SetEhrUser(userId [32]byte, ehrId [32]byte, nonce *big.Int, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrUser(&_EhrIndexer.TransactOpts, userId, ehrId, nonce, signer, signature)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0xeebfc00a.
//
// Solidity: function setEhrUser(bytes32 userId, bytes32 ehrId, uint256 nonce, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) SetEhrUser(userId [32]byte, ehrId [32]byte, nonce *big.Int, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrUser(&_EhrIndexer.TransactOpts, userId, ehrId, nonce, signer, signature)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_EhrIndexer *EhrIndexerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_EhrIndexer *EhrIndexerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _EhrIndexer.Contract.TransferOwnership(&_EhrIndexer.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _EhrIndexer.Contract.TransferOwnership(&_EhrIndexer.TransactOpts, newOwner)
}

// UserGroupCreate is a paid mutator transaction binding the contract method 0xc284aec1.
//
// Solidity: function userGroupCreate((bytes32,bytes,bytes,(bytes32,bytes)[],uint256,address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerTransactor) UserGroupCreate(opts *bind.TransactOpts, p UsersUserGroupCreateParams) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "userGroupCreate", p)
}

// UserGroupCreate is a paid mutator transaction binding the contract method 0xc284aec1.
//
// Solidity: function userGroupCreate((bytes32,bytes,bytes,(bytes32,bytes)[],uint256,address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerSession) UserGroupCreate(p UsersUserGroupCreateParams) (*types.Transaction, error) {
	return _EhrIndexer.Contract.UserGroupCreate(&_EhrIndexer.TransactOpts, p)
}

// UserGroupCreate is a paid mutator transaction binding the contract method 0xc284aec1.
//
// Solidity: function userGroupCreate((bytes32,bytes,bytes,(bytes32,bytes)[],uint256,address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) UserGroupCreate(p UsersUserGroupCreateParams) (*types.Transaction, error) {
	return _EhrIndexer.Contract.UserGroupCreate(&_EhrIndexer.TransactOpts, p)
}

// UserNew is a paid mutator transaction binding the contract method 0x94aacbad.
//
// Solidity: function userNew(address userAddr, bytes32 id, bytes32 systemID, uint8 role, bytes pwdHash, uint256 nonce, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactor) UserNew(opts *bind.TransactOpts, userAddr common.Address, id [32]byte, systemID [32]byte, role uint8, pwdHash []byte, nonce *big.Int, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "userNew", userAddr, id, systemID, role, pwdHash, nonce, signer, signature)
}

// UserNew is a paid mutator transaction binding the contract method 0x94aacbad.
//
// Solidity: function userNew(address userAddr, bytes32 id, bytes32 systemID, uint8 role, bytes pwdHash, uint256 nonce, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerSession) UserNew(userAddr common.Address, id [32]byte, systemID [32]byte, role uint8, pwdHash []byte, nonce *big.Int, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.UserNew(&_EhrIndexer.TransactOpts, userAddr, id, systemID, role, pwdHash, nonce, signer, signature)
}

// UserNew is a paid mutator transaction binding the contract method 0x94aacbad.
//
// Solidity: function userNew(address userAddr, bytes32 id, bytes32 systemID, uint8 role, bytes pwdHash, uint256 nonce, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) UserNew(userAddr common.Address, id [32]byte, systemID [32]byte, role uint8, pwdHash []byte, nonce *big.Int, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.UserNew(&_EhrIndexer.TransactOpts, userAddr, id, systemID, role, pwdHash, nonce, signer, signature)
}

// EhrIndexerOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the EhrIndexer contract.
type EhrIndexerOwnershipTransferredIterator struct {
	Event *EhrIndexerOwnershipTransferred // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *EhrIndexerOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EhrIndexerOwnershipTransferred)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(EhrIndexerOwnershipTransferred)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *EhrIndexerOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EhrIndexerOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EhrIndexerOwnershipTransferred represents a OwnershipTransferred event raised by the EhrIndexer contract.
type EhrIndexerOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_EhrIndexer *EhrIndexerFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*EhrIndexerOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _EhrIndexer.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &EhrIndexerOwnershipTransferredIterator{contract: _EhrIndexer.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_EhrIndexer *EhrIndexerFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *EhrIndexerOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _EhrIndexer.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EhrIndexerOwnershipTransferred)
				if err := _EhrIndexer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseOwnershipTransferred is a log parse operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_EhrIndexer *EhrIndexerFilterer) ParseOwnershipTransferred(log types.Log) (*EhrIndexerOwnershipTransferred, error) {
	event := new(EhrIndexerOwnershipTransferred)
	if err := _EhrIndexer.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
