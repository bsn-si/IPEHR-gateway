// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package dataStore

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

// DataStoreMetaData contains all meta data concerning the DataStore contract.
var DataStoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_users\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"groupID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"dataID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"ehrID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"DataUpdate\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"accessStore\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedChange\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"groupID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"dataID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"ehrID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"dataUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ehrIndex\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAllowed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"users\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// DataStoreABI is the input ABI used to generate the binding from.
// Deprecated: Use DataStoreMetaData.ABI instead.
var DataStoreABI = DataStoreMetaData.ABI

// DataStore is an auto generated Go binding around an Ethereum contract.
type DataStore struct {
	DataStoreCaller     // Read-only binding to the contract
	DataStoreTransactor // Write-only binding to the contract
	DataStoreFilterer   // Log filterer for contract events
}

// DataStoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type DataStoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataStoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DataStoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataStoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DataStoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DataStoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DataStoreSession struct {
	Contract     *DataStore        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DataStoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DataStoreCallerSession struct {
	Contract *DataStoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// DataStoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DataStoreTransactorSession struct {
	Contract     *DataStoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// DataStoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type DataStoreRaw struct {
	Contract *DataStore // Generic contract binding to access the raw methods on
}

// DataStoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DataStoreCallerRaw struct {
	Contract *DataStoreCaller // Generic read-only contract binding to access the raw methods on
}

// DataStoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DataStoreTransactorRaw struct {
	Contract *DataStoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDataStore creates a new instance of DataStore, bound to a specific deployed contract.
func NewDataStore(address common.Address, backend bind.ContractBackend) (*DataStore, error) {
	contract, err := bindDataStore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &DataStore{DataStoreCaller: DataStoreCaller{contract: contract}, DataStoreTransactor: DataStoreTransactor{contract: contract}, DataStoreFilterer: DataStoreFilterer{contract: contract}}, nil
}

// NewDataStoreCaller creates a new read-only instance of DataStore, bound to a specific deployed contract.
func NewDataStoreCaller(address common.Address, caller bind.ContractCaller) (*DataStoreCaller, error) {
	contract, err := bindDataStore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DataStoreCaller{contract: contract}, nil
}

// NewDataStoreTransactor creates a new write-only instance of DataStore, bound to a specific deployed contract.
func NewDataStoreTransactor(address common.Address, transactor bind.ContractTransactor) (*DataStoreTransactor, error) {
	contract, err := bindDataStore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DataStoreTransactor{contract: contract}, nil
}

// NewDataStoreFilterer creates a new log filterer instance of DataStore, bound to a specific deployed contract.
func NewDataStoreFilterer(address common.Address, filterer bind.ContractFilterer) (*DataStoreFilterer, error) {
	contract, err := bindDataStore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DataStoreFilterer{contract: contract}, nil
}

// bindDataStore binds a generic wrapper to an already deployed contract.
func bindDataStore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DataStoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DataStore *DataStoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DataStore.Contract.DataStoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DataStore *DataStoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataStore.Contract.DataStoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DataStore *DataStoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DataStore.Contract.DataStoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_DataStore *DataStoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _DataStore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_DataStore *DataStoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _DataStore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_DataStore *DataStoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _DataStore.Contract.contract.Transact(opts, method, params...)
}

// AccessStore is a free data retrieval call binding the contract method 0x45e9ea6f.
//
// Solidity: function accessStore() view returns(address)
func (_DataStore *DataStoreCaller) AccessStore(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "accessStore")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AccessStore is a free data retrieval call binding the contract method 0x45e9ea6f.
//
// Solidity: function accessStore() view returns(address)
func (_DataStore *DataStoreSession) AccessStore() (common.Address, error) {
	return _DataStore.Contract.AccessStore(&_DataStore.CallOpts)
}

// AccessStore is a free data retrieval call binding the contract method 0x45e9ea6f.
//
// Solidity: function accessStore() view returns(address)
func (_DataStore *DataStoreCallerSession) AccessStore() (common.Address, error) {
	return _DataStore.Contract.AccessStore(&_DataStore.CallOpts)
}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_DataStore *DataStoreCaller) AllowedChange(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "allowedChange", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_DataStore *DataStoreSession) AllowedChange(arg0 common.Address) (bool, error) {
	return _DataStore.Contract.AllowedChange(&_DataStore.CallOpts, arg0)
}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_DataStore *DataStoreCallerSession) AllowedChange(arg0 common.Address) (bool, error) {
	return _DataStore.Contract.AllowedChange(&_DataStore.CallOpts, arg0)
}

// EhrIndex is a free data retrieval call binding the contract method 0x98655acb.
//
// Solidity: function ehrIndex() view returns(address)
func (_DataStore *DataStoreCaller) EhrIndex(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "ehrIndex")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EhrIndex is a free data retrieval call binding the contract method 0x98655acb.
//
// Solidity: function ehrIndex() view returns(address)
func (_DataStore *DataStoreSession) EhrIndex() (common.Address, error) {
	return _DataStore.Contract.EhrIndex(&_DataStore.CallOpts)
}

// EhrIndex is a free data retrieval call binding the contract method 0x98655acb.
//
// Solidity: function ehrIndex() view returns(address)
func (_DataStore *DataStoreCallerSession) EhrIndex() (common.Address, error) {
	return _DataStore.Contract.EhrIndex(&_DataStore.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_DataStore *DataStoreCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "nonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_DataStore *DataStoreSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _DataStore.Contract.Nonces(&_DataStore.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_DataStore *DataStoreCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _DataStore.Contract.Nonces(&_DataStore.CallOpts, arg0)
}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_DataStore *DataStoreCaller) Users(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _DataStore.contract.Call(opts, &out, "users")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_DataStore *DataStoreSession) Users() (common.Address, error) {
	return _DataStore.Contract.Users(&_DataStore.CallOpts)
}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_DataStore *DataStoreCallerSession) Users() (common.Address, error) {
	return _DataStore.Contract.Users(&_DataStore.CallOpts)
}

