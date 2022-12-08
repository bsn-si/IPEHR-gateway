// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package accessStore

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

// IAccessStoreAccess is an auto generated low-level Go binding around an user-defined struct.
type IAccessStoreAccess struct {
	IdHash  [32]byte
	IdEncr  []byte
	KeyEncr []byte
	Level   uint8
}

// AccessStoreMetaData contains all meta data concerning the AccessStore contract.
var AccessStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"accessID\",\"type\":\"bytes32\"}],\"name\":\"getAccess\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumIAccessStore.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structIAccessStore.Access[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"accessID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"accessIdHash\",\"type\":\"bytes32\"}],\"name\":\"getAccessByIdHash\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumIAccessStore.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structIAccessStore.Access\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"accessID\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumIAccessStore.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structIAccessStore.Access\",\"name\":\"o\",\"type\":\"tuple\"}],\"name\":\"setAccess\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userID\",\"type\":\"bytes32\"},{\"internalType\":\"enumIAccessStore.AccessKind\",\"name\":\"kind\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"}],\"name\":\"userAccess\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumIAccessStore.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structIAccessStore.Access\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// AccessStoreABI is the input ABI used to generate the binding from.
// Deprecated: Use AccessStoreMetaData.ABI instead.
var AccessStoreABI = AccessStoreMetaData.ABI

// AccessStore is an auto generated Go binding around an Ethereum contract.
type AccessStore struct {
	AccessStoreCaller     // Read-only binding to the contract
	AccessStoreTransactor // Write-only binding to the contract
	AccessStoreFilterer   // Log filterer for contract events
}

// AccessStoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type AccessStoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessStoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type AccessStoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessStoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type AccessStoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// AccessStoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type AccessStoreSession struct {
	Contract     *AccessStore      // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// AccessStoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type AccessStoreCallerSession struct {
	Contract *AccessStoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts      // Call options to use throughout this session
}

// AccessStoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type AccessStoreTransactorSession struct {
	Contract     *AccessStoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts      // Transaction auth options to use throughout this session
}

// AccessStoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type AccessStoreRaw struct {
	Contract *AccessStore // Generic contract binding to access the raw methods on
}

// AccessStoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type AccessStoreCallerRaw struct {
	Contract *AccessStoreCaller // Generic read-only contract binding to access the raw methods on
}

// AccessStoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type AccessStoreTransactorRaw struct {
	Contract *AccessStoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewAccessStore creates a new instance of AccessStore, bound to a specific deployed contract.
func NewAccessStore(address common.Address, backend bind.ContractBackend) (*AccessStore, error) {
	contract, err := bindAccessStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &AccessStore{AccessStoreCaller: AccessStoreCaller{contract: contract}, AccessStoreTransactor: AccessStoreTransactor{contract: contract}, AccessStoreFilterer: AccessStoreFilterer{contract: contract}}, nil
}

// NewAccessStoreCaller creates a new read-only instance of AccessStore, bound to a specific deployed contract.
func NewAccessStoreCaller(address common.Address, caller bind.ContractCaller) (*AccessStoreCaller, error) {
	contract, err := bindAccessStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &AccessStoreCaller{contract: contract}, nil
}

// NewAccessStoreTransactor creates a new write-only instance of AccessStore, bound to a specific deployed contract.
func NewAccessStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*AccessStoreTransactor, error) {
	contract, err := bindAccessStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &AccessStoreTransactor{contract: contract}, nil
}

// NewAccessStoreFilterer creates a new log filterer instance of AccessStore, bound to a specific deployed contract.
func NewAccessStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*AccessStoreFilterer, error) {
	contract, err := bindAccessStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &AccessStoreFilterer{contract: contract}, nil
}

// bindAccessStore binds a generic wrapper to an already deployed contract.
func bindAccessStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(AccessStoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccessStore *AccessStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessStore.Contract.AccessStoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccessStore *AccessStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessStore.Contract.AccessStoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccessStore *AccessStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessStore.Contract.AccessStoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_AccessStore *AccessStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _AccessStore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_AccessStore *AccessStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _AccessStore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_AccessStore *AccessStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _AccessStore.Contract.contract.Transact(opts, method, params...)
}

// GetAccess is a free data retrieval call binding the contract method 0x3347fcbe.
//
// Solidity: function getAccess(bytes32 accessID) view returns((bytes32,bytes,bytes,uint8)[])
func (_AccessStore *AccessStoreCaller) GetAccess(opts *bind.CallOpts, accessID [32]byte) ([]IAccessStoreAccess, error) {
	var out []interface{}
	err := _AccessStore.contract.Call(opts, &out, "getAccess", accessID)

	if err != nil {
		return *new([]IAccessStoreAccess), err
	}

	out0 := *abi.ConvertType(out[0], new([]IAccessStoreAccess)).(*[]IAccessStoreAccess)

	return out0, err

}

// GetAccess is a free data retrieval call binding the contract method 0x3347fcbe.
//
// Solidity: function getAccess(bytes32 accessID) view returns((bytes32,bytes,bytes,uint8)[])
func (_AccessStore *AccessStoreSession) GetAccess(accessID [32]byte) ([]IAccessStoreAccess, error) {
	return _AccessStore.Contract.GetAccess(&_AccessStore.CallOpts, accessID)
}

