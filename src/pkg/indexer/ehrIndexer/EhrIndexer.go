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
	DocType        uint8
	Status         uint8
	StorageId      *big.Int
	DocIdEncrypted []byte
	Timestamp      uint32
}

// EhrIndexerMetaData contains all meta data concerning the EhrIndexer contract.
var EhrIndexerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"userId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"access\",\"type\":\"bytes\"}],\"name\":\"DataAccessChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"userId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"access\",\"type\":\"bytes\"}],\"name\":\"DocAccessChanged\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ehrId\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"storageId\",\"type\":\"uint256\"}],\"name\":\"EhrDocAdded\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"subjectKey\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"ehrId\",\"type\":\"uint256\"}],\"name\":\"EhrSubjectSet\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"ehrId\",\"type\":\"uint256\"},{\"components\":[{\"internalType\":\"uint8\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"storageId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"docIdEncrypted\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structEhrIndexer.DocumentMeta\",\"name\":\"docMeta\",\"type\":\"tuple\"}],\"name\":\"addEhrDoc\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedChange\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"dataAccess\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"docAccess\",\"outputs\":[{\"internalType\":\"bytes\",\"name\":\"\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ehrDocs\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"storageId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"docIdEncrypted\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ehrSubject\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"name\":\"ehrUsers\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"ehrId\",\"type\":\"uint256\"}],\"name\":\"getEhrDocs\",\"outputs\":[{\"components\":[{\"internalType\":\"uint8\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"storageId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"docIdEncrypted\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"internalType\":\"structEhrIndexer.DocumentMeta[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAllowed\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"userId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_access\",\"type\":\"bytes\"}],\"name\":\"setDataAccess\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"userId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"_access\",\"type\":\"bytes\"}],\"name\":\"setDocAccess\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"subjectKey\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"_ehrId\",\"type\":\"uint256\"}],\"name\":\"setEhrSubject\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"userId\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ehrId\",\"type\":\"uint256\"}],\"name\":\"setEhrUser\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
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

