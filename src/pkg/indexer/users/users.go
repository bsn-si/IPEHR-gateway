// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package users

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

// AttributesAttribute is an auto generated low-level Go binding around an user-defined struct.
type AttributesAttribute struct {
	Code  uint8
	Value []byte
}

// IUsersGroupAddUserParams is an auto generated low-level Go binding around an user-defined struct.
type IUsersGroupAddUserParams struct {
	GroupIDHash [32]byte
	UserIDHash  [32]byte
	Level       uint8
	UserIDEncr  []byte
	KeyEncr     []byte
	Signer      common.Address
	Signature   []byte
}

// IUsersGroupMember is an auto generated low-level Go binding around an user-defined struct.
type IUsersGroupMember struct {
	UserIDHash [32]byte
	UserIDEncr []byte
}

// IUsersUser is an auto generated low-level Go binding around an user-defined struct.
type IUsersUser struct {
	IDHash [32]byte
	Role   uint8
	Attrs  []AttributesAttribute
}

// IUsersUserGroup is an auto generated low-level Go binding around an user-defined struct.
type IUsersUserGroup struct {
	Attrs   []AttributesAttribute
	Members []IUsersGroupMember
}

// UsersMetaData contains all meta data concerning the Users contract.
var UsersMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_accessStore\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"accessStore\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedChange\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ehrIndex\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"}],\"name\":\"getUser\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"IDHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumIUsers.Role\",\"name\":\"role\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structIUsers.User\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"code\",\"type\":\"uint64\"}],\"name\":\"getUserByCode\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"IDHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumIUsers.Role\",\"name\":\"role\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structIUsers.User\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"groupIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"userIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumIAccessStore.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"userIDEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structIUsers.GroupAddUserParams\",\"name\":\"p\",\"type\":\"tuple\"}],\"name\":\"groupAddUser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"groupIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"userIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"groupRemoveUser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAllowed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"groupIdHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"userGroupCreate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"groupIdHash\",\"type\":\"bytes32\"}],\"name\":\"userGroupGetByID\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"userIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"userIDEncr\",\"type\":\"bytes\"}],\"internalType\":\"structIUsers.GroupMember[]\",\"name\":\"members\",\"type\":\"tuple[]\"}],\"internalType\":\"structIUsers.UserGroup\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"IDHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumIUsers.Role\",\"name\":\"role\",\"type\":\"uint8\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"userNew\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"users\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// UsersABI is the input ABI used to generate the binding from.
// Deprecated: Use UsersMetaData.ABI instead.
var UsersABI = UsersMetaData.ABI

// Users is an auto generated Go binding around an Ethereum contract.
type Users struct {
	UsersCaller     // Read-only binding to the contract
	UsersTransactor // Write-only binding to the contract
	UsersFilterer   // Log filterer for contract events
}

// UsersCaller is an auto generated read-only Go binding around an Ethereum contract.
type UsersCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UsersTransactor is an auto generated write-only Go binding around an Ethereum contract.
type UsersTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UsersFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type UsersFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// UsersSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type UsersSession struct {
	Contract     *Users            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UsersCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type UsersCallerSession struct {
	Contract *UsersCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// UsersTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type UsersTransactorSession struct {
	Contract     *UsersTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// UsersRaw is an auto generated low-level Go binding around an Ethereum contract.
type UsersRaw struct {
	Contract *Users // Generic contract binding to access the raw methods on
}

// UsersCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type UsersCallerRaw struct {
	Contract *UsersCaller // Generic read-only contract binding to access the raw methods on
}

// UsersTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type UsersTransactorRaw struct {
	Contract *UsersTransactor // Generic write-only contract binding to access the raw methods on
}

