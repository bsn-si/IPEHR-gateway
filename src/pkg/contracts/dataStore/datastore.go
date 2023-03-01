// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package datastore

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

// DatastoreMetaData contains all meta data concerning the Datastore contract.
var DatastoreMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_users\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"groupID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"dataID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes32\",\"name\":\"ehrID\",\"type\":\"bytes32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"}],\"name\":\"DataUpdate\",\"type\":\"event\"},{\"inputs\":[],\"name\":\"accessStore\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedChange\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"groupID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"dataID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"ehrID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"dataUpdate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"ehrIndex\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAllowed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"users\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// DatastoreABI is the input ABI used to generate the binding from.
// Deprecated: Use DatastoreMetaData.ABI instead.
var DatastoreABI = DatastoreMetaData.ABI

// Datastore is an auto generated Go binding around an Ethereum contract.
type Datastore struct {
	DatastoreCaller     // Read-only binding to the contract
	DatastoreTransactor // Write-only binding to the contract
	DatastoreFilterer   // Log filterer for contract events
}

// DatastoreCaller is an auto generated read-only Go binding around an Ethereum contract.
type DatastoreCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DatastoreTransactor is an auto generated write-only Go binding around an Ethereum contract.
type DatastoreTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DatastoreFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type DatastoreFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// DatastoreSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type DatastoreSession struct {
	Contract     *Datastore        // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// DatastoreCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type DatastoreCallerSession struct {
	Contract *DatastoreCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts    // Call options to use throughout this session
}

// DatastoreTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type DatastoreTransactorSession struct {
	Contract     *DatastoreTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// DatastoreRaw is an auto generated low-level Go binding around an Ethereum contract.
type DatastoreRaw struct {
	Contract *Datastore // Generic contract binding to access the raw methods on
}

// DatastoreCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type DatastoreCallerRaw struct {
	Contract *DatastoreCaller // Generic read-only contract binding to access the raw methods on
}

// DatastoreTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type DatastoreTransactorRaw struct {
	Contract *DatastoreTransactor // Generic write-only contract binding to access the raw methods on
}

