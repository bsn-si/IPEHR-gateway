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

// EhrIndexerDocumentMeta is an auto generated low-level Go binding around an user-defined struct.
type EhrIndexerDocumentMeta struct {
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

// EhrIndexerMetaData contains all meta data concerning the EhrIndexer contract.
var EhrIndexerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"access\",\"type\":\"bytes\"}],\"name\":\"DocAccessChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"CID\",\"type\":\"bytes\"}],\"name\":\"EhrDocAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"subjectKey\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"}],\"name\":\"EhrSubjectSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"access\",\"type\":\"bytes\"}],\"name\":\"GroupAccessChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"enumEhrIndexer.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"enumEhrIndexer.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"CID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"dealCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"minerAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"docUIDEncrypted\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structEhrIndexer.DocumentMeta\",\"name\":\"docMeta\",\"type\":\"tuple\"}],\"name\":\"addEhrDoc\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedChange\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"dataSearch\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"nodeType\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"nodeID\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumEhrIndexer.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"}],\"name\":\"deleteDoc\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"docAccess\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"enumEhrIndexer.DocType\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ehrDocs\",\"outputs\":[{\"internalType\":\"enumEhrIndexer.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"enumEhrIndexer.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"CID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"dealCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"minerAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"docUIDEncrypted\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"ehrSubject\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"ehrUsers\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumEhrIndexer.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"name\":\"getDocByTime\",\"outputs\":[{\"components\":[{\"internalType\":\"enumEhrIndexer.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"enumEhrIndexer.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"CID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"dealCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"minerAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"docUIDEncrypted\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structEhrIndexer.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumEhrIndexer.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"}],\"name\":\"getDocByVersion\",\"outputs\":[{\"components\":[{\"internalType\":\"enumEhrIndexer.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"enumEhrIndexer.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"CID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"dealCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"minerAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"docUIDEncrypted\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structEhrIndexer.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumEhrIndexer.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"}],\"name\":\"getDocLastByBaseID\",\"outputs\":[{\"components\":[{\"internalType\":\"enumEhrIndexer.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"enumEhrIndexer.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"CID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"dealCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"minerAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"docUIDEncrypted\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structEhrIndexer.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumEhrIndexer.DocType\",\"name\":\"docType\",\"type\":\"uint8\"}],\"name\":\"getEhrDocs\",\"outputs\":[{\"components\":[{\"internalType\":\"enumEhrIndexer.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"enumEhrIndexer.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"CID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"dealCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"minerAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"docUIDEncrypted\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structEhrIndexer.DocumentMeta[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumEhrIndexer.DocType\",\"name\":\"docType\",\"type\":\"uint8\"}],\"name\":\"getLastEhrDocByType\",\"outputs\":[{\"components\":[{\"internalType\":\"enumEhrIndexer.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"enumEhrIndexer.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"CID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"dealCID\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"minerAddress\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"docUIDEncrypted\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structEhrIndexer.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"groupAccess\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAllowed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"access\",\"type\":\"bytes\"}],\"name\":\"setDocAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"subjectKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"}],\"name\":\"setEhrSubject\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"}],\"name\":\"setEhrUser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"key\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"access\",\"type\":\"bytes\"}],\"name\":\"setGroupAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// DocAccess is a free data retrieval call binding the contract method 0x4b5a743b.
//
// Solidity: function docAccess(bytes32 ) view returns(bytes)
func (_EhrIndexer *EhrIndexerCaller) DocAccess(opts *bind.CallOpts, arg0 [32]byte) ([]byte, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "docAccess", arg0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// DocAccess is a free data retrieval call binding the contract method 0x4b5a743b.
//
// Solidity: function docAccess(bytes32 ) view returns(bytes)
func (_EhrIndexer *EhrIndexerSession) DocAccess(arg0 [32]byte) ([]byte, error) {
	return _EhrIndexer.Contract.DocAccess(&_EhrIndexer.CallOpts, arg0)
}

// DocAccess is a free data retrieval call binding the contract method 0x4b5a743b.
//
// Solidity: function docAccess(bytes32 ) view returns(bytes)
func (_EhrIndexer *EhrIndexerCallerSession) DocAccess(arg0 [32]byte) ([]byte, error) {
	return _EhrIndexer.Contract.DocAccess(&_EhrIndexer.CallOpts, arg0)
}

// EhrDocs is a free data retrieval call binding the contract method 0xa14b8188.
//
// Solidity: function ehrDocs(bytes32 , uint8 , uint256 ) view returns(uint8 docType, uint8 status, bytes CID, bytes dealCID, bytes minerAddress, bytes docUIDEncrypted, bytes32 docBaseUIDHash, bytes32 version, bool isLast, uint32 timestamp)
func (_EhrIndexer *EhrIndexerCaller) EhrDocs(opts *bind.CallOpts, arg0 [32]byte, arg1 uint8, arg2 *big.Int) (struct {
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
}, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "ehrDocs", arg0, arg1, arg2)

	outstruct := new(struct {
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
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.DocType = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.Status = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.CID = *abi.ConvertType(out[2], new([]byte)).(*[]byte)
	outstruct.DealCID = *abi.ConvertType(out[3], new([]byte)).(*[]byte)
	outstruct.MinerAddress = *abi.ConvertType(out[4], new([]byte)).(*[]byte)
	outstruct.DocUIDEncrypted = *abi.ConvertType(out[5], new([]byte)).(*[]byte)
	outstruct.DocBaseUIDHash = *abi.ConvertType(out[6], new([32]byte)).(*[32]byte)
	outstruct.Version = *abi.ConvertType(out[7], new([32]byte)).(*[32]byte)
	outstruct.IsLast = *abi.ConvertType(out[8], new(bool)).(*bool)
	outstruct.Timestamp = *abi.ConvertType(out[9], new(uint32)).(*uint32)

	return *outstruct, err

}

// EhrDocs is a free data retrieval call binding the contract method 0xa14b8188.
//
// Solidity: function ehrDocs(bytes32 , uint8 , uint256 ) view returns(uint8 docType, uint8 status, bytes CID, bytes dealCID, bytes minerAddress, bytes docUIDEncrypted, bytes32 docBaseUIDHash, bytes32 version, bool isLast, uint32 timestamp)
func (_EhrIndexer *EhrIndexerSession) EhrDocs(arg0 [32]byte, arg1 uint8, arg2 *big.Int) (struct {
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
}, error) {
	return _EhrIndexer.Contract.EhrDocs(&_EhrIndexer.CallOpts, arg0, arg1, arg2)
}

// EhrDocs is a free data retrieval call binding the contract method 0xa14b8188.
//
// Solidity: function ehrDocs(bytes32 , uint8 , uint256 ) view returns(uint8 docType, uint8 status, bytes CID, bytes dealCID, bytes minerAddress, bytes docUIDEncrypted, bytes32 docBaseUIDHash, bytes32 version, bool isLast, uint32 timestamp)
func (_EhrIndexer *EhrIndexerCallerSession) EhrDocs(arg0 [32]byte, arg1 uint8, arg2 *big.Int) (struct {
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
}, error) {
	return _EhrIndexer.Contract.EhrDocs(&_EhrIndexer.CallOpts, arg0, arg1, arg2)
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

// GetDocByTime is a free data retrieval call binding the contract method 0x4f722b16.
//
// Solidity: function getDocByTime(bytes32 ehrId, uint8 docType, uint32 timestamp) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerCaller) GetDocByTime(opts *bind.CallOpts, ehrId [32]byte, docType uint8, timestamp uint32) (EhrIndexerDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getDocByTime", ehrId, docType, timestamp)

	if err != nil {
		return *new(EhrIndexerDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(EhrIndexerDocumentMeta)).(*EhrIndexerDocumentMeta)

	return out0, err

}

// GetDocByTime is a free data retrieval call binding the contract method 0x4f722b16.
//
// Solidity: function getDocByTime(bytes32 ehrId, uint8 docType, uint32 timestamp) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerSession) GetDocByTime(ehrId [32]byte, docType uint8, timestamp uint32) (EhrIndexerDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocByTime(&_EhrIndexer.CallOpts, ehrId, docType, timestamp)
}

// GetDocByTime is a free data retrieval call binding the contract method 0x4f722b16.
//
// Solidity: function getDocByTime(bytes32 ehrId, uint8 docType, uint32 timestamp) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerCallerSession) GetDocByTime(ehrId [32]byte, docType uint8, timestamp uint32) (EhrIndexerDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocByTime(&_EhrIndexer.CallOpts, ehrId, docType, timestamp)
}

// GetDocByVersion is a free data retrieval call binding the contract method 0x179fcacf.
//
// Solidity: function getDocByVersion(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerCaller) GetDocByVersion(opts *bind.CallOpts, ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte) (EhrIndexerDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getDocByVersion", ehrId, docType, docBaseUIDHash, version)

	if err != nil {
		return *new(EhrIndexerDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(EhrIndexerDocumentMeta)).(*EhrIndexerDocumentMeta)

	return out0, err

}

// GetDocByVersion is a free data retrieval call binding the contract method 0x179fcacf.
//
// Solidity: function getDocByVersion(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerSession) GetDocByVersion(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte) (EhrIndexerDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocByVersion(&_EhrIndexer.CallOpts, ehrId, docType, docBaseUIDHash, version)
}

