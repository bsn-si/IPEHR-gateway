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

// AttributesAttribute is an auto generated low-level Go binding around an user-defined struct.
type AttributesAttribute struct {
	Code  uint8
	Value []byte
}

// DocGroupsDocGroupCreateParams is an auto generated low-level Go binding around an user-defined struct.
type DocGroupsDocGroupCreateParams struct {
	GroupIDHash [32]byte
	Attrs       []AttributesAttribute
	Signer      common.Address
	Signature   []byte
}

// DocsAddEhrDocParams is an auto generated low-level Go binding around an user-defined struct.
type DocsAddEhrDocParams struct {
	DocType   uint8
	Id        []byte
	Version   []byte
	Timestamp uint32
	Attrs     []AttributesAttribute
	Signer    common.Address
	Signature []byte
}

// DocsDocumentMeta is an auto generated low-level Go binding around an user-defined struct.
type DocsDocumentMeta struct {
	Status    uint8
	Id        []byte
	Version   []byte
	Timestamp uint32
	IsLast    bool
	Attrs     []AttributesAttribute
}

// IAccessStoreAccess is an auto generated low-level Go binding around an user-defined struct.
type IAccessStoreAccess struct {
	IdHash  [32]byte
	IdEncr  []byte
	KeyEncr []byte
	Level   uint8
}

// EhrIndexerMetaData contains all meta data concerning the EhrIndexer contract.
var EhrIndexerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_accessStore\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_users\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"accessStore\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structDocs.AddEhrDocParams\",\"name\":\"p\",\"type\":\"tuple\"}],\"name\":\"addEhrDoc\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedChange\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"deleteDoc\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"groupIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"docCIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"docCIDEncr\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"docGroupAddDoc\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"groupIDHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structDocGroups.DocGroupCreateParams\",\"name\":\"p\",\"type\":\"tuple\"}],\"name\":\"docGroupCreate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"groupIdHash\",\"type\":\"bytes32\"}],\"name\":\"docGroupGetAttrs\",\"outputs\":[{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"groupIdHash\",\"type\":\"bytes32\"}],\"name\":\"docGroupGetDocs\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"ehrSubject\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrID\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"name\":\"getDocByTime\",\"outputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structDocs.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"}],\"name\":\"getDocByVersion\",\"outputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structDocs.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"UIDHash\",\"type\":\"bytes32\"}],\"name\":\"getDocLastByBaseID\",\"outputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structDocs.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"}],\"name\":\"getEhrDocs\",\"outputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structDocs.DocumentMeta[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userIDHash\",\"type\":\"bytes32\"}],\"name\":\"getEhrUser\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"}],\"name\":\"getLastEhrDocByType\",\"outputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structDocs.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAllowed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"CIDHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumIAccessStore.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structIAccessStore.Access\",\"name\":\"access\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"userAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"setDocAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"subjectKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"setEhrSubject\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"IDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"setEhrUser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"users\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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

// AccessStore is a free data retrieval call binding the contract method 0x45e9ea6f.
//
// Solidity: function accessStore() view returns(address)
func (_EhrIndexer *EhrIndexerCaller) AccessStore(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "accessStore")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AccessStore is a free data retrieval call binding the contract method 0x45e9ea6f.
//
// Solidity: function accessStore() view returns(address)
func (_EhrIndexer *EhrIndexerSession) AccessStore() (common.Address, error) {
	return _EhrIndexer.Contract.AccessStore(&_EhrIndexer.CallOpts)
}