// NewUsers creates a new instance of Users, bound to a specific deployed contract.
func NewUsers(address common.Address, backend bind.ContractBackend) (*Users, error) {
	contract, err := bindUsers(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Users{UsersCaller: UsersCaller{contract: contract}, UsersTransactor: UsersTransactor{contract: contract}, UsersFilterer: UsersFilterer{contract: contract}}, nil
}

// NewUsersCaller creates a new read-only instance of Users, bound to a specific deployed contract.
func NewUsersCaller(address common.Address, caller bind.ContractCaller) (*UsersCaller, error) {
	contract, err := bindUsers(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &UsersCaller{contract: contract}, nil
}

// NewUsersTransactor creates a new write-only instance of Users, bound to a specific deployed contract.
func NewUsersTransactor(address common.Address, transactor bind.ContractTransactor) (*UsersTransactor, error) {
	contract, err := bindUsers(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &UsersTransactor{contract: contract}, nil
}

// NewUsersFilterer creates a new log filterer instance of Users, bound to a specific deployed contract.
func NewUsersFilterer(address common.Address, filterer bind.ContractFilterer) (*UsersFilterer, error) {
	contract, err := bindUsers(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &UsersFilterer{contract: contract}, nil
}

// bindUsers binds a generic wrapper to an already deployed contract.
func bindUsers(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(UsersABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Users *UsersRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Users.Contract.UsersCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Users *UsersRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Users.Contract.UsersTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Users *UsersRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Users.Contract.UsersTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Users *UsersCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Users.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Users *UsersTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Users.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Users *UsersTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Users.Contract.contract.Transact(opts, method, params...)
}

// AccessStore is a free data retrieval call binding the contract method 0x45e9ea6f.
//
// Solidity: function accessStore() view returns(address)
func (_Users *UsersCaller) AccessStore(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Users.contract.Call(opts, &out, "accessStore")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AccessStore is a free data retrieval call binding the contract method 0x45e9ea6f.
//
// Solidity: function accessStore() view returns(address)
func (_Users *UsersSession) AccessStore() (common.Address, error) {
	return _Users.Contract.AccessStore(&_Users.CallOpts)
}

// AccessStore is a free data retrieval call binding the contract method 0x45e9ea6f.
//
// Solidity: function accessStore() view returns(address)
func (_Users *UsersCallerSession) AccessStore() (common.Address, error) {
	return _Users.Contract.AccessStore(&_Users.CallOpts)
}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_Users *UsersCaller) AllowedChange(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Users.contract.Call(opts, &out, "allowedChange", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_Users *UsersSession) AllowedChange(arg0 common.Address) (bool, error) {
	return _Users.Contract.AllowedChange(&_Users.CallOpts, arg0)
}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_Users *UsersCallerSession) AllowedChange(arg0 common.Address) (bool, error) {
	return _Users.Contract.AllowedChange(&_Users.CallOpts, arg0)
}

// EhrIndex is a free data retrieval call binding the contract method 0x98655acb.
//
// Solidity: function ehrIndex() view returns(address)
func (_Users *UsersCaller) EhrIndex(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Users.contract.Call(opts, &out, "ehrIndex")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EhrIndex is a free data retrieval call binding the contract method 0x98655acb.
//
// Solidity: function ehrIndex() view returns(address)
func (_Users *UsersSession) EhrIndex() (common.Address, error) {
	return _Users.Contract.EhrIndex(&_Users.CallOpts)
}

// EhrIndex is a free data retrieval call binding the contract method 0x98655acb.
//
// Solidity: function ehrIndex() view returns(address)
func (_Users *UsersCallerSession) EhrIndex() (common.Address, error) {
	return _Users.Contract.EhrIndex(&_Users.CallOpts)
}

// GetUser is a free data retrieval call binding the contract method 0x6f77926b.
//
// Solidity: function getUser(address addr) view returns((bytes32,uint8,(uint8,bytes)[]))
func (_Users *UsersCaller) GetUser(opts *bind.CallOpts, addr common.Address) (IUsersUser, error) {
	var out []interface{}
	err := _Users.contract.Call(opts, &out, "getUser", addr)

	if err != nil {
		return *new(IUsersUser), err
	}

	out0 := *abi.ConvertType(out[0], new(IUsersUser)).(*IUsersUser)

	return out0, err

}

// GetUser is a free data retrieval call binding the contract method 0x6f77926b.
//
// Solidity: function getUser(address addr) view returns((bytes32,uint8,(uint8,bytes)[]))
func (_Users *UsersSession) GetUser(addr common.Address) (IUsersUser, error) {
	return _Users.Contract.GetUser(&_Users.CallOpts, addr)
}

// GetUser is a free data retrieval call binding the contract method 0x6f77926b.
//
// Solidity: function getUser(address addr) view returns((bytes32,uint8,(uint8,bytes)[]))
func (_Users *UsersCallerSession) GetUser(addr common.Address) (IUsersUser, error) {
	return _Users.Contract.GetUser(&_Users.CallOpts, addr)
}

// GetUserByCode is a free data retrieval call binding the contract method 0x138d131f.
//
// Solidity: function getUserByCode(uint64 code) view returns((bytes32,uint8,(uint8,bytes)[]))
func (_Users *UsersCaller) GetUserByCode(opts *bind.CallOpts, code uint64) (IUsersUser, error) {
	var out []interface{}
	err := _Users.contract.Call(opts, &out, "getUserByCode", code)

	if err != nil {
		return *new(IUsersUser), err
	}

	out0 := *abi.ConvertType(out[0], new(IUsersUser)).(*IUsersUser)

	return out0, err

}

// GetUserByCode is a free data retrieval call binding the contract method 0x138d131f.
//
// Solidity: function getUserByCode(uint64 code) view returns((bytes32,uint8,(uint8,bytes)[]))
func (_Users *UsersSession) GetUserByCode(code uint64) (IUsersUser, error) {
	return _Users.Contract.GetUserByCode(&_Users.CallOpts, code)
}

// GetUserByCode is a free data retrieval call binding the contract method 0x138d131f.
//
// Solidity: function getUserByCode(uint64 code) view returns((bytes32,uint8,(uint8,bytes)[]))
func (_Users *UsersCallerSession) GetUserByCode(code uint64) (IUsersUser, error) {
	return _Users.Contract.GetUserByCode(&_Users.CallOpts, code)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_Users *UsersCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Users.contract.Call(opts, &out, "nonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_Users *UsersSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _Users.Contract.Nonces(&_Users.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_Users *UsersCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _Users.Contract.Nonces(&_Users.CallOpts, arg0)
}

// UserGroupGetByID is a free data retrieval call binding the contract method 0xc1dfff99.
//
// Solidity: function userGroupGetByID(bytes32 groupIdHash) view returns(((uint8,bytes)[],(bytes32,bytes)[]))
func (_Users *UsersCaller) UserGroupGetByID(opts *bind.CallOpts, groupIdHash [32]byte) (IUsersUserGroup, error) {
	var out []interface{}
	err := _Users.contract.Call(opts, &out, "userGroupGetByID", groupIdHash)

	if err != nil {
		return *new(IUsersUserGroup), err
	}

	out0 := *abi.ConvertType(out[0], new(IUsersUserGroup)).(*IUsersUserGroup)

	return out0, err

}

// UserGroupGetByID is a free data retrieval call binding the contract method 0xc1dfff99.
//
// Solidity: function userGroupGetByID(bytes32 groupIdHash) view returns(((uint8,bytes)[],(bytes32,bytes)[]))
func (_Users *UsersSession) UserGroupGetByID(groupIdHash [32]byte) (IUsersUserGroup, error) {
	return _Users.Contract.UserGroupGetByID(&_Users.CallOpts, groupIdHash)
}

// UserGroupGetByID is a free data retrieval call binding the contract method 0xc1dfff99.
//
// Solidity: function userGroupGetByID(bytes32 groupIdHash) view returns(((uint8,bytes)[],(bytes32,bytes)[]))
func (_Users *UsersCallerSession) UserGroupGetByID(groupIdHash [32]byte) (IUsersUserGroup, error) {
	return _Users.Contract.UserGroupGetByID(&_Users.CallOpts, groupIdHash)
}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_Users *UsersCaller) Users(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Users.contract.Call(opts, &out, "users")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_Users *UsersSession) Users() (common.Address, error) {
	return _Users.Contract.Users(&_Users.CallOpts)
}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_Users *UsersCallerSession) Users() (common.Address, error) {
	return _Users.Contract.Users(&_Users.CallOpts)
}

// GroupAddUser is a paid mutator transaction binding the contract method 0xa544050a.
//
// Solidity: function groupAddUser((bytes32,bytes32,uint8,bytes,bytes,address,bytes) p) returns()
func (_Users *UsersTransactor) GroupAddUser(opts *bind.TransactOpts, p IUsersGroupAddUserParams) (*types.Transaction, error) {
	return _Users.contract.Transact(opts, "groupAddUser", p)
}

// GroupAddUser is a paid mutator transaction binding the contract method 0xa544050a.
//
// Solidity: function groupAddUser((bytes32,bytes32,uint8,bytes,bytes,address,bytes) p) returns()
func (_Users *UsersSession) GroupAddUser(p IUsersGroupAddUserParams) (*types.Transaction, error) {
	return _Users.Contract.GroupAddUser(&_Users.TransactOpts, p)
}

// GroupAddUser is a paid mutator transaction binding the contract method 0xa544050a.
//
// Solidity: function groupAddUser((bytes32,bytes32,uint8,bytes,bytes,address,bytes) p) returns()
func (_Users *UsersTransactorSession) GroupAddUser(p IUsersGroupAddUserParams) (*types.Transaction, error) {
	return _Users.Contract.GroupAddUser(&_Users.TransactOpts, p)
}

// GroupRemoveUser is a paid mutator transaction binding the contract method 0xb94e5cee.
//
// Solidity: function groupRemoveUser(bytes32 groupIDHash, bytes32 userIDHash, address signer, bytes signature) returns()
func (_Users *UsersTransactor) GroupRemoveUser(opts *bind.TransactOpts, groupIDHash [32]byte, userIDHash [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Users.contract.Transact(opts, "groupRemoveUser", groupIDHash, userIDHash, signer, signature)
}

// GroupRemoveUser is a paid mutator transaction binding the contract method 0xb94e5cee.
//
// Solidity: function groupRemoveUser(bytes32 groupIDHash, bytes32 userIDHash, address signer, bytes signature) returns()
func (_Users *UsersSession) GroupRemoveUser(groupIDHash [32]byte, userIDHash [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Users.Contract.GroupRemoveUser(&_Users.TransactOpts, groupIDHash, userIDHash, signer, signature)
}

// GroupRemoveUser is a paid mutator transaction binding the contract method 0xb94e5cee.
//
// Solidity: function groupRemoveUser(bytes32 groupIDHash, bytes32 userIDHash, address signer, bytes signature) returns()
func (_Users *UsersTransactorSession) GroupRemoveUser(groupIDHash [32]byte, userIDHash [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Users.Contract.GroupRemoveUser(&_Users.TransactOpts, groupIDHash, userIDHash, signer, signature)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Users *UsersTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _Users.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Users *UsersSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _Users.Contract.Multicall(&_Users.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Users *UsersTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _Users.Contract.Multicall(&_Users.TransactOpts, data)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_Users *UsersTransactor) SetAllowed(opts *bind.TransactOpts, addr common.Address, allowed bool) (*types.Transaction, error) {
	return _Users.contract.Transact(opts, "setAllowed", addr, allowed)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_Users *UsersSession) SetAllowed(addr common.Address, allowed bool) (*types.Transaction, error) {
	return _Users.Contract.SetAllowed(&_Users.TransactOpts, addr, allowed)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_Users *UsersTransactorSession) SetAllowed(addr common.Address, allowed bool) (*types.Transaction, error) {
	return _Users.Contract.SetAllowed(&_Users.TransactOpts, addr, allowed)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Users *UsersTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Users.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Users *UsersSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Users.Contract.TransferOwnership(&_Users.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Users *UsersTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Users.Contract.TransferOwnership(&_Users.TransactOpts, newOwner)
}

// UserGroupCreate is a paid mutator transaction binding the contract method 0xc430bb90.
//
// Solidity: function userGroupCreate(bytes32 groupIdHash, (uint8,bytes)[] attrs, address signer, bytes signature) returns()
func (_Users *UsersTransactor) UserGroupCreate(opts *bind.TransactOpts, groupIdHash [32]byte, attrs []AttributesAttribute, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Users.contract.Transact(opts, "userGroupCreate", groupIdHash, attrs, signer, signature)
}

// UserGroupCreate is a paid mutator transaction binding the contract method 0xc430bb90.
//
// Solidity: function userGroupCreate(bytes32 groupIdHash, (uint8,bytes)[] attrs, address signer, bytes signature) returns()
func (_Users *UsersSession) UserGroupCreate(groupIdHash [32]byte, attrs []AttributesAttribute, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Users.Contract.UserGroupCreate(&_Users.TransactOpts, groupIdHash, attrs, signer, signature)
}

// UserGroupCreate is a paid mutator transaction binding the contract method 0xc430bb90.
//
// Solidity: function userGroupCreate(bytes32 groupIdHash, (uint8,bytes)[] attrs, address signer, bytes signature) returns()
func (_Users *UsersTransactorSession) UserGroupCreate(groupIdHash [32]byte, attrs []AttributesAttribute, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Users.Contract.UserGroupCreate(&_Users.TransactOpts, groupIdHash, attrs, signer, signature)
}

// UserNew is a paid mutator transaction binding the contract method 0x1bd5a4f9.
//
// Solidity: function userNew(address addr, bytes32 IDHash, uint8 role, (uint8,bytes)[] attrs, address signer, bytes signature) returns()
func (_Users *UsersTransactor) UserNew(opts *bind.TransactOpts, addr common.Address, IDHash [32]byte, role uint8, attrs []AttributesAttribute, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Users.contract.Transact(opts, "userNew", addr, IDHash, role, attrs, signer, signature)
}

// UserNew is a paid mutator transaction binding the contract method 0x1bd5a4f9.
//
// Solidity: function userNew(address addr, bytes32 IDHash, uint8 role, (uint8,bytes)[] attrs, address signer, bytes signature) returns()
func (_Users *UsersSession) UserNew(addr common.Address, IDHash [32]byte, role uint8, attrs []AttributesAttribute, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Users.Contract.UserNew(&_Users.TransactOpts, addr, IDHash, role, attrs, signer, signature)
}

// UserNew is a paid mutator transaction binding the contract method 0x1bd5a4f9.
//
// Solidity: function userNew(address addr, bytes32 IDHash, uint8 role, (uint8,bytes)[] attrs, address signer, bytes signature) returns()
func (_Users *UsersTransactorSession) UserNew(addr common.Address, IDHash [32]byte, role uint8, attrs []AttributesAttribute, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Users.Contract.UserNew(&_Users.TransactOpts, addr, IDHash, role, attrs, signer, signature)
}