// GetDocByVersion is a free data retrieval call binding the contract method 0x179fcacf.
//
// Solidity: function getDocByVersion(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerCallerSession) GetDocByVersion(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte) (EhrIndexerDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocByVersion(&_EhrIndexer.CallOpts, ehrId, docType, docBaseUIDHash, version)
}

// GetDocLastByBaseID is a free data retrieval call binding the contract method 0x949c8e77.
//
// Solidity: function getDocLastByBaseID(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerCaller) GetDocLastByBaseID(opts *bind.CallOpts, ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte) (EhrIndexerDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getDocLastByBaseID", ehrId, docType, docBaseUIDHash)

	if err != nil {
		return *new(EhrIndexerDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(EhrIndexerDocumentMeta)).(*EhrIndexerDocumentMeta)

	return out0, err

}

// GetDocLastByBaseID is a free data retrieval call binding the contract method 0x949c8e77.
//
// Solidity: function getDocLastByBaseID(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerSession) GetDocLastByBaseID(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte) (EhrIndexerDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocLastByBaseID(&_EhrIndexer.CallOpts, ehrId, docType, docBaseUIDHash)
}

// GetDocLastByBaseID is a free data retrieval call binding the contract method 0x949c8e77.
//
// Solidity: function getDocLastByBaseID(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerCallerSession) GetDocLastByBaseID(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte) (EhrIndexerDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocLastByBaseID(&_EhrIndexer.CallOpts, ehrId, docType, docBaseUIDHash)
}