// AccessStore is a free data retrieval call binding the contract method 0x45e9ea6f.
//
// Solidity: function accessStore() view returns(address)
func (_EhrIndexer *EhrIndexerCallerSession) AccessStore() (common.Address, error) {
	return _EhrIndexer.Contract.AccessStore(&_EhrIndexer.CallOpts)
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

// DocGroupGetAttrs is a free data retrieval call binding the contract method 0x891c952e.
//
// Solidity: function docGroupGetAttrs(bytes32 groupIdHash) view returns((uint8,bytes)[])
func (_EhrIndexer *EhrIndexerCaller) DocGroupGetAttrs(opts *bind.CallOpts, groupIdHash [32]byte) ([]AttributesAttribute, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "docGroupGetAttrs", groupIdHash)

	if err != nil {
		return *new([]AttributesAttribute), err
	}

	out0 := *abi.ConvertType(out[0], new([]AttributesAttribute)).(*[]AttributesAttribute)

	return out0, err

}

// DocGroupGetAttrs is a free data retrieval call binding the contract method 0x891c952e.
//
// Solidity: function docGroupGetAttrs(bytes32 groupIdHash) view returns((uint8,bytes)[])
func (_EhrIndexer *EhrIndexerSession) DocGroupGetAttrs(groupIdHash [32]byte) ([]AttributesAttribute, error) {
	return _EhrIndexer.Contract.DocGroupGetAttrs(&_EhrIndexer.CallOpts, groupIdHash)
}

// DocGroupGetAttrs is a free data retrieval call binding the contract method 0x891c952e.
//
// Solidity: function docGroupGetAttrs(bytes32 groupIdHash) view returns((uint8,bytes)[])
func (_EhrIndexer *EhrIndexerCallerSession) DocGroupGetAttrs(groupIdHash [32]byte) ([]AttributesAttribute, error) {
	return _EhrIndexer.Contract.DocGroupGetAttrs(&_EhrIndexer.CallOpts, groupIdHash)
}

// DocGroupGetDocs is a free data retrieval call binding the contract method 0x88216834.
//
// Solidity: function docGroupGetDocs(bytes32 groupIdHash) view returns(bytes[])
func (_EhrIndexer *EhrIndexerCaller) DocGroupGetDocs(opts *bind.CallOpts, groupIdHash [32]byte) ([][]byte, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "docGroupGetDocs", groupIdHash)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

// DocGroupGetDocs is a free data retrieval call binding the contract method 0x88216834.
//
// Solidity: function docGroupGetDocs(bytes32 groupIdHash) view returns(bytes[])
func (_EhrIndexer *EhrIndexerSession) DocGroupGetDocs(groupIdHash [32]byte) ([][]byte, error) {
	return _EhrIndexer.Contract.DocGroupGetDocs(&_EhrIndexer.CallOpts, groupIdHash)
}

// DocGroupGetDocs is a free data retrieval call binding the contract method 0x88216834.
//
// Solidity: function docGroupGetDocs(bytes32 groupIdHash) view returns(bytes[])
func (_EhrIndexer *EhrIndexerCallerSession) DocGroupGetDocs(groupIdHash [32]byte) ([][]byte, error) {
	return _EhrIndexer.Contract.DocGroupGetDocs(&_EhrIndexer.CallOpts, groupIdHash)
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

// GetDocByTime is a free data retrieval call binding the contract method 0x4f722b16.
//
// Solidity: function getDocByTime(bytes32 ehrID, uint8 docType, uint32 timestamp) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_EhrIndexer *EhrIndexerCaller) GetDocByTime(opts *bind.CallOpts, ehrID [32]byte, docType uint8, timestamp uint32) (DocsDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getDocByTime", ehrID, docType, timestamp)

	if err != nil {
		return *new(DocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(DocsDocumentMeta)).(*DocsDocumentMeta)

	return out0, err

}

// GetDocByTime is a free data retrieval call binding the contract method 0x4f722b16.
//
// Solidity: function getDocByTime(bytes32 ehrID, uint8 docType, uint32 timestamp) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_EhrIndexer *EhrIndexerSession) GetDocByTime(ehrID [32]byte, docType uint8, timestamp uint32) (DocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocByTime(&_EhrIndexer.CallOpts, ehrID, docType, timestamp)
}