// GetAccess is a free data retrieval call binding the contract method 0x3347fcbe.
//
// Solidity: function getAccess(bytes32 accessID) view returns((bytes32,bytes,bytes,uint8)[])
func (_AccessStore *AccessStoreCallerSession) GetAccess(accessID [32]byte) ([]IAccessStoreAccess, error) {
	return _AccessStore.Contract.GetAccess(&_AccessStore.CallOpts, accessID)
}

// GetAccessByIdHash is a free data retrieval call binding the contract method 0x9ae2da76.
//
// Solidity: function getAccessByIdHash(bytes32 accessID, bytes32 accessIdHash) view returns((bytes32,bytes,bytes,uint8))
func (_AccessStore *AccessStoreCaller) GetAccessByIdHash(opts *bind.CallOpts, accessID [32]byte, accessIdHash [32]byte) (IAccessStoreAccess, error) {
	var out []interface{}
	err := _AccessStore.contract.Call(opts, &out, "getAccessByIdHash", accessID, accessIdHash)

	if err != nil {
		return *new(IAccessStoreAccess), err
	}

	out0 := *abi.ConvertType(out[0], new(IAccessStoreAccess)).(*IAccessStoreAccess)

	return out0, err

}

// GetAccessByIdHash is a free data retrieval call binding the contract method 0x9ae2da76.
//
// Solidity: function getAccessByIdHash(bytes32 accessID, bytes32 accessIdHash) view returns((bytes32,bytes,bytes,uint8))
func (_AccessStore *AccessStoreSession) GetAccessByIdHash(accessID [32]byte, accessIdHash [32]byte) (IAccessStoreAccess, error) {
	return _AccessStore.Contract.GetAccessByIdHash(&_AccessStore.CallOpts, accessID, accessIdHash)
}

// GetAccessByIdHash is a free data retrieval call binding the contract method 0x9ae2da76.
//
// Solidity: function getAccessByIdHash(bytes32 accessID, bytes32 accessIdHash) view returns((bytes32,bytes,bytes,uint8))
func (_AccessStore *AccessStoreCallerSession) GetAccessByIdHash(accessID [32]byte, accessIdHash [32]byte) (IAccessStoreAccess, error) {
	return _AccessStore.Contract.GetAccessByIdHash(&_AccessStore.CallOpts, accessID, accessIdHash)
}

// UserAccess is a free data retrieval call binding the contract method 0xa93f7898.
//
// Solidity: function userAccess(bytes32 userID, uint8 kind, bytes32 idHash) view returns((bytes32,bytes,bytes,uint8))
func (_AccessStore *AccessStoreCaller) UserAccess(opts *bind.CallOpts, userID [32]byte, kind uint8, idHash [32]byte) (IAccessStoreAccess, error) {
	var out []interface{}
	err := _AccessStore.contract.Call(opts, &out, "userAccess", userID, kind, idHash)

	if err != nil {
		return *new(IAccessStoreAccess), err
	}

	out0 := *abi.ConvertType(out[0], new(IAccessStoreAccess)).(*IAccessStoreAccess)

	return out0, err

}

// UserAccess is a free data retrieval call binding the contract method 0xa93f7898.
//
// Solidity: function userAccess(bytes32 userID, uint8 kind, bytes32 idHash) view returns((bytes32,bytes,bytes,uint8))
func (_AccessStore *AccessStoreSession) UserAccess(userID [32]byte, kind uint8, idHash [32]byte) (IAccessStoreAccess, error) {
	return _AccessStore.Contract.UserAccess(&_AccessStore.CallOpts, userID, kind, idHash)
}

// UserAccess is a free data retrieval call binding the contract method 0xa93f7898.
//
// Solidity: function userAccess(bytes32 userID, uint8 kind, bytes32 idHash) view returns((bytes32,bytes,bytes,uint8))
func (_AccessStore *AccessStoreCallerSession) UserAccess(userID [32]byte, kind uint8, idHash [32]byte) (IAccessStoreAccess, error) {
	return _AccessStore.Contract.UserAccess(&_AccessStore.CallOpts, userID, kind, idHash)
}

// SetAccess is a paid mutator transaction binding the contract method 0x2547e2fd.
//
// Solidity: function setAccess(bytes32 accessID, (bytes32,bytes,bytes,uint8) o) returns(uint8)
func (_AccessStore *AccessStoreTransactor) SetAccess(opts *bind.TransactOpts, accessID [32]byte, o IAccessStoreAccess) (*types.Transaction, error) {
	return _AccessStore.contract.Transact(opts, "setAccess", accessID, o)
}

// SetAccess is a paid mutator transaction binding the contract method 0x2547e2fd.
//
// Solidity: function setAccess(bytes32 accessID, (bytes32,bytes,bytes,uint8) o) returns(uint8)
func (_AccessStore *AccessStoreSession) SetAccess(accessID [32]byte, o IAccessStoreAccess) (*types.Transaction, error) {
	return _AccessStore.Contract.SetAccess(&_AccessStore.TransactOpts, accessID, o)
}

// SetAccess is a paid mutator transaction binding the contract method 0x2547e2fd.
//
// Solidity: function setAccess(bytes32 accessID, (bytes32,bytes,bytes,uint8) o) returns(uint8)
func (_AccessStore *AccessStoreTransactorSession) SetAccess(accessID [32]byte, o IAccessStoreAccess) (*types.Transaction, error) {
	return _AccessStore.Contract.SetAccess(&_AccessStore.TransactOpts, accessID, o)
}