// GetEhrDocs is a free data retrieval call binding the contract method 0xeec6b331.
//
// Solidity: function getEhrDocs(bytes32 ehrId, uint8 docType) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32)[])
func (_EhrIndexer *EhrIndexerCaller) GetEhrDocs(opts *bind.CallOpts, ehrId [32]byte, docType uint8) ([]EhrIndexerDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getEhrDocs", ehrId, docType)

	if err != nil {
		return *new([]EhrIndexerDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new([]EhrIndexerDocumentMeta)).(*[]EhrIndexerDocumentMeta)

	return out0, err

}

// GetEhrDocs is a free data retrieval call binding the contract method 0xeec6b331.
//
// Solidity: function getEhrDocs(bytes32 ehrId, uint8 docType) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32)[])
func (_EhrIndexer *EhrIndexerSession) GetEhrDocs(ehrId [32]byte, docType uint8) ([]EhrIndexerDocumentMeta, error) {
	return _EhrIndexer.Contract.GetEhrDocs(&_EhrIndexer.CallOpts, ehrId, docType)
}

// GetEhrDocs is a free data retrieval call binding the contract method 0xeec6b331.
//
// Solidity: function getEhrDocs(bytes32 ehrId, uint8 docType) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32)[])
func (_EhrIndexer *EhrIndexerCallerSession) GetEhrDocs(ehrId [32]byte, docType uint8) ([]EhrIndexerDocumentMeta, error) {
	return _EhrIndexer.Contract.GetEhrDocs(&_EhrIndexer.CallOpts, ehrId, docType)
}

// GetLastEhrDocByType is a free data retrieval call binding the contract method 0x15dbcf1a.
//
// Solidity: function getLastEhrDocByType(bytes32 ehrId, uint8 docType) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerCaller) GetLastEhrDocByType(opts *bind.CallOpts, ehrId [32]byte, docType uint8) (EhrIndexerDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getLastEhrDocByType", ehrId, docType)

	if err != nil {
		return *new(EhrIndexerDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(EhrIndexerDocumentMeta)).(*EhrIndexerDocumentMeta)

	return out0, err

}

// GetLastEhrDocByType is a free data retrieval call binding the contract method 0x15dbcf1a.
//
// Solidity: function getLastEhrDocByType(bytes32 ehrId, uint8 docType) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerSession) GetLastEhrDocByType(ehrId [32]byte, docType uint8) (EhrIndexerDocumentMeta, error) {
	return _EhrIndexer.Contract.GetLastEhrDocByType(&_EhrIndexer.CallOpts, ehrId, docType)
}