// DataAccess is a free data retrieval call binding the contract method 0x47f5105d.
//
// Solidity: function dataAccess(uint256 ) view returns(bytes)
func (_EhrIndexer *EhrIndexerCaller) DataAccess(opts *bind.CallOpts, arg0 *big.Int) ([]byte, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "dataAccess", arg0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// DataAccess is a free data retrieval call binding the contract method 0x47f5105d.
//
// Solidity: function dataAccess(uint256 ) view returns(bytes)
func (_EhrIndexer *EhrIndexerSession) DataAccess(arg0 *big.Int) ([]byte, error) {
	return _EhrIndexer.Contract.DataAccess(&_EhrIndexer.CallOpts, arg0)
}

// DataAccess is a free data retrieval call binding the contract method 0x47f5105d.
//
// Solidity: function dataAccess(uint256 ) view returns(bytes)
func (_EhrIndexer *EhrIndexerCallerSession) DataAccess(arg0 *big.Int) ([]byte, error) {
	return _EhrIndexer.Contract.DataAccess(&_EhrIndexer.CallOpts, arg0)
}

// DocAccess is a free data retrieval call binding the contract method 0x8ae39fb9.
//
// Solidity: function docAccess(uint256 ) view returns(bytes)
func (_EhrIndexer *EhrIndexerCaller) DocAccess(opts *bind.CallOpts, arg0 *big.Int) ([]byte, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "docAccess", arg0)

	if err != nil {
		return *new([]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([]byte)).(*[]byte)

	return out0, err

}

// DocAccess is a free data retrieval call binding the contract method 0x8ae39fb9.
//
// Solidity: function docAccess(uint256 ) view returns(bytes)
func (_EhrIndexer *EhrIndexerSession) DocAccess(arg0 *big.Int) ([]byte, error) {
	return _EhrIndexer.Contract.DocAccess(&_EhrIndexer.CallOpts, arg0)
}

// DocAccess is a free data retrieval call binding the contract method 0x8ae39fb9.
//
// Solidity: function docAccess(uint256 ) view returns(bytes)
func (_EhrIndexer *EhrIndexerCallerSession) DocAccess(arg0 *big.Int) ([]byte, error) {
	return _EhrIndexer.Contract.DocAccess(&_EhrIndexer.CallOpts, arg0)
}

// EhrDocs is a free data retrieval call binding the contract method 0x4c4df71f.
//
// Solidity: function ehrDocs(uint256 , uint256 ) view returns(uint8 docType, uint8 status, uint256 storageId, bytes docIdEncrypted, uint32 timestamp)
func (_EhrIndexer *EhrIndexerCaller) EhrDocs(opts *bind.CallOpts, arg0 *big.Int, arg1 *big.Int) (struct {
	DocType        uint8
	Status         uint8
	StorageId      *big.Int
	DocIdEncrypted []byte
	Timestamp      uint32
}, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "ehrDocs", arg0, arg1)

	outstruct := new(struct {
		DocType        uint8
		Status         uint8
		StorageId      *big.Int
		DocIdEncrypted []byte
		Timestamp      uint32
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.DocType = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.Status = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.StorageId = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.DocIdEncrypted = *abi.ConvertType(out[3], new([]byte)).(*[]byte)
	outstruct.Timestamp = *abi.ConvertType(out[4], new(uint32)).(*uint32)

	return *outstruct, err

}

// EhrDocs is a free data retrieval call binding the contract method 0x4c4df71f.
//
// Solidity: function ehrDocs(uint256 , uint256 ) view returns(uint8 docType, uint8 status, uint256 storageId, bytes docIdEncrypted, uint32 timestamp)
func (_EhrIndexer *EhrIndexerSession) EhrDocs(arg0 *big.Int, arg1 *big.Int) (struct {
	DocType        uint8
	Status         uint8
	StorageId      *big.Int
	DocIdEncrypted []byte
	Timestamp      uint32
}, error) {
	return _EhrIndexer.Contract.EhrDocs(&_EhrIndexer.CallOpts, arg0, arg1)
}

// EhrDocs is a free data retrieval call binding the contract method 0x4c4df71f.
//
// Solidity: function ehrDocs(uint256 , uint256 ) view returns(uint8 docType, uint8 status, uint256 storageId, bytes docIdEncrypted, uint32 timestamp)
func (_EhrIndexer *EhrIndexerCallerSession) EhrDocs(arg0 *big.Int, arg1 *big.Int) (struct {
	DocType        uint8
	Status         uint8
	StorageId      *big.Int
	DocIdEncrypted []byte
	Timestamp      uint32
}, error) {
	return _EhrIndexer.Contract.EhrDocs(&_EhrIndexer.CallOpts, arg0, arg1)
}

// EhrSubject is a free data retrieval call binding the contract method 0xdbe3281f.
//
// Solidity: function ehrSubject(uint256 ) view returns(uint256)
func (_EhrIndexer *EhrIndexerCaller) EhrSubject(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "ehrSubject", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EhrSubject is a free data retrieval call binding the contract method 0xdbe3281f.
//
// Solidity: function ehrSubject(uint256 ) view returns(uint256)
func (_EhrIndexer *EhrIndexerSession) EhrSubject(arg0 *big.Int) (*big.Int, error) {
	return _EhrIndexer.Contract.EhrSubject(&_EhrIndexer.CallOpts, arg0)
}

// EhrSubject is a free data retrieval call binding the contract method 0xdbe3281f.
//
// Solidity: function ehrSubject(uint256 ) view returns(uint256)
func (_EhrIndexer *EhrIndexerCallerSession) EhrSubject(arg0 *big.Int) (*big.Int, error) {
	return _EhrIndexer.Contract.EhrSubject(&_EhrIndexer.CallOpts, arg0)
}

// EhrUsers is a free data retrieval call binding the contract method 0x2046667a.
//
// Solidity: function ehrUsers(uint256 ) view returns(uint256)
func (_EhrIndexer *EhrIndexerCaller) EhrUsers(opts *bind.CallOpts, arg0 *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "ehrUsers", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// EhrUsers is a free data retrieval call binding the contract method 0x2046667a.
//
// Solidity: function ehrUsers(uint256 ) view returns(uint256)
func (_EhrIndexer *EhrIndexerSession) EhrUsers(arg0 *big.Int) (*big.Int, error) {
	return _EhrIndexer.Contract.EhrUsers(&_EhrIndexer.CallOpts, arg0)
}

// EhrUsers is a free data retrieval call binding the contract method 0x2046667a.
//
// Solidity: function ehrUsers(uint256 ) view returns(uint256)
func (_EhrIndexer *EhrIndexerCallerSession) EhrUsers(arg0 *big.Int) (*big.Int, error) {
	return _EhrIndexer.Contract.EhrUsers(&_EhrIndexer.CallOpts, arg0)
}

// GetEhrDocs is a free data retrieval call binding the contract method 0xc0e06345.
//
// Solidity: function getEhrDocs(uint256 ehrId) view returns((uint8,uint8,uint256,bytes,uint32)[])
func (_EhrIndexer *EhrIndexerCaller) GetEhrDocs(opts *bind.CallOpts, ehrId *big.Int) ([]EhrIndexerDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getEhrDocs", ehrId)

	if err != nil {
		return *new([]EhrIndexerDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new([]EhrIndexerDocumentMeta)).(*[]EhrIndexerDocumentMeta)

	return out0, err

}

// GetEhrDocs is a free data retrieval call binding the contract method 0xc0e06345.
//
// Solidity: function getEhrDocs(uint256 ehrId) view returns((uint8,uint8,uint256,bytes,uint32)[])
func (_EhrIndexer *EhrIndexerSession) GetEhrDocs(ehrId *big.Int) ([]EhrIndexerDocumentMeta, error) {
	return _EhrIndexer.Contract.GetEhrDocs(&_EhrIndexer.CallOpts, ehrId)
}

// GetEhrDocs is a free data retrieval call binding the contract method 0xc0e06345.
//
// Solidity: function getEhrDocs(uint256 ehrId) view returns((uint8,uint8,uint256,bytes,uint32)[])
func (_EhrIndexer *EhrIndexerCallerSession) GetEhrDocs(ehrId *big.Int) ([]EhrIndexerDocumentMeta, error) {
	return _EhrIndexer.Contract.GetEhrDocs(&_EhrIndexer.CallOpts, ehrId)
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

// AddEhrDoc is a paid mutator transaction binding the contract method 0xe94feafc.
//
// Solidity: function addEhrDoc(uint256 ehrId, (uint8,uint8,uint256,bytes,uint32) docMeta) returns()
func (_EhrIndexer *EhrIndexerTransactor) AddEhrDoc(opts *bind.TransactOpts, ehrId *big.Int, docMeta EhrIndexerDocumentMeta) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "addEhrDoc", ehrId, docMeta)
}

// AddEhrDoc is a paid mutator transaction binding the contract method 0xe94feafc.
//
// Solidity: function addEhrDoc(uint256 ehrId, (uint8,uint8,uint256,bytes,uint32) docMeta) returns()
func (_EhrIndexer *EhrIndexerSession) AddEhrDoc(ehrId *big.Int, docMeta EhrIndexerDocumentMeta) (*types.Transaction, error) {
	return _EhrIndexer.Contract.AddEhrDoc(&_EhrIndexer.TransactOpts, ehrId, docMeta)
}

// AddEhrDoc is a paid mutator transaction binding the contract method 0xe94feafc.
//
// Solidity: function addEhrDoc(uint256 ehrId, (uint8,uint8,uint256,bytes,uint32) docMeta) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) AddEhrDoc(ehrId *big.Int, docMeta EhrIndexerDocumentMeta) (*types.Transaction, error) {
	return _EhrIndexer.Contract.AddEhrDoc(&_EhrIndexer.TransactOpts, ehrId, docMeta)
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
// Solidity: function setAllowed(address addr, bool allowed) returns(bool)
func (_EhrIndexer *EhrIndexerTransactor) SetAllowed(opts *bind.TransactOpts, addr common.Address, allowed bool) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setAllowed", addr, allowed)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns(bool)
func (_EhrIndexer *EhrIndexerSession) SetAllowed(addr common.Address, allowed bool) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetAllowed(&_EhrIndexer.TransactOpts, addr, allowed)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns(bool)
func (_EhrIndexer *EhrIndexerTransactorSession) SetAllowed(addr common.Address, allowed bool) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetAllowed(&_EhrIndexer.TransactOpts, addr, allowed)
}

// SetDataAccess is a paid mutator transaction binding the contract method 0x4bb5a1ba.
//
// Solidity: function setDataAccess(uint256 userId, bytes _access) returns(uint256)
func (_EhrIndexer *EhrIndexerTransactor) SetDataAccess(opts *bind.TransactOpts, userId *big.Int, _access []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setDataAccess", userId, _access)
}

// SetDataAccess is a paid mutator transaction binding the contract method 0x4bb5a1ba.
//
// Solidity: function setDataAccess(uint256 userId, bytes _access) returns(uint256)
func (_EhrIndexer *EhrIndexerSession) SetDataAccess(userId *big.Int, _access []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetDataAccess(&_EhrIndexer.TransactOpts, userId, _access)
}

// SetDataAccess is a paid mutator transaction binding the contract method 0x4bb5a1ba.
//
// Solidity: function setDataAccess(uint256 userId, bytes _access) returns(uint256)
func (_EhrIndexer *EhrIndexerTransactorSession) SetDataAccess(userId *big.Int, _access []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetDataAccess(&_EhrIndexer.TransactOpts, userId, _access)
}

// SetDocAccess is a paid mutator transaction binding the contract method 0x1efa95ce.
//
// Solidity: function setDocAccess(uint256 userId, bytes _access) returns(uint256)
func (_EhrIndexer *EhrIndexerTransactor) SetDocAccess(opts *bind.TransactOpts, userId *big.Int, _access []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setDocAccess", userId, _access)
}

// SetDocAccess is a paid mutator transaction binding the contract method 0x1efa95ce.
//
// Solidity: function setDocAccess(uint256 userId, bytes _access) returns(uint256)
func (_EhrIndexer *EhrIndexerSession) SetDocAccess(userId *big.Int, _access []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetDocAccess(&_EhrIndexer.TransactOpts, userId, _access)
}

// SetDocAccess is a paid mutator transaction binding the contract method 0x1efa95ce.
//
// Solidity: function setDocAccess(uint256 userId, bytes _access) returns(uint256)
func (_EhrIndexer *EhrIndexerTransactorSession) SetDocAccess(userId *big.Int, _access []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetDocAccess(&_EhrIndexer.TransactOpts, userId, _access)
}

// SetEhrSubject is a paid mutator transaction binding the contract method 0x46db77a1.
//
// Solidity: function setEhrSubject(uint256 subjectKey, uint256 _ehrId) returns(uint256)
func (_EhrIndexer *EhrIndexerTransactor) SetEhrSubject(opts *bind.TransactOpts, subjectKey *big.Int, _ehrId *big.Int) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setEhrSubject", subjectKey, _ehrId)
}

// SetEhrSubject is a paid mutator transaction binding the contract method 0x46db77a1.
//
// Solidity: function setEhrSubject(uint256 subjectKey, uint256 _ehrId) returns(uint256)
func (_EhrIndexer *EhrIndexerSession) SetEhrSubject(subjectKey *big.Int, _ehrId *big.Int) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrSubject(&_EhrIndexer.TransactOpts, subjectKey, _ehrId)
}

// SetEhrSubject is a paid mutator transaction binding the contract method 0x46db77a1.
//
// Solidity: function setEhrSubject(uint256 subjectKey, uint256 _ehrId) returns(uint256)
func (_EhrIndexer *EhrIndexerTransactorSession) SetEhrSubject(subjectKey *big.Int, _ehrId *big.Int) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrSubject(&_EhrIndexer.TransactOpts, subjectKey, _ehrId)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0x651f1972.
//
// Solidity: function setEhrUser(uint256 userId, uint256 ehrId) returns(uint256)
func (_EhrIndexer *EhrIndexerTransactor) SetEhrUser(opts *bind.TransactOpts, userId *big.Int, ehrId *big.Int) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setEhrUser", userId, ehrId)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0x651f1972.
//
// Solidity: function setEhrUser(uint256 userId, uint256 ehrId) returns(uint256)
func (_EhrIndexer *EhrIndexerSession) SetEhrUser(userId *big.Int, ehrId *big.Int) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrUser(&_EhrIndexer.TransactOpts, userId, ehrId)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0x651f1972.
//
// Solidity: function setEhrUser(uint256 userId, uint256 ehrId) returns(uint256)
func (_EhrIndexer *EhrIndexerTransactorSession) SetEhrUser(userId *big.Int, ehrId *big.Int) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrUser(&_EhrIndexer.TransactOpts, userId, ehrId)
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

// EhrIndexerDataAccessChangedIterator is returned from FilterDataAccessChanged and is used to iterate over the raw logs and unpacked data for DataAccessChanged events raised by the EhrIndexer contract.
type EhrIndexerDataAccessChangedIterator struct {
	Event *EhrIndexerDataAccessChanged // Event containing the contract specifics and raw log

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
func (it *EhrIndexerDataAccessChangedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(EhrIndexerDataAccessChanged)
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
		it.Event = new(EhrIndexerDataAccessChanged)
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
func (it *EhrIndexerDataAccessChangedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *EhrIndexerDataAccessChangedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// EhrIndexerDataAccessChanged represents a DataAccessChanged event raised by the EhrIndexer contract.
type EhrIndexerDataAccessChanged struct {
	UserId *big.Int
	Access []byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDataAccessChanged is a free log retrieval operation binding the contract event 0xa9d658ce394dcc3ebc71e2f72f6440b8e02d9831bb72f3c59a622d6130e150bb.
//
// Solidity: event DataAccessChanged(uint256 userId, bytes access)
func (_EhrIndexer *EhrIndexerFilterer) FilterDataAccessChanged(opts *bind.FilterOpts) (*EhrIndexerDataAccessChangedIterator, error) {

	logs, sub, err := _EhrIndexer.contract.FilterLogs(opts, "DataAccessChanged")
	if err != nil {
		return nil, err
	}
	return &EhrIndexerDataAccessChangedIterator{contract: _EhrIndexer.contract, event: "DataAccessChanged", logs: logs, sub: sub}, nil
}

// WatchDataAccessChanged is a free log subscription operation binding the contract event 0xa9d658ce394dcc3ebc71e2f72f6440b8e02d9831bb72f3c59a622d6130e150bb.
//
// Solidity: event DataAccessChanged(uint256 userId, bytes access)
func (_EhrIndexer *EhrIndexerFilterer) WatchDataAccessChanged(opts *bind.WatchOpts, sink chan<- *EhrIndexerDataAccessChanged) (event.Subscription, error) {

	logs, sub, err := _EhrIndexer.contract.WatchLogs(opts, "DataAccessChanged")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(EhrIndexerDataAccessChanged)
				if err := _EhrIndexer.contract.UnpackLog(event, "DataAccessChanged", log); err != nil {
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

// ParseDataAccessChanged is a log parse operation binding the contract event 0xa9d658ce394dcc3ebc71e2f72f6440b8e02d9831bb72f3c59a622d6130e150bb.
//
// Solidity: event DataAccessChanged(uint256 userId, bytes access)
func (_EhrIndexer *EhrIndexerFilterer) ParseDataAccessChanged(log types.Log) (*EhrIndexerDataAccessChanged, error) {
	event := new(EhrIndexerDataAccessChanged)
	if err := _EhrIndexer.contract.UnpackLog(event, "DataAccessChanged", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
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
	UserId *big.Int
	Access []byte
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterDocAccessChanged is a free log retrieval operation binding the contract event 0xdfd7c68e099c5d969e8a9417e6319e2f985e5d2657e7f54b4af4a309aefcccbf.
//
// Solidity: event DocAccessChanged(uint256 userId, bytes access)
func (_EhrIndexer *EhrIndexerFilterer) FilterDocAccessChanged(opts *bind.FilterOpts) (*EhrIndexerDocAccessChangedIterator, error) {

	logs, sub, err := _EhrIndexer.contract.FilterLogs(opts, "DocAccessChanged")
	if err != nil {
		return nil, err
	}
	return &EhrIndexerDocAccessChangedIterator{contract: _EhrIndexer.contract, event: "DocAccessChanged", logs: logs, sub: sub}, nil
}

// WatchDocAccessChanged is a free log subscription operation binding the contract event 0xdfd7c68e099c5d969e8a9417e6319e2f985e5d2657e7f54b4af4a309aefcccbf.
//
// Solidity: event DocAccessChanged(uint256 userId, bytes access)
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

// ParseDocAccessChanged is a log parse operation binding the contract event 0xdfd7c68e099c5d969e8a9417e6319e2f985e5d2657e7f54b4af4a309aefcccbf.
//
// Solidity: event DocAccessChanged(uint256 userId, bytes access)
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
	EhrId     *big.Int
	StorageId *big.Int
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterEhrDocAdded is a free log retrieval operation binding the contract event 0xa6ed15209306a2d07d3ca4688cf4fb7d3fcd021baa6c2fb66b0af5e5598ed437.
//
// Solidity: event EhrDocAdded(uint256 ehrId, uint256 storageId)
func (_EhrIndexer *EhrIndexerFilterer) FilterEhrDocAdded(opts *bind.FilterOpts) (*EhrIndexerEhrDocAddedIterator, error) {

	logs, sub, err := _EhrIndexer.contract.FilterLogs(opts, "EhrDocAdded")
	if err != nil {
		return nil, err
	}
	return &EhrIndexerEhrDocAddedIterator{contract: _EhrIndexer.contract, event: "EhrDocAdded", logs: logs, sub: sub}, nil
}

// WatchEhrDocAdded is a free log subscription operation binding the contract event 0xa6ed15209306a2d07d3ca4688cf4fb7d3fcd021baa6c2fb66b0af5e5598ed437.
//
// Solidity: event EhrDocAdded(uint256 ehrId, uint256 storageId)
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

// ParseEhrDocAdded is a log parse operation binding the contract event 0xa6ed15209306a2d07d3ca4688cf4fb7d3fcd021baa6c2fb66b0af5e5598ed437.
//
// Solidity: event EhrDocAdded(uint256 ehrId, uint256 storageId)
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
	SubjectKey *big.Int
	EhrId      *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterEhrSubjectSet is a free log retrieval operation binding the contract event 0x7ab0e5c4dd7715a5cbe4c5dadf9bedd539a85f9bbf0bea823cbaf83d9dce54ea.
//
// Solidity: event EhrSubjectSet(uint256 subjectKey, uint256 ehrId)
func (_EhrIndexer *EhrIndexerFilterer) FilterEhrSubjectSet(opts *bind.FilterOpts) (*EhrIndexerEhrSubjectSetIterator, error) {

	logs, sub, err := _EhrIndexer.contract.FilterLogs(opts, "EhrSubjectSet")
	if err != nil {
		return nil, err
	}
	return &EhrIndexerEhrSubjectSetIterator{contract: _EhrIndexer.contract, event: "EhrSubjectSet", logs: logs, sub: sub}, nil
}

// WatchEhrSubjectSet is a free log subscription operation binding the contract event 0x7ab0e5c4dd7715a5cbe4c5dadf9bedd539a85f9bbf0bea823cbaf83d9dce54ea.
//
// Solidity: event EhrSubjectSet(uint256 subjectKey, uint256 ehrId)
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

// ParseEhrSubjectSet is a log parse operation binding the contract event 0x7ab0e5c4dd7715a5cbe4c5dadf9bedd539a85f9bbf0bea823cbaf83d9dce54ea.
//
// Solidity: event EhrSubjectSet(uint256 subjectKey, uint256 ehrId)
func (_EhrIndexer *EhrIndexerFilterer) ParseEhrSubjectSet(log types.Log) (*EhrIndexerEhrSubjectSet, error) {
	event := new(EhrIndexerEhrSubjectSet)
	if err := _EhrIndexer.contract.UnpackLog(event, "EhrSubjectSet", log); err != nil {
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