// GetDocByTime is a free data retrieval call binding the contract method 0x4f722b16.
//
// Solidity: function getDocByTime(bytes32 ehrID, uint8 docType, uint32 timestamp) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_EhrIndexer *EhrIndexerCallerSession) GetDocByTime(ehrID [32]byte, docType uint8, timestamp uint32) (DocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocByTime(&_EhrIndexer.CallOpts, ehrID, docType, timestamp)
}

// GetDocByVersion is a free data retrieval call binding the contract method 0x179fcacf.
//
// Solidity: function getDocByVersion(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_EhrIndexer *EhrIndexerCaller) GetDocByVersion(opts *bind.CallOpts, ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte) (DocsDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getDocByVersion", ehrId, docType, docBaseUIDHash, version)

	if err != nil {
		return *new(DocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(DocsDocumentMeta)).(*DocsDocumentMeta)

	return out0, err

}

// GetDocByVersion is a free data retrieval call binding the contract method 0x179fcacf.
//
// Solidity: function getDocByVersion(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_EhrIndexer *EhrIndexerSession) GetDocByVersion(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte) (DocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocByVersion(&_EhrIndexer.CallOpts, ehrId, docType, docBaseUIDHash, version)
}

// GetDocByVersion is a free data retrieval call binding the contract method 0x179fcacf.
//
// Solidity: function getDocByVersion(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_EhrIndexer *EhrIndexerCallerSession) GetDocByVersion(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte) (DocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocByVersion(&_EhrIndexer.CallOpts, ehrId, docType, docBaseUIDHash, version)
}

// GetDocLastByBaseID is a free data retrieval call binding the contract method 0x949c8e77.
//
// Solidity: function getDocLastByBaseID(bytes32 userIDHash, uint8 docType, bytes32 UIDHash) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_EhrIndexer *EhrIndexerCaller) GetDocLastByBaseID(opts *bind.CallOpts, userIDHash [32]byte, docType uint8, UIDHash [32]byte) (DocsDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getDocLastByBaseID", userIDHash, docType, UIDHash)

	if err != nil {
		return *new(DocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(DocsDocumentMeta)).(*DocsDocumentMeta)

	return out0, err

}

// GetDocLastByBaseID is a free data retrieval call binding the contract method 0x949c8e77.
//
// Solidity: function getDocLastByBaseID(bytes32 userIDHash, uint8 docType, bytes32 UIDHash) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_EhrIndexer *EhrIndexerSession) GetDocLastByBaseID(userIDHash [32]byte, docType uint8, UIDHash [32]byte) (DocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocLastByBaseID(&_EhrIndexer.CallOpts, userIDHash, docType, UIDHash)
}

// GetDocLastByBaseID is a free data retrieval call binding the contract method 0x949c8e77.
//
// Solidity: function getDocLastByBaseID(bytes32 userIDHash, uint8 docType, bytes32 UIDHash) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_EhrIndexer *EhrIndexerCallerSession) GetDocLastByBaseID(userIDHash [32]byte, docType uint8, UIDHash [32]byte) (DocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocLastByBaseID(&_EhrIndexer.CallOpts, userIDHash, docType, UIDHash)
}