// GetLastEhrDocByType is a free data retrieval call binding the contract method 0x15dbcf1a.
//
// Solidity: function getLastEhrDocByType(bytes32 ehrId, uint8 docType) view returns((uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32))
func (_EhrIndexer *EhrIndexerCallerSession) GetLastEhrDocByType(ehrId [32]byte, docType uint8) (EhrIndexerDocumentMeta, error) {
	return _EhrIndexer.Contract.GetLastEhrDocByType(&_EhrIndexer.CallOpts, ehrId, docType)
}

// GroupAccess is a free data retrieval call binding the contract method 0xe51ad625.
//
// Solidity: function groupAccess(bytes32 ) view returns(bytes)
func (_EhrIndexer *EhrIndexerCaller) GroupAccess(opts *bind.CallOpts, arg0 [32]byte) ([]byte, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "groupAccess", arg0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// GroupAccess is a free data retrieval call binding the contract method 0xe51ad625.
//
// Solidity: function groupAccess(bytes32 ) view returns(bytes)
func (_EhrIndexer *EhrIndexerSession) GroupAccess(arg0 [32]byte) ([]byte, error) {
	return _EhrIndexer.Contract.GroupAccess(&_EhrIndexer.CallOpts, arg0)
}

// GroupAccess is a free data retrieval call binding the contract method 0xe51ad625.
//
// Solidity: function groupAccess(bytes32 ) view returns(bytes)
func (_EhrIndexer *EhrIndexerCallerSession) GroupAccess(arg0 [32]byte) ([]byte, error) {
	return _EhrIndexer.Contract.GroupAccess(&_EhrIndexer.CallOpts, arg0)
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

// AddEhrDoc is a paid mutator transaction binding the contract method 0xf69ffe32.
//
// Solidity: function addEhrDoc(bytes32 ehrId, (uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32) docMeta) returns()
func (_EhrIndexer *EhrIndexerTransactor) AddEhrDoc(opts *bind.TransactOpts, ehrId [32]byte, docMeta EhrIndexerDocumentMeta) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "addEhrDoc", ehrId, docMeta)
}

// AddEhrDoc is a paid mutator transaction binding the contract method 0xf69ffe32.
//
// Solidity: function addEhrDoc(bytes32 ehrId, (uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32) docMeta) returns()
func (_EhrIndexer *EhrIndexerSession) AddEhrDoc(ehrId [32]byte, docMeta EhrIndexerDocumentMeta) (*types.Transaction, error) {
	return _EhrIndexer.Contract.AddEhrDoc(&_EhrIndexer.TransactOpts, ehrId, docMeta)
}