// NewDatastore creates a new instance of Datastore, bound to a specific deployed contract.
func NewDatastore(address common.Address, backend bind.ContractBackend) (*Datastore, error) {
	contract, err := bindDatastore(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Datastore{DatastoreCaller: DatastoreCaller{contract: contract}, DatastoreTransactor: DatastoreTransactor{contract: contract}, DatastoreFilterer: DatastoreFilterer{contract: contract}}, nil
}

// NewDatastoreCaller creates a new read-only instance of Datastore, bound to a specific deployed contract.
func NewDatastoreCaller(address common.Address, caller bind.ContractCaller) (*DatastoreCaller, error) {
	contract, err := bindDatastore(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &DatastoreCaller{contract: contract}, nil
}

// NewDatastoreTransactor creates a new write-only instance of Datastore, bound to a specific deployed contract.
func NewDatastoreTransactor(address common.Address, transactor bind.ContractTransactor) (*DatastoreTransactor, error) {
	contract, err := bindDatastore(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &DatastoreTransactor{contract: contract}, nil
}

// NewDatastoreFilterer creates a new log filterer instance of Datastore, bound to a specific deployed contract.
func NewDatastoreFilterer(address common.Address, filterer bind.ContractFilterer) (*DatastoreFilterer, error) {
	contract, err := bindDatastore(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &DatastoreFilterer{contract: contract}, nil
}

// bindDatastore binds a generic wrapper to an already deployed contract.
func bindDatastore(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(DatastoreABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Datastore *DatastoreRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Datastore.Contract.DatastoreCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Datastore *DatastoreRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Datastore.Contract.DatastoreTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Datastore *DatastoreRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Datastore.Contract.DatastoreTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Datastore *DatastoreCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Datastore.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Datastore *DatastoreTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Datastore.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Datastore *DatastoreTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Datastore.Contract.contract.Transact(opts, method, params...)
}

// AccessStore is a free data retrieval call binding the contract method 0x45e9ea6f.
//
// Solidity: function accessStore() view returns(address)
func (_Datastore *DatastoreCaller) AccessStore(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Datastore.contract.Call(opts, &out, "accessStore")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AccessStore is a free data retrieval call binding the contract method 0x45e9ea6f.
//
// Solidity: function accessStore() view returns(address)
func (_Datastore *DatastoreSession) AccessStore() (common.Address, error) {
	return _Datastore.Contract.AccessStore(&_Datastore.CallOpts)
}

// AccessStore is a free data retrieval call binding the contract method 0x45e9ea6f.
//
// Solidity: function accessStore() view returns(address)
func (_Datastore *DatastoreCallerSession) AccessStore() (common.Address, error) {
	return _Datastore.Contract.AccessStore(&_Datastore.CallOpts)
}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_Datastore *DatastoreCaller) AllowedChange(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Datastore.contract.Call(opts, &out, "allowedChange", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_Datastore *DatastoreSession) AllowedChange(arg0 common.Address) (bool, error) {
	return _Datastore.Contract.AllowedChange(&_Datastore.CallOpts, arg0)
}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_Datastore *DatastoreCallerSession) AllowedChange(arg0 common.Address) (bool, error) {
	return _Datastore.Contract.AllowedChange(&_Datastore.CallOpts, arg0)
}

// EhrIndex is a free data retrieval call binding the contract method 0x98655acb.
//
// Solidity: function ehrIndex() view returns(address)
func (_Datastore *DatastoreCaller) EhrIndex(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Datastore.contract.Call(opts, &out, "ehrIndex")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// EhrIndex is a free data retrieval call binding the contract method 0x98655acb.
//
// Solidity: function ehrIndex() view returns(address)
func (_Datastore *DatastoreSession) EhrIndex() (common.Address, error) {
	return _Datastore.Contract.EhrIndex(&_Datastore.CallOpts)
}

// EhrIndex is a free data retrieval call binding the contract method 0x98655acb.
//
// Solidity: function ehrIndex() view returns(address)
func (_Datastore *DatastoreCallerSession) EhrIndex() (common.Address, error) {
	return _Datastore.Contract.EhrIndex(&_Datastore.CallOpts)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_Datastore *DatastoreCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Datastore.contract.Call(opts, &out, "nonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_Datastore *DatastoreSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _Datastore.Contract.Nonces(&_Datastore.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_Datastore *DatastoreCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _Datastore.Contract.Nonces(&_Datastore.CallOpts, arg0)
}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_Datastore *DatastoreCaller) Users(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Datastore.contract.Call(opts, &out, "users")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_Datastore *DatastoreSession) Users() (common.Address, error) {
	return _Datastore.Contract.Users(&_Datastore.CallOpts)
}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_Datastore *DatastoreCallerSession) Users() (common.Address, error) {
	return _Datastore.Contract.Users(&_Datastore.CallOpts)
}

// DataUpdate is a paid mutator transaction binding the contract method 0xbd97b9c2.
//
// Solidity: function dataUpdate(bytes32 groupID, bytes32 dataID, bytes32 ehrID, bytes data, address signer, bytes signature) returns()
func (_Datastore *DatastoreTransactor) DataUpdate(opts *bind.TransactOpts, groupID [32]byte, dataID [32]byte, ehrID [32]byte, data []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Datastore.contract.Transact(opts, "dataUpdate", groupID, dataID, ehrID, data, signer, signature)
}

// DataUpdate is a paid mutator transaction binding the contract method 0xbd97b9c2.
//
// Solidity: function dataUpdate(bytes32 groupID, bytes32 dataID, bytes32 ehrID, bytes data, address signer, bytes signature) returns()
func (_Datastore *DatastoreSession) DataUpdate(groupID [32]byte, dataID [32]byte, ehrID [32]byte, data []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Datastore.Contract.DataUpdate(&_Datastore.TransactOpts, groupID, dataID, ehrID, data, signer, signature)
}

// DataUpdate is a paid mutator transaction binding the contract method 0xbd97b9c2.
//
// Solidity: function dataUpdate(bytes32 groupID, bytes32 dataID, bytes32 ehrID, bytes data, address signer, bytes signature) returns()
func (_Datastore *DatastoreTransactorSession) DataUpdate(groupID [32]byte, dataID [32]byte, ehrID [32]byte, data []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Datastore.Contract.DataUpdate(&_Datastore.TransactOpts, groupID, dataID, ehrID, data, signer, signature)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_Datastore *DatastoreTransactor) SetAllowed(opts *bind.TransactOpts, addr common.Address, allowed bool) (*types.Transaction, error) {
	return _Datastore.contract.Transact(opts, "setAllowed", addr, allowed)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_Datastore *DatastoreSession) SetAllowed(addr common.Address, allowed bool) (*types.Transaction, error) {
	return _Datastore.Contract.SetAllowed(&_Datastore.TransactOpts, addr, allowed)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_Datastore *DatastoreTransactorSession) SetAllowed(addr common.Address, allowed bool) (*types.Transaction, error) {
	return _Datastore.Contract.SetAllowed(&_Datastore.TransactOpts, addr, allowed)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Datastore *DatastoreTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Datastore.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Datastore *DatastoreSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Datastore.Contract.TransferOwnership(&_Datastore.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Datastore *DatastoreTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Datastore.Contract.TransferOwnership(&_Datastore.TransactOpts, newOwner)
}

// DatastoreDataUpdateIterator is returned from FilterDataUpdate and is used to iterate over the raw logs and unpacked data for DataUpdate events raised by the Datastore contract.
type DatastoreDataUpdateIterator struct {
	Event *DatastoreDataUpdate // Event containing the contract specifics and raw log

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
func (it *DatastoreDataUpdateIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(DatastoreDataUpdate)
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
		it.Event = new(DatastoreDataUpdate)
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
func (it *DatastoreDataUpdateIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *DatastoreDataUpdateIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// DatastoreDataUpdate represents a DataUpdate event raised by the Datastore contract.
type DatastoreDataUpdate struct {
	GroupID [32]byte
	DataID  [32]byte
	EhrID   [32]byte
	Data    []byte
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterDataUpdate is a free log retrieval operation binding the contract event 0x1412a9906d5462af327cae12f530e25ae94bdfd0ea6eb9aa4b8ad371ed56f8a2.
//
// Solidity: event DataUpdate(bytes32 groupID, bytes32 dataID, bytes32 ehrID, bytes data)
func (_Datastore *DatastoreFilterer) FilterDataUpdate(opts *bind.FilterOpts) (*DatastoreDataUpdateIterator, error) {

	logs, sub, err := _Datastore.contract.FilterLogs(opts, "DataUpdate")
	if err != nil {
		return nil, err
	}
	return &DatastoreDataUpdateIterator{contract: _Datastore.contract, event: "DataUpdate", logs: logs, sub: sub}, nil
}

// WatchDataUpdate is a free log subscription operation binding the contract event 0x1412a9906d5462af327cae12f530e25ae94bdfd0ea6eb9aa4b8ad371ed56f8a2.
//
// Solidity: event DataUpdate(bytes32 groupID, bytes32 dataID, bytes32 ehrID, bytes data)
func (_Datastore *DatastoreFilterer) WatchDataUpdate(opts *bind.WatchOpts, sink chan<- *DatastoreDataUpdate) (event.Subscription, error) {

	logs, sub, err := _Datastore.contract.WatchLogs(opts, "DataUpdate")
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(DatastoreDataUpdate)
				if err := _Datastore.contract.UnpackLog(event, "DataUpdate", log); err != nil {
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
func (_Datastore *DatastoreFilterer) ParseDataUpdate(log types.Log) (*DatastoreDataUpdate, error) {
	event := new(DatastoreDataUpdate)
	if err := _Datastore.contract.UnpackLog(event, "DataUpdate", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