// GetEhrDocs is a free data retrieval call binding the contract method 0xeec6b331.
//
// Solidity: function getEhrDocs(bytes32 userIDHash, uint8 docType) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[])[])
func (_EhrIndexer *EhrIndexerCaller) GetEhrDocs(opts *bind.CallOpts, userIDHash [32]byte, docType uint8) ([]DocsDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getEhrDocs", userIDHash, docType)

	if err != nil {
		return *new([]DocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new([]DocsDocumentMeta)).(*[]DocsDocumentMeta)

	return out0, err

}

// GetEhrDocs is a free data retrieval call binding the contract method 0xeec6b331.
//
// Solidity: function getEhrDocs(bytes32 userIDHash, uint8 docType) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[])[])
func (_EhrIndexer *EhrIndexerSession) GetEhrDocs(userIDHash [32]byte, docType uint8) ([]DocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetEhrDocs(&_EhrIndexer.CallOpts, userIDHash, docType)
}

// GetEhrDocs is a free data retrieval call binding the contract method 0xeec6b331.
//
// Solidity: function getEhrDocs(bytes32 userIDHash, uint8 docType) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[])[])
func (_EhrIndexer *EhrIndexerCallerSession) GetEhrDocs(userIDHash [32]byte, docType uint8) ([]DocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetEhrDocs(&_EhrIndexer.CallOpts, userIDHash, docType)
}

// GetEhrUser is a free data retrieval call binding the contract method 0xf228d6f0.
//
// Solidity: function getEhrUser(bytes32 userIDHash) view returns(bytes32)
func (_EhrIndexer *EhrIndexerCaller) GetEhrUser(opts *bind.CallOpts, userIDHash [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getEhrUser", userIDHash)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetEhrUser is a free data retrieval call binding the contract method 0xf228d6f0.
//
// Solidity: function getEhrUser(bytes32 userIDHash) view returns(bytes32)
func (_EhrIndexer *EhrIndexerSession) GetEhrUser(userIDHash [32]byte) ([32]byte, error) {
	return _EhrIndexer.Contract.GetEhrUser(&_EhrIndexer.CallOpts, userIDHash)
}

// GetEhrUser is a free data retrieval call binding the contract method 0xf228d6f0.
//
// Solidity: function getEhrUser(bytes32 userIDHash) view returns(bytes32)
func (_EhrIndexer *EhrIndexerCallerSession) GetEhrUser(userIDHash [32]byte) ([32]byte, error) {
	return _EhrIndexer.Contract.GetEhrUser(&_EhrIndexer.CallOpts, userIDHash)
}

// GetLastEhrDocByType is a free data retrieval call binding the contract method 0x15dbcf1a.
//
// Solidity: function getLastEhrDocByType(bytes32 ehrId, uint8 docType) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_EhrIndexer *EhrIndexerCaller) GetLastEhrDocByType(opts *bind.CallOpts, ehrId [32]byte, docType uint8) (DocsDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getLastEhrDocByType", ehrId, docType)

	if err != nil {
		return *new(DocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(DocsDocumentMeta)).(*DocsDocumentMeta)

	return out0, err

}

// GetLastEhrDocByType is a free data retrieval call binding the contract method 0x15dbcf1a.
//
// Solidity: function getLastEhrDocByType(bytes32 ehrId, uint8 docType) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_EhrIndexer *EhrIndexerSession) GetLastEhrDocByType(ehrId [32]byte, docType uint8) (DocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetLastEhrDocByType(&_EhrIndexer.CallOpts, ehrId, docType)
}

// GetLastEhrDocByType is a free data retrieval call binding the contract method 0x15dbcf1a.
//
// Solidity: function getLastEhrDocByType(bytes32 ehrId, uint8 docType) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_EhrIndexer *EhrIndexerCallerSession) GetLastEhrDocByType(ehrId [32]byte, docType uint8) (DocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetLastEhrDocByType(&_EhrIndexer.CallOpts, ehrId, docType)
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

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_EhrIndexer *EhrIndexerCaller) Users(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "users")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_EhrIndexer *EhrIndexerSession) Users() (common.Address, error) {
	return _EhrIndexer.Contract.Users(&_EhrIndexer.CallOpts)
}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_EhrIndexer *EhrIndexerCallerSession) Users() (common.Address, error) {
	return _EhrIndexer.Contract.Users(&_EhrIndexer.CallOpts)
}