// AddEhrDoc is a paid mutator transaction binding the contract method 0xf69ffe32.
//
// Solidity: function addEhrDoc(bytes32 ehrId, (uint8,uint8,bytes,bytes,bytes,bytes,bytes32,bytes32,bool,uint32) docMeta) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) AddEhrDoc(ehrId [32]byte, docMeta EhrIndexerDocumentMeta) (*types.Transaction, error) {
	return _EhrIndexer.Contract.AddEhrDoc(&_EhrIndexer.TransactOpts, ehrId, docMeta)
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

// SetDocAccess is a paid mutator transaction binding the contract method 0x65bde879.
//
// Solidity: function setDocAccess(bytes32 key, bytes access) returns()
func (_EhrIndexer *EhrIndexerTransactor) SetDocAccess(opts *bind.TransactOpts, key [32]byte, access []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setDocAccess", key, access)
}

// SetDocAccess is a paid mutator transaction binding the contract method 0x65bde879.
//
// Solidity: function setDocAccess(bytes32 key, bytes access) returns()
func (_EhrIndexer *EhrIndexerSession) SetDocAccess(key [32]byte, access []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetDocAccess(&_EhrIndexer.TransactOpts, key, access)
}

// SetDocAccess is a paid mutator transaction binding the contract method 0x65bde879.
//
// Solidity: function setDocAccess(bytes32 key, bytes access) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) SetDocAccess(key [32]byte, access []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetDocAccess(&_EhrIndexer.TransactOpts, key, access)
}

// SetEhrSubject is a paid mutator transaction binding the contract method 0xfa0e2b69.
//
// Solidity: function setEhrSubject(bytes32 subjectKey, bytes32 ehrId) returns()
func (_EhrIndexer *EhrIndexerTransactor) SetEhrSubject(opts *bind.TransactOpts, subjectKey [32]byte, ehrId [32]byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setEhrSubject", subjectKey, ehrId)
}

// SetEhrSubject is a paid mutator transaction binding the contract method 0xfa0e2b69.
//
// Solidity: function setEhrSubject(bytes32 subjectKey, bytes32 ehrId) returns()
func (_EhrIndexer *EhrIndexerSession) SetEhrSubject(subjectKey [32]byte, ehrId [32]byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrSubject(&_EhrIndexer.TransactOpts, subjectKey, ehrId)
}

