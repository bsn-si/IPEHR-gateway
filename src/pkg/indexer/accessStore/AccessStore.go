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
	Kind    uint8
	IdHash  [32]byte
	IdEncr  []byte
	KeyEncr []byte
	Level   uint8
}

// AccessStoreMetaData contains all meta data concerning the AccessStore contract.
var AccessStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedChange\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"accessID\",\"type\":\"bytes32\"}],\"name\":\"getAccess\",\"outputs\":[{\"components\":[{\"internalType\":\"enumIAccessStore.AccessKind\",\"name\":\"kind\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumIAccessStore.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structIAccessStore.Access[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"accessID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"accessIdHash\",\"type\":\"bytes32\"}],\"name\":\"getAccessByIdHash\",\"outputs\":[{\"components\":[{\"internalType\":\"enumIAccessStore.AccessKind\",\"name\":\"kind\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumIAccessStore.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structIAccessStore.Access\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"accessID\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"enumIAccessStore.AccessKind\",\"name\":\"kind\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumIAccessStore.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structIAccessStore.Access\",\"name\":\"a\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"setAccess\",\"outputs\":[{\"internalType\":\"enumIAccessStore.AccessAction\",\"name\":\"\",\"type\":\"uint8\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAllowed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_users\",\"type\":\"address\"}],\"name\":\"setUsersContractAddress\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumIAccessStore.AccessKind\",\"name\":\"kind\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"}],\"name\":\"userAccess\",\"outputs\":[{\"components\":[{\"internalType\":\"enumIAccessStore.AccessKind\",\"name\":\"kind\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumIAccessStore.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structIAccessStore.Access\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"users\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_AccessStore *AccessStoreCaller) AllowedChange(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _AccessStore.contract.Call(opts, &out, "allowedChange", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_AccessStore *AccessStoreSession) AllowedChange(arg0 common.Address) (bool, error) {
	return _AccessStore.Contract.AllowedChange(&_AccessStore.CallOpts, arg0)
}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_AccessStore *AccessStoreCallerSession) AllowedChange(arg0 common.Address) (bool, error) {
	return _AccessStore.Contract.AllowedChange(&_AccessStore.CallOpts, arg0)
}

// GetAccess is a free data retrieval call binding the contract method 0x3347fcbe.
//
// Solidity: function getAccess(bytes32 accessID) view returns((uint8,bytes32,bytes,bytes,uint8)[])
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
// Solidity: function getAccess(bytes32 accessID) view returns((uint8,bytes32,bytes,bytes,uint8)[])
func (_AccessStore *AccessStoreSession) GetAccess(accessID [32]byte) ([]IAccessStoreAccess, error) {
	return _AccessStore.Contract.GetAccess(&_AccessStore.CallOpts, accessID)
}

// GetAccess is a free data retrieval call binding the contract method 0x3347fcbe.
//
// Solidity: function getAccess(bytes32 accessID) view returns((uint8,bytes32,bytes,bytes,uint8)[])
func (_AccessStore *AccessStoreCallerSession) GetAccess(accessID [32]byte) ([]IAccessStoreAccess, error) {
	return _AccessStore.Contract.GetAccess(&_AccessStore.CallOpts, accessID)
}

// GetAccessByIdHash is a free data retrieval call binding the contract method 0x9ae2da76.
//
// Solidity: function getAccessByIdHash(bytes32 accessID, bytes32 accessIdHash) view returns((uint8,bytes32,bytes,bytes,uint8))
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
// Solidity: function getAccessByIdHash(bytes32 accessID, bytes32 accessIdHash) view returns((uint8,bytes32,bytes,bytes,uint8))
func (_AccessStore *AccessStoreSession) GetAccessByIdHash(accessID [32]byte, accessIdHash [32]byte) (IAccessStoreAccess, error) {
	return _AccessStore.Contract.GetAccessByIdHash(&_AccessStore.CallOpts, accessID, accessIdHash)
}