// AddEhrDoc is a paid mutator transaction binding the contract method 0x8708f61f.
//
// Solidity: function addEhrDoc((uint8,bytes,bytes,uint32,(uint8,bytes)[],address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerTransactor) AddEhrDoc(opts *bind.TransactOpts, p DocsAddEhrDocParams) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "addEhrDoc", p)
}

// AddEhrDoc is a paid mutator transaction binding the contract method 0x8708f61f.
//
// Solidity: function addEhrDoc((uint8,bytes,bytes,uint32,(uint8,bytes)[],address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerSession) AddEhrDoc(p DocsAddEhrDocParams) (*types.Transaction, error) {
	return _EhrIndexer.Contract.AddEhrDoc(&_EhrIndexer.TransactOpts, p)
}

// AddEhrDoc is a paid mutator transaction binding the contract method 0x8708f61f.
//
// Solidity: function addEhrDoc((uint8,bytes,bytes,uint32,(uint8,bytes)[],address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) AddEhrDoc(p DocsAddEhrDocParams) (*types.Transaction, error) {
	return _EhrIndexer.Contract.AddEhrDoc(&_EhrIndexer.TransactOpts, p)
}

// DeleteDoc is a paid mutator transaction binding the contract method 0x54d3f64f.
//
// Solidity: function deleteDoc(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactor) DeleteDoc(opts *bind.TransactOpts, ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "deleteDoc", ehrId, docType, docBaseUIDHash, version, signer, signature)
}

// DeleteDoc is a paid mutator transaction binding the contract method 0x54d3f64f.
//
// Solidity: function deleteDoc(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerSession) DeleteDoc(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.DeleteDoc(&_EhrIndexer.TransactOpts, ehrId, docType, docBaseUIDHash, version, signer, signature)
}

// DeleteDoc is a paid mutator transaction binding the contract method 0x54d3f64f.
//
// Solidity: function deleteDoc(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) DeleteDoc(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.DeleteDoc(&_EhrIndexer.TransactOpts, ehrId, docType, docBaseUIDHash, version, signer, signature)
}

// DocGroupAddDoc is a paid mutator transaction binding the contract method 0x14ce75d6.
//
// Solidity: function docGroupAddDoc(bytes32 groupIDHash, bytes32 docCIDHash, bytes docCIDEncr, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactor) DocGroupAddDoc(opts *bind.TransactOpts, groupIDHash [32]byte, docCIDHash [32]byte, docCIDEncr []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "docGroupAddDoc", groupIDHash, docCIDHash, docCIDEncr, signer, signature)
}

// DocGroupAddDoc is a paid mutator transaction binding the contract method 0x14ce75d6.
//
// Solidity: function docGroupAddDoc(bytes32 groupIDHash, bytes32 docCIDHash, bytes docCIDEncr, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerSession) DocGroupAddDoc(groupIDHash [32]byte, docCIDHash [32]byte, docCIDEncr []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.DocGroupAddDoc(&_EhrIndexer.TransactOpts, groupIDHash, docCIDHash, docCIDEncr, signer, signature)
}

// DocGroupAddDoc is a paid mutator transaction binding the contract method 0x14ce75d6.
//
// Solidity: function docGroupAddDoc(bytes32 groupIDHash, bytes32 docCIDHash, bytes docCIDEncr, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) DocGroupAddDoc(groupIDHash [32]byte, docCIDHash [32]byte, docCIDEncr []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.DocGroupAddDoc(&_EhrIndexer.TransactOpts, groupIDHash, docCIDHash, docCIDEncr, signer, signature)
}

// DocGroupCreate is a paid mutator transaction binding the contract method 0x88695a65.
//
// Solidity: function docGroupCreate((bytes32,(uint8,bytes)[],address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerTransactor) DocGroupCreate(opts *bind.TransactOpts, p DocGroupsDocGroupCreateParams) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "docGroupCreate", p)
}