// SetEhrSubject is a paid mutator transaction binding the contract method 0xfa0e2b69.
//
// Solidity: function setEhrSubject(bytes32 subjectKey, bytes32 ehrId) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) SetEhrSubject(subjectKey [32]byte, ehrId [32]byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrSubject(&_EhrIndexer.TransactOpts, subjectKey, ehrId)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0xfbdbd13b.
//
// Solidity: function setEhrUser(bytes32 userId, bytes32 ehrId) returns()
func (_EhrIndexer *EhrIndexerTransactor) SetEhrUser(opts *bind.TransactOpts, userId [32]byte, ehrId [32]byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setEhrUser", userId, ehrId)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0xfbdbd13b.
//
// Solidity: function setEhrUser(bytes32 userId, bytes32 ehrId) returns()
func (_EhrIndexer *EhrIndexerSession) SetEhrUser(userId [32]byte, ehrId [32]byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrUser(&_EhrIndexer.TransactOpts, userId, ehrId)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0xfbdbd13b.
//
// Solidity: function setEhrUser(bytes32 userId, bytes32 ehrId) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) SetEhrUser(userId [32]byte, ehrId [32]byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrUser(&_EhrIndexer.TransactOpts, userId, ehrId)
}

// SetGroupAccess is a paid mutator transaction binding the contract method 0xe40671ac.
//
// Solidity: function setGroupAccess(bytes32 key, bytes access) returns()
func (_EhrIndexer *EhrIndexerTransactor) SetGroupAccess(opts *bind.TransactOpts, key [32]byte, access []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setGroupAccess", key, access)
}

// SetGroupAccess is a paid mutator transaction binding the contract method 0xe40671ac.
//
// Solidity: function setGroupAccess(bytes32 key, bytes access) returns()
func (_EhrIndexer *EhrIndexerSession) SetGroupAccess(key [32]byte, access []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetGroupAccess(&_EhrIndexer.TransactOpts, key, access)
}

// SetGroupAccess is a paid mutator transaction binding the contract method 0xe40671ac.
//
// Solidity: function setGroupAccess(bytes32 key, bytes access) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) SetGroupAccess(key [32]byte, access []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetGroupAccess(&_EhrIndexer.TransactOpts, key, access)
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

// EhrIndexerDocAccessChangedIterator is returned from FilterDocAccessChanged and is used to iterate over the raw logs and unpacked data for DocAccessChanged events raised by the EhrIndexer contract.
type EhrIndexerDocAccessChangedIterator struct {
	Event *EhrIndexerDocAccessChanged // Event containing the contract specifics and raw log

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
func (it *EhrIndexerDocAccessChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EhrIndexerDocAccessChanged)
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
		it.Event = new(EhrIndexerDocAccessChanged)
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
func (it *EhrIndexerDocAccessChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EhrIndexerDocAccessChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EhrIndexerDocAccessChanged represents a DocAccessChanged event raised by the EhrIndexer contract.
type EhrIndexerDocAccessChanged struct {
	Key    [32]byte
	Access []byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDocAccessChanged is a free log retrieval operation binding the contract event 0x96a63afebe577b3909fe26d5e194804fa47790fc5d1fc6c80f31b67c69b878fb.
//
// Solidity: event DocAccessChanged(bytes32 key, bytes access)
func (_EhrIndexer *EhrIndexerFilterer) FilterDocAccessChanged(opts *bind.FilterOpts) (*EhrIndexerDocAccessChangedIterator, error) {

	logs, sub, err := _EhrIndexer.contract.FilterLogs(opts, "DocAccessChanged")
	if err != nil {
		return nil, err
	}
	return &EhrIndexerDocAccessChangedIterator{contract: _EhrIndexer.contract, event: "DocAccessChanged", logs: logs, sub: sub}, nil
}

// WatchDocAccessChanged is a free log subscription operation binding the contract event 0x96a63afebe577b3909fe26d5e194804fa47790fc5d1fc6c80f31b67c69b878fb.
//
// Solidity: event DocAccessChanged(bytes32 key, bytes access)
func (_EhrIndexer *EhrIndexerFilterer) WatchDocAccessChanged(opts *bind.WatchOpts, sink chan<- *EhrIndexerDocAccessChanged) (event.Subscription, error) {

	logs, sub, err := _EhrIndexer.contract.WatchLogs(opts, "DocAccessChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EhrIndexerDocAccessChanged)
				if err := _EhrIndexer.contract.UnpackLog(event, "DocAccessChanged", log); err != nil {
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

// ParseDocAccessChanged is a log parse operation binding the contract event 0x96a63afebe577b3909fe26d5e194804fa47790fc5d1fc6c80f31b67c69b878fb.
//
// Solidity: event DocAccessChanged(bytes32 key, bytes access)
func (_EhrIndexer *EhrIndexerFilterer) ParseDocAccessChanged(log types.Log) (*EhrIndexerDocAccessChanged, error) {
	event := new(EhrIndexerDocAccessChanged)
	if err := _EhrIndexer.contract.UnpackLog(event, "DocAccessChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EhrIndexerEhrDocAddedIterator is returned from FilterEhrDocAdded and is used to iterate over the raw logs and unpacked data for EhrDocAdded events raised by the EhrIndexer contract.
type EhrIndexerEhrDocAddedIterator struct {
	Event *EhrIndexerEhrDocAdded // Event containing the contract specifics and raw log

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
func (it *EhrIndexerEhrDocAddedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EhrIndexerEhrDocAdded)
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
		it.Event = new(EhrIndexerEhrDocAdded)
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
func (it *EhrIndexerEhrDocAddedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EhrIndexerEhrDocAddedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EhrIndexerEhrDocAdded represents a EhrDocAdded event raised by the EhrIndexer contract.
type EhrIndexerEhrDocAdded struct {
	EhrId [32]byte
	CID   []byte
	Raw   types.Log // Blockchain specific contextual infos
}

// FilterEhrDocAdded is a free log retrieval operation binding the contract event 0x0477812eb2952d83d4085b67b0669f6c20e2a052588fbf64e7d8113a8bfb2442.
//
// Solidity: event EhrDocAdded(bytes32 ehrId, bytes CID)
func (_EhrIndexer *EhrIndexerFilterer) FilterEhrDocAdded(opts *bind.FilterOpts) (*EhrIndexerEhrDocAddedIterator, error) {

	logs, sub, err := _EhrIndexer.contract.FilterLogs(opts, "EhrDocAdded")
	if err != nil {
		return nil, err
	}
	return &EhrIndexerEhrDocAddedIterator{contract: _EhrIndexer.contract, event: "EhrDocAdded", logs: logs, sub: sub}, nil
}

// WatchEhrDocAdded is a free log subscription operation binding the contract event 0x0477812eb2952d83d4085b67b0669f6c20e2a052588fbf64e7d8113a8bfb2442.
//
// Solidity: event EhrDocAdded(bytes32 ehrId, bytes CID)
func (_EhrIndexer *EhrIndexerFilterer) WatchEhrDocAdded(opts *bind.WatchOpts, sink chan<- *EhrIndexerEhrDocAdded) (event.Subscription, error) {

	logs, sub, err := _EhrIndexer.contract.WatchLogs(opts, "EhrDocAdded")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EhrIndexerEhrDocAdded)
				if err := _EhrIndexer.contract.UnpackLog(event, "EhrDocAdded", log); err != nil {
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

// ParseEhrDocAdded is a log parse operation binding the contract event 0x0477812eb2952d83d4085b67b0669f6c20e2a052588fbf64e7d8113a8bfb2442.
//
// Solidity: event EhrDocAdded(bytes32 ehrId, bytes CID)
func (_EhrIndexer *EhrIndexerFilterer) ParseEhrDocAdded(log types.Log) (*EhrIndexerEhrDocAdded, error) {
	event := new(EhrIndexerEhrDocAdded)
	if err := _EhrIndexer.contract.UnpackLog(event, "EhrDocAdded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EhrIndexerEhrSubjectSetIterator is returned from FilterEhrSubjectSet and is used to iterate over the raw logs and unpacked data for EhrSubjectSet events raised by the EhrIndexer contract.
type EhrIndexerEhrSubjectSetIterator struct {
	Event *EhrIndexerEhrSubjectSet // Event containing the contract specifics and raw log

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
func (it *EhrIndexerEhrSubjectSetIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EhrIndexerEhrSubjectSet)
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
		it.Event = new(EhrIndexerEhrSubjectSet)
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
func (it *EhrIndexerEhrSubjectSetIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EhrIndexerEhrSubjectSetIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EhrIndexerEhrSubjectSet represents a EhrSubjectSet event raised by the EhrIndexer contract.
type EhrIndexerEhrSubjectSet struct {
	SubjectKey [32]byte
	EhrId      [32]byte
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterEhrSubjectSet is a free log retrieval operation binding the contract event 0xa162c4f1c6de6f45639be0d9608d4603bc8ff3d7890ce493368734ee95710a5b.
//
// Solidity: event EhrSubjectSet(bytes32 subjectKey, bytes32 ehrId)
func (_EhrIndexer *EhrIndexerFilterer) FilterEhrSubjectSet(opts *bind.FilterOpts) (*EhrIndexerEhrSubjectSetIterator, error) {

	logs, sub, err := _EhrIndexer.contract.FilterLogs(opts, "EhrSubjectSet")
	if err != nil {
		return nil, err
	}
	return &EhrIndexerEhrSubjectSetIterator{contract: _EhrIndexer.contract, event: "EhrSubjectSet", logs: logs, sub: sub}, nil
}

// WatchEhrSubjectSet is a free log subscription operation binding the contract event 0xa162c4f1c6de6f45639be0d9608d4603bc8ff3d7890ce493368734ee95710a5b.
//
// Solidity: event EhrSubjectSet(bytes32 subjectKey, bytes32 ehrId)
func (_EhrIndexer *EhrIndexerFilterer) WatchEhrSubjectSet(opts *bind.WatchOpts, sink chan<- *EhrIndexerEhrSubjectSet) (event.Subscription, error) {

	logs, sub, err := _EhrIndexer.contract.WatchLogs(opts, "EhrSubjectSet")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EhrIndexerEhrSubjectSet)
				if err := _EhrIndexer.contract.UnpackLog(event, "EhrSubjectSet", log); err != nil {
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

// ParseEhrSubjectSet is a log parse operation binding the contract event 0xa162c4f1c6de6f45639be0d9608d4603bc8ff3d7890ce493368734ee95710a5b.
//
// Solidity: event EhrSubjectSet(bytes32 subjectKey, bytes32 ehrId)
func (_EhrIndexer *EhrIndexerFilterer) ParseEhrSubjectSet(log types.Log) (*EhrIndexerEhrSubjectSet, error) {
	event := new(EhrIndexerEhrSubjectSet)
	if err := _EhrIndexer.contract.UnpackLog(event, "EhrSubjectSet", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// EhrIndexerGroupAccessChangedIterator is returned from FilterGroupAccessChanged and is used to iterate over the raw logs and unpacked data for GroupAccessChanged events raised by the EhrIndexer contract.
type EhrIndexerGroupAccessChangedIterator struct {
	Event *EhrIndexerGroupAccessChanged // Event containing the contract specifics and raw log

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
func (it *EhrIndexerGroupAccessChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EhrIndexerGroupAccessChanged)
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
		it.Event = new(EhrIndexerGroupAccessChanged)
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
func (it *EhrIndexerGroupAccessChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EhrIndexerGroupAccessChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EhrIndexerGroupAccessChanged represents a GroupAccessChanged event raised by the EhrIndexer contract.
type EhrIndexerGroupAccessChanged struct {
	Key    [32]byte
	Access []byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterGroupAccessChanged is a free log retrieval operation binding the contract event 0x1ade2dc6533b425f330a611adaf78bb5cf637acdcdb83feecbf619232691da86.
//
// Solidity: event GroupAccessChanged(bytes32 key, bytes access)
func (_EhrIndexer *EhrIndexerFilterer) FilterGroupAccessChanged(opts *bind.FilterOpts) (*EhrIndexerGroupAccessChangedIterator, error) {

	logs, sub, err := _EhrIndexer.contract.FilterLogs(opts, "GroupAccessChanged")
	if err != nil {
		return nil, err
	}
	return &EhrIndexerGroupAccessChangedIterator{contract: _EhrIndexer.contract, event: "GroupAccessChanged", logs: logs, sub: sub}, nil
}

// WatchGroupAccessChanged is a free log subscription operation binding the contract event 0x1ade2dc6533b425f330a611adaf78bb5cf637acdcdb83feecbf619232691da86.
//
// Solidity: event GroupAccessChanged(bytes32 key, bytes access)
func (_EhrIndexer *EhrIndexerFilterer) WatchGroupAccessChanged(opts *bind.WatchOpts, sink chan<- *EhrIndexerGroupAccessChanged) (event.Subscription, error) {

	logs, sub, err := _EhrIndexer.contract.WatchLogs(opts, "GroupAccessChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EhrIndexerGroupAccessChanged)
				if err := _EhrIndexer.contract.UnpackLog(event, "GroupAccessChanged", log); err != nil {
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

// ParseGroupAccessChanged is a log parse operation binding the contract event 0x1ade2dc6533b425f330a611adaf78bb5cf637acdcdb83feecbf619232691da86.
//
// Solidity: event GroupAccessChanged(bytes32 key, bytes access)
func (_EhrIndexer *EhrIndexerFilterer) ParseGroupAccessChanged(log types.Log) (*EhrIndexerGroupAccessChanged, error) {
	event := new(EhrIndexerGroupAccessChanged)
	if err := _EhrIndexer.contract.UnpackLog(event, "GroupAccessChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