// GetAccessByIdHash is a free data retrieval call binding the contract method 0x9ae2da76.
//
// Solidity: function getAccessByIdHash(bytes32 accessID, bytes32 accessIdHash) view returns((uint8,bytes32,bytes,bytes,uint8))
func (_AccessStore *AccessStoreCallerSession) GetAccessByIdHash(accessID [32]byte, accessIdHash [32]byte) (IAccessStoreAccess, error) {
	return _AccessStore.Contract.GetAccessByIdHash(&_AccessStore.CallOpts, accessID, accessIdHash)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_AccessStore *AccessStoreCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _AccessStore.contract.Call(opts, &out, "nonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_AccessStore *AccessStoreSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _AccessStore.Contract.Nonces(&_AccessStore.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_AccessStore *AccessStoreCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _AccessStore.Contract.Nonces(&_AccessStore.CallOpts, arg0)
}

// UserAccess is a free data retrieval call binding the contract method 0xa93f7898.
//
// Solidity: function userAccess(bytes32 userIDHash, uint8 kind, bytes32 idHash) view returns((uint8,bytes32,bytes,bytes,uint8))
func (_AccessStore *AccessStoreCaller) UserAccess(opts *bind.CallOpts, userIDHash [32]byte, kind uint8, idHash [32]byte) (IAccessStoreAccess, error) {
	var out []interface{}
	err := _AccessStore.contract.Call(opts, &out, "userAccess", userIDHash, kind, idHash)

	if err != nil {
		return *new(IAccessStoreAccess), err
	}

	out0 := *abi.ConvertType(out[0], new(IAccessStoreAccess)).(*IAccessStoreAccess)

	return out0, err

}

// UserAccess is a free data retrieval call binding the contract method 0xa93f7898.
//
// Solidity: function userAccess(bytes32 userIDHash, uint8 kind, bytes32 idHash) view returns((uint8,bytes32,bytes,bytes,uint8))
func (_AccessStore *AccessStoreSession) UserAccess(userIDHash [32]byte, kind uint8, idHash [32]byte) (IAccessStoreAccess, error) {
	return _AccessStore.Contract.UserAccess(&_AccessStore.CallOpts, userIDHash, kind, idHash)
}

// UserAccess is a free data retrieval call binding the contract method 0xa93f7898.
//
// Solidity: function userAccess(bytes32 userIDHash, uint8 kind, bytes32 idHash) view returns((uint8,bytes32,bytes,bytes,uint8))
func (_AccessStore *AccessStoreCallerSession) UserAccess(userIDHash [32]byte, kind uint8, idHash [32]byte) (IAccessStoreAccess, error) {
	return _AccessStore.Contract.UserAccess(&_AccessStore.CallOpts, userIDHash, kind, idHash)
}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_AccessStore *AccessStoreCaller) Users(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _AccessStore.contract.Call(opts, &out, "users")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_AccessStore *AccessStoreSession) Users() (common.Address, error) {
	return _AccessStore.Contract.Users(&_AccessStore.CallOpts)
}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_AccessStore *AccessStoreCallerSession) Users() (common.Address, error) {
	return _AccessStore.Contract.Users(&_AccessStore.CallOpts)
}

// SetAccess is a paid mutator transaction binding the contract method 0xf44e860f.
//
// Solidity: function setAccess(bytes32 accessID, (uint8,bytes32,bytes,bytes,uint8) a, address signer, bytes signature) returns(uint8)
func (_AccessStore *AccessStoreTransactor) SetAccess(opts *bind.TransactOpts, accessID [32]byte, a IAccessStoreAccess, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _AccessStore.contract.Transact(opts, "setAccess", accessID, a, signer, signature)
}

// SetAccess is a paid mutator transaction binding the contract method 0xf44e860f.
//
// Solidity: function setAccess(bytes32 accessID, (uint8,bytes32,bytes,bytes,uint8) a, address signer, bytes signature) returns(uint8)
func (_AccessStore *AccessStoreSession) SetAccess(accessID [32]byte, a IAccessStoreAccess, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _AccessStore.Contract.SetAccess(&_AccessStore.TransactOpts, accessID, a, signer, signature)
}

// SetAccess is a paid mutator transaction binding the contract method 0xf44e860f.
//
// Solidity: function setAccess(bytes32 accessID, (uint8,bytes32,bytes,bytes,uint8) a, address signer, bytes signature) returns(uint8)
func (_AccessStore *AccessStoreTransactorSession) SetAccess(accessID [32]byte, a IAccessStoreAccess, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _AccessStore.Contract.SetAccess(&_AccessStore.TransactOpts, accessID, a, signer, signature)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_AccessStore *AccessStoreTransactor) SetAllowed(opts *bind.TransactOpts, addr common.Address, allowed bool) (*types.Transaction, error) {
	return _AccessStore.contract.Transact(opts, "setAllowed", addr, allowed)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_AccessStore *AccessStoreSession) SetAllowed(addr common.Address, allowed bool) (*types.Transaction, error) {
	return _AccessStore.Contract.SetAllowed(&_AccessStore.TransactOpts, addr, allowed)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_AccessStore *AccessStoreTransactorSession) SetAllowed(addr common.Address, allowed bool) (*types.Transaction, error) {
	return _AccessStore.Contract.SetAllowed(&_AccessStore.TransactOpts, addr, allowed)
}

// SetUsersContractAddress is a paid mutator transaction binding the contract method 0xdbfadfe2.
//
// Solidity: function setUsersContractAddress(address _users) returns()
func (_AccessStore *AccessStoreTransactor) SetUsersContractAddress(opts *bind.TransactOpts, _users common.Address) (*types.Transaction, error) {
	return _AccessStore.contract.Transact(opts, "setUsersContractAddress", _users)
}

// SetUsersContractAddress is a paid mutator transaction binding the contract method 0xdbfadfe2.
//
// Solidity: function setUsersContractAddress(address _users) returns()
func (_AccessStore *AccessStoreSession) SetUsersContractAddress(_users common.Address) (*types.Transaction, error) {
	return _AccessStore.Contract.SetUsersContractAddress(&_AccessStore.TransactOpts, _users)
}

// SetUsersContractAddress is a paid mutator transaction binding the contract method 0xdbfadfe2.
//
// Solidity: function setUsersContractAddress(address _users) returns()
func (_AccessStore *AccessStoreTransactorSession) SetUsersContractAddress(_users common.Address) (*types.Transaction, error) {
	return _AccessStore.Contract.SetUsersContractAddress(&_AccessStore.TransactOpts, _users)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AccessStore *AccessStoreTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _AccessStore.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AccessStore *AccessStoreSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AccessStore.Contract.TransferOwnership(&_AccessStore.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_AccessStore *AccessStoreTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _AccessStore.Contract.TransferOwnership(&_AccessStore.TransactOpts, newOwner)
}