// DocGroupCreate is a paid mutator transaction binding the contract method 0x88695a65.
//
// Solidity: function docGroupCreate((bytes32,(uint8,bytes)[],address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerSession) DocGroupCreate(p DocGroupsDocGroupCreateParams) (*types.Transaction, error) {
	return _EhrIndexer.Contract.DocGroupCreate(&_EhrIndexer.TransactOpts, p)
}

// DocGroupCreate is a paid mutator transaction binding the contract method 0x88695a65.
//
// Solidity: function docGroupCreate((bytes32,(uint8,bytes)[],address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) DocGroupCreate(p DocGroupsDocGroupCreateParams) (*types.Transaction, error) {
	return _EhrIndexer.Contract.DocGroupCreate(&_EhrIndexer.TransactOpts, p)
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

// SetDocAccess is a paid mutator transaction binding the contract method 0x3591ac33.
//
// Solidity: function setDocAccess(bytes32 CIDHash, (bytes32,bytes,bytes,uint8) access, address userAddr, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactor) SetDocAccess(opts *bind.TransactOpts, CIDHash [32]byte, access IAccessStoreAccess, userAddr common.Address, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setDocAccess", CIDHash, access, userAddr, signer, signature)
}

// SetDocAccess is a paid mutator transaction binding the contract method 0x3591ac33.
//
// Solidity: function setDocAccess(bytes32 CIDHash, (bytes32,bytes,bytes,uint8) access, address userAddr, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerSession) SetDocAccess(CIDHash [32]byte, access IAccessStoreAccess, userAddr common.Address, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetDocAccess(&_EhrIndexer.TransactOpts, CIDHash, access, userAddr, signer, signature)
}

// SetDocAccess is a paid mutator transaction binding the contract method 0x3591ac33.
//
// Solidity: function setDocAccess(bytes32 CIDHash, (bytes32,bytes,bytes,uint8) access, address userAddr, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) SetDocAccess(CIDHash [32]byte, access IAccessStoreAccess, userAddr common.Address, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetDocAccess(&_EhrIndexer.TransactOpts, CIDHash, access, userAddr, signer, signature)
}

// SetEhrSubject is a paid mutator transaction binding the contract method 0x9975202f.
//
// Solidity: function setEhrSubject(bytes32 subjectKey, bytes32 ehrId, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactor) SetEhrSubject(opts *bind.TransactOpts, subjectKey [32]byte, ehrId [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setEhrSubject", subjectKey, ehrId, signer, signature)
}

// SetEhrSubject is a paid mutator transaction binding the contract method 0x9975202f.
//
// Solidity: function setEhrSubject(bytes32 subjectKey, bytes32 ehrId, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerSession) SetEhrSubject(subjectKey [32]byte, ehrId [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrSubject(&_EhrIndexer.TransactOpts, subjectKey, ehrId, signer, signature)
}

// SetEhrSubject is a paid mutator transaction binding the contract method 0x9975202f.
//
// Solidity: function setEhrSubject(bytes32 subjectKey, bytes32 ehrId, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) SetEhrSubject(subjectKey [32]byte, ehrId [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrSubject(&_EhrIndexer.TransactOpts, subjectKey, ehrId, signer, signature)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0x3f157693.
//
// Solidity: function setEhrUser(bytes32 IDHash, bytes32 ehrId, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactor) SetEhrUser(opts *bind.TransactOpts, IDHash [32]byte, ehrId [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setEhrUser", IDHash, ehrId, signer, signature)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0x3f157693.
//
// Solidity: function setEhrUser(bytes32 IDHash, bytes32 ehrId, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerSession) SetEhrUser(IDHash [32]byte, ehrId [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrUser(&_EhrIndexer.TransactOpts, IDHash, ehrId, signer, signature)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0x3f157693.
//
// Solidity: function setEhrUser(bytes32 IDHash, bytes32 ehrId, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) SetEhrUser(IDHash [32]byte, ehrId [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrUser(&_EhrIndexer.TransactOpts, IDHash, ehrId, signer, signature)
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