// DataUpdate is a paid mutator transaction binding the contract method 0xbd97b9c2.
//
// Solidity: function dataUpdate(bytes32 groupID, bytes32 dataID, bytes32 ehrID, bytes data, address signer, bytes signature) returns()
func (_DataStore *DataStoreTransactor) DataUpdate(opts *bind.TransactOpts, groupID [32]byte, dataID [32]byte, ehrID [32]byte, data []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "dataUpdate", groupID, dataID, ehrID, data, signer, signature)
}

// DataUpdate is a paid mutator transaction binding the contract method 0xbd97b9c2.
//
// Solidity: function dataUpdate(bytes32 groupID, bytes32 dataID, bytes32 ehrID, bytes data, address signer, bytes signature) returns()
func (_DataStore *DataStoreSession) DataUpdate(groupID [32]byte, dataID [32]byte, ehrID [32]byte, data []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _DataStore.Contract.DataUpdate(&_DataStore.TransactOpts, groupID, dataID, ehrID, data, signer, signature)
}

// DataUpdate is a paid mutator transaction binding the contract method 0xbd97b9c2.
//
// Solidity: function dataUpdate(bytes32 groupID, bytes32 dataID, bytes32 ehrID, bytes data, address signer, bytes signature) returns()
func (_DataStore *DataStoreTransactorSession) DataUpdate(groupID [32]byte, dataID [32]byte, ehrID [32]byte, data []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _DataStore.Contract.DataUpdate(&_DataStore.TransactOpts, groupID, dataID, ehrID, data, signer, signature)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_DataStore *DataStoreTransactor) SetAllowed(opts *bind.TransactOpts, addr common.Address, allowed bool) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "setAllowed", addr, allowed)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_DataStore *DataStoreSession) SetAllowed(addr common.Address, allowed bool) (*types.Transaction, error) {
	return _DataStore.Contract.SetAllowed(&_DataStore.TransactOpts, addr, allowed)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_DataStore *DataStoreTransactorSession) SetAllowed(addr common.Address, allowed bool) (*types.Transaction, error) {
	return _DataStore.Contract.SetAllowed(&_DataStore.TransactOpts, addr, allowed)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DataStore *DataStoreTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _DataStore.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DataStore *DataStoreSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DataStore.Contract.TransferOwnership(&_DataStore.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_DataStore *DataStoreTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _DataStore.Contract.TransferOwnership(&_DataStore.TransactOpts, newOwner)
}

// DataStoreDataUpdateIterator is returned from FilterDataUpdate and is used to iterate over the raw logs and unpacked data for DataUpdate events raised by the DataStore contract.
type DataStoreDataUpdateIterator struct {
	Event *DataStoreDataUpdate // Event containing the contract specifics and raw log

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
func (it *DataStoreDataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DataStoreDataUpdate)
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
		it.Event = new(DataStoreDataUpdate)
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
func (it *DataStoreDataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DataStoreDataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DataStoreDataUpdate represents a DataUpdate event raised by the DataStore contract.
type DataStoreDataUpdate struct {
	GroupID [32]byte
	DataID  [32]byte
	EhrID   [32]byte
	Data    []byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterDataUpdate is a free log retrieval operation binding the contract event 0x1412a9906d5462af327cae12f530e25ae94bdfd0ea6eb9aa4b8ad371ed56f8a2.
//
// Solidity: event DataUpdate(bytes32 groupID, bytes32 dataID, bytes32 ehrID, bytes data)
func (_DataStore *DataStoreFilterer) FilterDataUpdate(opts *bind.FilterOpts) (*DataStoreDataUpdateIterator, error) {

	logs, sub, err := _DataStore.contract.FilterLogs(opts, "DataUpdate")
	if err != nil {
		return nil, err
	}
	return &DataStoreDataUpdateIterator{contract: _DataStore.contract, event: "DataUpdate", logs: logs, sub: sub}, nil
}

// WatchDataUpdate is a free log subscription operation binding the contract event 0x1412a9906d5462af327cae12f530e25ae94bdfd0ea6eb9aa4b8ad371ed56f8a2.
//
// Solidity: event DataUpdate(bytes32 groupID, bytes32 dataID, bytes32 ehrID, bytes data)
func (_DataStore *DataStoreFilterer) WatchDataUpdate(opts *bind.WatchOpts, sink chan<- *DataStoreDataUpdate) (event.Subscription, error) {

	logs, sub, err := _DataStore.contract.WatchLogs(opts, "DataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DataStoreDataUpdate)
				if err := _DataStore.contract.UnpackLog(event, "DataUpdate", log); err != nil {
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

// ParseDataUpdate is a log parse operation binding the contract event 0x1412a9906d5462af327cae12f530e25ae94bdfd0ea6eb9aa4b8ad371ed56f8a2.
//
// Solidity: event DataUpdate(bytes32 groupID, bytes32 dataID, bytes32 ehrID, bytes data)
func (_DataStore *DataStoreFilterer) ParseDataUpdate(log types.Log) (*DataStoreDataUpdate, error) {
	event := new(DataStoreDataUpdate)
	if err := _DataStore.contract.UnpackLog(event, "DataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
