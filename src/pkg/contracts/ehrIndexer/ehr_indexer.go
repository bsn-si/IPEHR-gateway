// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ehrindexer

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

// EhrindexerMetaData contains all meta data concerning the Ehrindexer contract.
var EhrindexerMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_accessStore\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_users\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"accessStore\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structDocs.AddEhrDocParams\",\"name\":\"p\",\"type\":\"tuple\"}],\"name\":\"addEhrDoc\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedChange\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"deleteDoc\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"groupIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"docCIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"docCIDEncr\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"docGroupAddDoc\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"groupIDHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structDocGroups.DocGroupCreateParams\",\"name\":\"p\",\"type\":\"tuple\"}],\"name\":\"docGroupCreate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"groupIdHash\",\"type\":\"bytes32\"}],\"name\":\"docGroupGetAttrs\",\"outputs\":[{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"groupIdHash\",\"type\":\"bytes32\"}],\"name\":\"docGroupGetDocs\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"ehrSubject\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrID\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"name\":\"getDocByTime\",\"outputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structDocs.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"}],\"name\":\"getDocByVersion\",\"outputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structDocs.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"UIDHash\",\"type\":\"bytes32\"}],\"name\":\"getDocLastByBaseID\",\"outputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structDocs.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"}],\"name\":\"getEhrDocs\",\"outputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structDocs.DocumentMeta[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userIDHash\",\"type\":\"bytes32\"}],\"name\":\"getEhrUser\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"}],\"name\":\"getLastEhrDocByType\",\"outputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structDocs.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAllowed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"CIDHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumIAccessStore.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structIAccessStore.Access\",\"name\":\"access\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"userAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"setDocAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"subjectKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"setEhrSubject\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"IDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"setEhrUser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"users\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// EhrindexerABI is the input ABI used to generate the binding from.
// Deprecated: Use EhrindexerMetaData.ABI instead.
var EhrindexerABI = EhrindexerMetaData.ABI

// Ehrindexer is an auto generated Go binding around an Ethereum contract.
type Ehrindexer struct {
	EhrindexerCaller     // Read-only binding to the contract
	EhrindexerTransactor // Write-only binding to the contract
	EhrindexerFilterer   // Log filterer for contract events
}

// EhrindexerCaller is an auto generated read-only Go binding around an Ethereum contract.
type EhrindexerCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EhrindexerTransactor is an auto generated write-only Go binding around an Ethereum contract.
type EhrindexerTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EhrindexerFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type EhrindexerFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// EhrindexerSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type EhrindexerSession struct {
	Contract     *Ehrindexer       // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// EhrindexerCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type EhrindexerCallerSession struct {
	Contract *EhrindexerCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts     // Call options to use throughout this session
}

// EhrindexerTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type EhrindexerTransactorSession struct {
	Contract     *EhrindexerTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts     // Transaction auth options to use throughout this session
}

// EhrindexerRaw is an auto generated low-level Go binding around an Ethereum contract.
type EhrindexerRaw struct {
	Contract *Ehrindexer // Generic contract binding to access the raw methods on
}

// EhrindexerCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type EhrindexerCallerRaw struct {
	Contract *EhrindexerCaller // Generic read-only contract binding to access the raw methods on
}

// EhrindexerTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type EhrindexerTransactorRaw struct {
	Contract *EhrindexerTransactor // Generic write-only contract binding to access the raw methods on
}

// NewEhrindexer creates a new instance of Ehrindexer, bound to a specific deployed contract.
func NewEhrindexer(address common.Address, backend bind.ContractBackend) (*Ehrindexer, error) {
	contract, err := bindEhrindexer(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Ehrindexer{EhrindexerCaller: EhrindexerCaller{contract: contract}, EhrindexerTransactor: EhrindexerTransactor{contract: contract}, EhrindexerFilterer: EhrindexerFilterer{contract: contract}}, nil
}

// NewEhrindexerCaller creates a new read-only instance of Ehrindexer, bound to a specific deployed contract.
func NewEhrindexerCaller(address common.Address, caller bind.ContractCaller) (*EhrindexerCaller, error) {
	contract, err := bindEhrindexer(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &EhrindexerCaller{contract: contract}, nil
}

// NewEhrindexerTransactor creates a new write-only instance of Ehrindexer, bound to a specific deployed contract.
func NewEhrindexerTransactor(address common.Address, transactor bind.ContractTransactor) (*EhrindexerTransactor, error) {
	contract, err := bindEhrindexer(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &EhrindexerTransactor{contract: contract}, nil
}

// NewEhrindexerFilterer creates a new log filterer instance of Ehrindexer, bound to a specific deployed contract.
func NewEhrindexerFilterer(address common.Address, filterer bind.ContractFilterer) (*EhrindexerFilterer, error) {
	contract, err := bindEhrindexer(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &EhrindexerFilterer{contract: contract}, nil
}

// bindEhrindexer binds a generic wrapper to an already deployed contract.
func bindEhrindexer(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(EhrindexerABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ehrindexer *EhrindexerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ehrindexer.Contract.EhrindexerCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ehrindexer *EhrindexerRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ehrindexer.Contract.EhrindexerTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ehrindexer *EhrindexerRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ehrindexer.Contract.EhrindexerTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Ehrindexer *EhrindexerCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Ehrindexer.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Ehrindexer *EhrindexerTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Ehrindexer.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Ehrindexer *EhrindexerTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Ehrindexer.Contract.contract.Transact(opts, method, params...)
}

// AccessStore is a free data retrieval call binding the contract method 0x45e9ea6f.
//
// Solidity: function accessStore() view returns(address)
func (_Ehrindexer *EhrindexerCaller) AccessStore(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ehrindexer.contract.Call(opts, &out, "accessStore")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AccessStore is a free data retrieval call binding the contract method 0x45e9ea6f.
//
// Solidity: function accessStore() view returns(address)
func (_Ehrindexer *EhrindexerSession) AccessStore() (common.Address, error) {
	return _Ehrindexer.Contract.AccessStore(&_Ehrindexer.CallOpts)
}

// AccessStore is a free data retrieval call binding the contract method 0x45e9ea6f.
//
// Solidity: function accessStore() view returns(address)
func (_Ehrindexer *EhrindexerCallerSession) AccessStore() (common.Address, error) {
	return _Ehrindexer.Contract.AccessStore(&_Ehrindexer.CallOpts)
}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_Ehrindexer *EhrindexerCaller) AllowedChange(opts *bind.CallOpts, arg0 common.Address) (bool, error) {
	var out []interface{}
	err := _Ehrindexer.contract.Call(opts, &out, "allowedChange", arg0)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_Ehrindexer *EhrindexerSession) AllowedChange(arg0 common.Address) (bool, error) {
	return _Ehrindexer.Contract.AllowedChange(&_Ehrindexer.CallOpts, arg0)
}

// AllowedChange is a free data retrieval call binding the contract method 0xe9b5b29a.
//
// Solidity: function allowedChange(address ) view returns(bool)
func (_Ehrindexer *EhrindexerCallerSession) AllowedChange(arg0 common.Address) (bool, error) {
	return _Ehrindexer.Contract.AllowedChange(&_Ehrindexer.CallOpts, arg0)
}

// DocGroupGetAttrs is a free data retrieval call binding the contract method 0x891c952e.
//
// Solidity: function docGroupGetAttrs(bytes32 groupIdHash) view returns((uint8,bytes)[])
func (_Ehrindexer *EhrindexerCaller) DocGroupGetAttrs(opts *bind.CallOpts, groupIdHash [32]byte) ([]AttributesAttribute, error) {
	var out []interface{}
	err := _Ehrindexer.contract.Call(opts, &out, "docGroupGetAttrs", groupIdHash)

	if err != nil {
		return *new([]AttributesAttribute), err
	}

	out0 := *abi.ConvertType(out[0], new([]AttributesAttribute)).(*[]AttributesAttribute)

	return out0, err

}

// DocGroupGetAttrs is a free data retrieval call binding the contract method 0x891c952e.
//
// Solidity: function docGroupGetAttrs(bytes32 groupIdHash) view returns((uint8,bytes)[])
func (_Ehrindexer *EhrindexerSession) DocGroupGetAttrs(groupIdHash [32]byte) ([]AttributesAttribute, error) {
	return _Ehrindexer.Contract.DocGroupGetAttrs(&_Ehrindexer.CallOpts, groupIdHash)
}

// DocGroupGetAttrs is a free data retrieval call binding the contract method 0x891c952e.
//
// Solidity: function docGroupGetAttrs(bytes32 groupIdHash) view returns((uint8,bytes)[])
func (_Ehrindexer *EhrindexerCallerSession) DocGroupGetAttrs(groupIdHash [32]byte) ([]AttributesAttribute, error) {
	return _Ehrindexer.Contract.DocGroupGetAttrs(&_Ehrindexer.CallOpts, groupIdHash)
}

// DocGroupGetDocs is a free data retrieval call binding the contract method 0x88216834.
//
// Solidity: function docGroupGetDocs(bytes32 groupIdHash) view returns(bytes[])
func (_Ehrindexer *EhrindexerCaller) DocGroupGetDocs(opts *bind.CallOpts, groupIdHash [32]byte) ([][]byte, error) {
	var out []interface{}
	err := _Ehrindexer.contract.Call(opts, &out, "docGroupGetDocs", groupIdHash)

	if err != nil {
		return *new([][]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([][]byte)).(*[][]byte)

	return out0, err

}

// DocGroupGetDocs is a free data retrieval call binding the contract method 0x88216834.
//
// Solidity: function docGroupGetDocs(bytes32 groupIdHash) view returns(bytes[])
func (_Ehrindexer *EhrindexerSession) DocGroupGetDocs(groupIdHash [32]byte) ([][]byte, error) {
	return _Ehrindexer.Contract.DocGroupGetDocs(&_Ehrindexer.CallOpts, groupIdHash)
}

// DocGroupGetDocs is a free data retrieval call binding the contract method 0x88216834.
//
// Solidity: function docGroupGetDocs(bytes32 groupIdHash) view returns(bytes[])
func (_Ehrindexer *EhrindexerCallerSession) DocGroupGetDocs(groupIdHash [32]byte) ([][]byte, error) {
	return _Ehrindexer.Contract.DocGroupGetDocs(&_Ehrindexer.CallOpts, groupIdHash)
}

// EhrSubject is a free data retrieval call binding the contract method 0xfe1b5580.
//
// Solidity: function ehrSubject(bytes32 ) view returns(bytes32)
func (_Ehrindexer *EhrindexerCaller) EhrSubject(opts *bind.CallOpts, arg0 [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Ehrindexer.contract.Call(opts, &out, "ehrSubject", arg0)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// EhrSubject is a free data retrieval call binding the contract method 0xfe1b5580.
//
// Solidity: function ehrSubject(bytes32 ) view returns(bytes32)
func (_Ehrindexer *EhrindexerSession) EhrSubject(arg0 [32]byte) ([32]byte, error) {
	return _Ehrindexer.Contract.EhrSubject(&_Ehrindexer.CallOpts, arg0)
}

// EhrSubject is a free data retrieval call binding the contract method 0xfe1b5580.
//
// Solidity: function ehrSubject(bytes32 ) view returns(bytes32)
func (_Ehrindexer *EhrindexerCallerSession) EhrSubject(arg0 [32]byte) ([32]byte, error) {
	return _Ehrindexer.Contract.EhrSubject(&_Ehrindexer.CallOpts, arg0)
}

// GetDocByTime is a free data retrieval call binding the contract method 0x4f722b16.
//
// Solidity: function getDocByTime(bytes32 ehrID, uint8 docType, uint32 timestamp) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_Ehrindexer *EhrindexerCaller) GetDocByTime(opts *bind.CallOpts, ehrID [32]byte, docType uint8, timestamp uint32) (DocsDocumentMeta, error) {
	var out []interface{}
	err := _Ehrindexer.contract.Call(opts, &out, "getDocByTime", ehrID, docType, timestamp)

	if err != nil {
		return *new(DocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(DocsDocumentMeta)).(*DocsDocumentMeta)

	return out0, err

}

// GetDocByTime is a free data retrieval call binding the contract method 0x4f722b16.
//
// Solidity: function getDocByTime(bytes32 ehrID, uint8 docType, uint32 timestamp) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_Ehrindexer *EhrindexerSession) GetDocByTime(ehrID [32]byte, docType uint8, timestamp uint32) (DocsDocumentMeta, error) {
	return _Ehrindexer.Contract.GetDocByTime(&_Ehrindexer.CallOpts, ehrID, docType, timestamp)
}

// GetDocByTime is a free data retrieval call binding the contract method 0x4f722b16.
//
// Solidity: function getDocByTime(bytes32 ehrID, uint8 docType, uint32 timestamp) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_Ehrindexer *EhrindexerCallerSession) GetDocByTime(ehrID [32]byte, docType uint8, timestamp uint32) (DocsDocumentMeta, error) {
	return _Ehrindexer.Contract.GetDocByTime(&_Ehrindexer.CallOpts, ehrID, docType, timestamp)
}

// GetDocByVersion is a free data retrieval call binding the contract method 0x179fcacf.
//
// Solidity: function getDocByVersion(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_Ehrindexer *EhrindexerCaller) GetDocByVersion(opts *bind.CallOpts, ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte) (DocsDocumentMeta, error) {
	var out []interface{}
	err := _Ehrindexer.contract.Call(opts, &out, "getDocByVersion", ehrId, docType, docBaseUIDHash, version)

	if err != nil {
		return *new(DocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(DocsDocumentMeta)).(*DocsDocumentMeta)

	return out0, err

}

// GetDocByVersion is a free data retrieval call binding the contract method 0x179fcacf.
//
// Solidity: function getDocByVersion(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_Ehrindexer *EhrindexerSession) GetDocByVersion(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte) (DocsDocumentMeta, error) {
	return _Ehrindexer.Contract.GetDocByVersion(&_Ehrindexer.CallOpts, ehrId, docType, docBaseUIDHash, version)
}

// GetDocByVersion is a free data retrieval call binding the contract method 0x179fcacf.
//
// Solidity: function getDocByVersion(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_Ehrindexer *EhrindexerCallerSession) GetDocByVersion(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte) (DocsDocumentMeta, error) {
	return _Ehrindexer.Contract.GetDocByVersion(&_Ehrindexer.CallOpts, ehrId, docType, docBaseUIDHash, version)
}

// GetDocLastByBaseID is a free data retrieval call binding the contract method 0x949c8e77.
//
// Solidity: function getDocLastByBaseID(bytes32 userIDHash, uint8 docType, bytes32 UIDHash) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_Ehrindexer *EhrindexerCaller) GetDocLastByBaseID(opts *bind.CallOpts, userIDHash [32]byte, docType uint8, UIDHash [32]byte) (DocsDocumentMeta, error) {
	var out []interface{}
	err := _Ehrindexer.contract.Call(opts, &out, "getDocLastByBaseID", userIDHash, docType, UIDHash)

	if err != nil {
		return *new(DocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(DocsDocumentMeta)).(*DocsDocumentMeta)

	return out0, err

}

// GetDocLastByBaseID is a free data retrieval call binding the contract method 0x949c8e77.
//
// Solidity: function getDocLastByBaseID(bytes32 userIDHash, uint8 docType, bytes32 UIDHash) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_Ehrindexer *EhrindexerSession) GetDocLastByBaseID(userIDHash [32]byte, docType uint8, UIDHash [32]byte) (DocsDocumentMeta, error) {
	return _Ehrindexer.Contract.GetDocLastByBaseID(&_Ehrindexer.CallOpts, userIDHash, docType, UIDHash)
}

// GetDocLastByBaseID is a free data retrieval call binding the contract method 0x949c8e77.
//
// Solidity: function getDocLastByBaseID(bytes32 userIDHash, uint8 docType, bytes32 UIDHash) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_Ehrindexer *EhrindexerCallerSession) GetDocLastByBaseID(userIDHash [32]byte, docType uint8, UIDHash [32]byte) (DocsDocumentMeta, error) {
	return _Ehrindexer.Contract.GetDocLastByBaseID(&_Ehrindexer.CallOpts, userIDHash, docType, UIDHash)
}

// GetEhrDocs is a free data retrieval call binding the contract method 0xeec6b331.
//
// Solidity: function getEhrDocs(bytes32 userIDHash, uint8 docType) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[])[])
func (_Ehrindexer *EhrindexerCaller) GetEhrDocs(opts *bind.CallOpts, userIDHash [32]byte, docType uint8) ([]DocsDocumentMeta, error) {
	var out []interface{}
	err := _Ehrindexer.contract.Call(opts, &out, "getEhrDocs", userIDHash, docType)

	if err != nil {
		return *new([]DocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new([]DocsDocumentMeta)).(*[]DocsDocumentMeta)

	return out0, err

}

// GetEhrDocs is a free data retrieval call binding the contract method 0xeec6b331.
//
// Solidity: function getEhrDocs(bytes32 userIDHash, uint8 docType) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[])[])
func (_Ehrindexer *EhrindexerSession) GetEhrDocs(userIDHash [32]byte, docType uint8) ([]DocsDocumentMeta, error) {
	return _Ehrindexer.Contract.GetEhrDocs(&_Ehrindexer.CallOpts, userIDHash, docType)
}

// GetEhrDocs is a free data retrieval call binding the contract method 0xeec6b331.
//
// Solidity: function getEhrDocs(bytes32 userIDHash, uint8 docType) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[])[])
func (_Ehrindexer *EhrindexerCallerSession) GetEhrDocs(userIDHash [32]byte, docType uint8) ([]DocsDocumentMeta, error) {
	return _Ehrindexer.Contract.GetEhrDocs(&_Ehrindexer.CallOpts, userIDHash, docType)
}

// GetEhrUser is a free data retrieval call binding the contract method 0xf228d6f0.
//
// Solidity: function getEhrUser(bytes32 userIDHash) view returns(bytes32)
func (_Ehrindexer *EhrindexerCaller) GetEhrUser(opts *bind.CallOpts, userIDHash [32]byte) ([32]byte, error) {
	var out []interface{}
	err := _Ehrindexer.contract.Call(opts, &out, "getEhrUser", userIDHash)

	if err != nil {
		return *new([32]byte), err
	}

	out0 := *abi.ConvertType(out[0], new([32]byte)).(*[32]byte)

	return out0, err

}

// GetEhrUser is a free data retrieval call binding the contract method 0xf228d6f0.
//
// Solidity: function getEhrUser(bytes32 userIDHash) view returns(bytes32)
func (_Ehrindexer *EhrindexerSession) GetEhrUser(userIDHash [32]byte) ([32]byte, error) {
	return _Ehrindexer.Contract.GetEhrUser(&_Ehrindexer.CallOpts, userIDHash)
}

// GetEhrUser is a free data retrieval call binding the contract method 0xf228d6f0.
//
// Solidity: function getEhrUser(bytes32 userIDHash) view returns(bytes32)
func (_Ehrindexer *EhrindexerCallerSession) GetEhrUser(userIDHash [32]byte) ([32]byte, error) {
	return _Ehrindexer.Contract.GetEhrUser(&_Ehrindexer.CallOpts, userIDHash)
}

// GetLastEhrDocByType is a free data retrieval call binding the contract method 0x15dbcf1a.
//
// Solidity: function getLastEhrDocByType(bytes32 ehrId, uint8 docType) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_Ehrindexer *EhrindexerCaller) GetLastEhrDocByType(opts *bind.CallOpts, ehrId [32]byte, docType uint8) (DocsDocumentMeta, error) {
	var out []interface{}
	err := _Ehrindexer.contract.Call(opts, &out, "getLastEhrDocByType", ehrId, docType)

	if err != nil {
		return *new(DocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(DocsDocumentMeta)).(*DocsDocumentMeta)

	return out0, err

}

// GetLastEhrDocByType is a free data retrieval call binding the contract method 0x15dbcf1a.
//
// Solidity: function getLastEhrDocByType(bytes32 ehrId, uint8 docType) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_Ehrindexer *EhrindexerSession) GetLastEhrDocByType(ehrId [32]byte, docType uint8) (DocsDocumentMeta, error) {
	return _Ehrindexer.Contract.GetLastEhrDocByType(&_Ehrindexer.CallOpts, ehrId, docType)
}

// GetLastEhrDocByType is a free data retrieval call binding the contract method 0x15dbcf1a.
//
// Solidity: function getLastEhrDocByType(bytes32 ehrId, uint8 docType) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_Ehrindexer *EhrindexerCallerSession) GetLastEhrDocByType(ehrId [32]byte, docType uint8) (DocsDocumentMeta, error) {
	return _Ehrindexer.Contract.GetLastEhrDocByType(&_Ehrindexer.CallOpts, ehrId, docType)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_Ehrindexer *EhrindexerCaller) Nonces(opts *bind.CallOpts, arg0 common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Ehrindexer.contract.Call(opts, &out, "nonces", arg0)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_Ehrindexer *EhrindexerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _Ehrindexer.Contract.Nonces(&_Ehrindexer.CallOpts, arg0)
}

// Nonces is a free data retrieval call binding the contract method 0x7ecebe00.
//
// Solidity: function nonces(address ) view returns(uint256)
func (_Ehrindexer *EhrindexerCallerSession) Nonces(arg0 common.Address) (*big.Int, error) {
	return _Ehrindexer.Contract.Nonces(&_Ehrindexer.CallOpts, arg0)
}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_Ehrindexer *EhrindexerCaller) Users(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Ehrindexer.contract.Call(opts, &out, "users")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_Ehrindexer *EhrindexerSession) Users() (common.Address, error) {
	return _Ehrindexer.Contract.Users(&_Ehrindexer.CallOpts)
}

// Users is a free data retrieval call binding the contract method 0xf2020275.
//
// Solidity: function users() view returns(address)
func (_Ehrindexer *EhrindexerCallerSession) Users() (common.Address, error) {
	return _Ehrindexer.Contract.Users(&_Ehrindexer.CallOpts)
}

// AddEhrDoc is a paid mutator transaction binding the contract method 0x8708f61f.
//
// Solidity: function addEhrDoc((uint8,bytes,bytes,uint32,(uint8,bytes)[],address,bytes) p) returns()
func (_Ehrindexer *EhrindexerTransactor) AddEhrDoc(opts *bind.TransactOpts, p DocsAddEhrDocParams) (*types.Transaction, error) {
	return _Ehrindexer.contract.Transact(opts, "addEhrDoc", p)
}

// AddEhrDoc is a paid mutator transaction binding the contract method 0x8708f61f.
//
// Solidity: function addEhrDoc((uint8,bytes,bytes,uint32,(uint8,bytes)[],address,bytes) p) returns()
func (_Ehrindexer *EhrindexerSession) AddEhrDoc(p DocsAddEhrDocParams) (*types.Transaction, error) {
	return _Ehrindexer.Contract.AddEhrDoc(&_Ehrindexer.TransactOpts, p)
}

// AddEhrDoc is a paid mutator transaction binding the contract method 0x8708f61f.
//
// Solidity: function addEhrDoc((uint8,bytes,bytes,uint32,(uint8,bytes)[],address,bytes) p) returns()
func (_Ehrindexer *EhrindexerTransactorSession) AddEhrDoc(p DocsAddEhrDocParams) (*types.Transaction, error) {
	return _Ehrindexer.Contract.AddEhrDoc(&_Ehrindexer.TransactOpts, p)
}

// DeleteDoc is a paid mutator transaction binding the contract method 0x54d3f64f.
//
// Solidity: function deleteDoc(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version, address signer, bytes signature) returns()
func (_Ehrindexer *EhrindexerTransactor) DeleteDoc(opts *bind.TransactOpts, ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Ehrindexer.contract.Transact(opts, "deleteDoc", ehrId, docType, docBaseUIDHash, version, signer, signature)
}

// DeleteDoc is a paid mutator transaction binding the contract method 0x54d3f64f.
//
// Solidity: function deleteDoc(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version, address signer, bytes signature) returns()
func (_Ehrindexer *EhrindexerSession) DeleteDoc(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Ehrindexer.Contract.DeleteDoc(&_Ehrindexer.TransactOpts, ehrId, docType, docBaseUIDHash, version, signer, signature)
}

// DeleteDoc is a paid mutator transaction binding the contract method 0x54d3f64f.
//
// Solidity: function deleteDoc(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash, bytes32 version, address signer, bytes signature) returns()
func (_Ehrindexer *EhrindexerTransactorSession) DeleteDoc(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte, version [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Ehrindexer.Contract.DeleteDoc(&_Ehrindexer.TransactOpts, ehrId, docType, docBaseUIDHash, version, signer, signature)
}

// DocGroupAddDoc is a paid mutator transaction binding the contract method 0x14ce75d6.
//
// Solidity: function docGroupAddDoc(bytes32 groupIDHash, bytes32 docCIDHash, bytes docCIDEncr, address signer, bytes signature) returns()
func (_Ehrindexer *EhrindexerTransactor) DocGroupAddDoc(opts *bind.TransactOpts, groupIDHash [32]byte, docCIDHash [32]byte, docCIDEncr []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Ehrindexer.contract.Transact(opts, "docGroupAddDoc", groupIDHash, docCIDHash, docCIDEncr, signer, signature)
}

// DocGroupAddDoc is a paid mutator transaction binding the contract method 0x14ce75d6.
//
// Solidity: function docGroupAddDoc(bytes32 groupIDHash, bytes32 docCIDHash, bytes docCIDEncr, address signer, bytes signature) returns()
func (_Ehrindexer *EhrindexerSession) DocGroupAddDoc(groupIDHash [32]byte, docCIDHash [32]byte, docCIDEncr []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Ehrindexer.Contract.DocGroupAddDoc(&_Ehrindexer.TransactOpts, groupIDHash, docCIDHash, docCIDEncr, signer, signature)
}

// DocGroupAddDoc is a paid mutator transaction binding the contract method 0x14ce75d6.
//
// Solidity: function docGroupAddDoc(bytes32 groupIDHash, bytes32 docCIDHash, bytes docCIDEncr, address signer, bytes signature) returns()
func (_Ehrindexer *EhrindexerTransactorSession) DocGroupAddDoc(groupIDHash [32]byte, docCIDHash [32]byte, docCIDEncr []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Ehrindexer.Contract.DocGroupAddDoc(&_Ehrindexer.TransactOpts, groupIDHash, docCIDHash, docCIDEncr, signer, signature)
}

// DocGroupCreate is a paid mutator transaction binding the contract method 0x88695a65.
//
// Solidity: function docGroupCreate((bytes32,(uint8,bytes)[],address,bytes) p) returns()
func (_Ehrindexer *EhrindexerTransactor) DocGroupCreate(opts *bind.TransactOpts, p DocGroupsDocGroupCreateParams) (*types.Transaction, error) {
	return _Ehrindexer.contract.Transact(opts, "docGroupCreate", p)
}

// DocGroupCreate is a paid mutator transaction binding the contract method 0x88695a65.
//
// Solidity: function docGroupCreate((bytes32,(uint8,bytes)[],address,bytes) p) returns()
func (_Ehrindexer *EhrindexerSession) DocGroupCreate(p DocGroupsDocGroupCreateParams) (*types.Transaction, error) {
	return _Ehrindexer.Contract.DocGroupCreate(&_Ehrindexer.TransactOpts, p)
}

// DocGroupCreate is a paid mutator transaction binding the contract method 0x88695a65.
//
// Solidity: function docGroupCreate((bytes32,(uint8,bytes)[],address,bytes) p) returns()
func (_Ehrindexer *EhrindexerTransactorSession) DocGroupCreate(p DocGroupsDocGroupCreateParams) (*types.Transaction, error) {
	return _Ehrindexer.Contract.DocGroupCreate(&_Ehrindexer.TransactOpts, p)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Ehrindexer *EhrindexerTransactor) Multicall(opts *bind.TransactOpts, data [][]byte) (*types.Transaction, error) {
	return _Ehrindexer.contract.Transact(opts, "multicall", data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Ehrindexer *EhrindexerSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _Ehrindexer.Contract.Multicall(&_Ehrindexer.TransactOpts, data)
}

// Multicall is a paid mutator transaction binding the contract method 0xac9650d8.
//
// Solidity: function multicall(bytes[] data) returns(bytes[] results)
func (_Ehrindexer *EhrindexerTransactorSession) Multicall(data [][]byte) (*types.Transaction, error) {
	return _Ehrindexer.Contract.Multicall(&_Ehrindexer.TransactOpts, data)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_Ehrindexer *EhrindexerTransactor) SetAllowed(opts *bind.TransactOpts, addr common.Address, allowed bool) (*types.Transaction, error) {
	return _Ehrindexer.contract.Transact(opts, "setAllowed", addr, allowed)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_Ehrindexer *EhrindexerSession) SetAllowed(addr common.Address, allowed bool) (*types.Transaction, error) {
	return _Ehrindexer.Contract.SetAllowed(&_Ehrindexer.TransactOpts, addr, allowed)
}

// SetAllowed is a paid mutator transaction binding the contract method 0x4697f05d.
//
// Solidity: function setAllowed(address addr, bool allowed) returns()
func (_Ehrindexer *EhrindexerTransactorSession) SetAllowed(addr common.Address, allowed bool) (*types.Transaction, error) {
	return _Ehrindexer.Contract.SetAllowed(&_Ehrindexer.TransactOpts, addr, allowed)
}

// SetDocAccess is a paid mutator transaction binding the contract method 0x3591ac33.
//
// Solidity: function setDocAccess(bytes32 CIDHash, (bytes32,bytes,bytes,uint8) access, address userAddr, address signer, bytes signature) returns()
func (_Ehrindexer *EhrindexerTransactor) SetDocAccess(opts *bind.TransactOpts, CIDHash [32]byte, access IAccessStoreAccess, userAddr common.Address, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Ehrindexer.contract.Transact(opts, "setDocAccess", CIDHash, access, userAddr, signer, signature)
}

// SetDocAccess is a paid mutator transaction binding the contract method 0x3591ac33.
//
// Solidity: function setDocAccess(bytes32 CIDHash, (bytes32,bytes,bytes,uint8) access, address userAddr, address signer, bytes signature) returns()
func (_Ehrindexer *EhrindexerSession) SetDocAccess(CIDHash [32]byte, access IAccessStoreAccess, userAddr common.Address, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Ehrindexer.Contract.SetDocAccess(&_Ehrindexer.TransactOpts, CIDHash, access, userAddr, signer, signature)
}

// SetDocAccess is a paid mutator transaction binding the contract method 0x3591ac33.
//
// Solidity: function setDocAccess(bytes32 CIDHash, (bytes32,bytes,bytes,uint8) access, address userAddr, address signer, bytes signature) returns()
func (_Ehrindexer *EhrindexerTransactorSession) SetDocAccess(CIDHash [32]byte, access IAccessStoreAccess, userAddr common.Address, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Ehrindexer.Contract.SetDocAccess(&_Ehrindexer.TransactOpts, CIDHash, access, userAddr, signer, signature)
}

// SetEhrSubject is a paid mutator transaction binding the contract method 0x9975202f.
//
// Solidity: function setEhrSubject(bytes32 subjectKey, bytes32 ehrId, address signer, bytes signature) returns()
func (_Ehrindexer *EhrindexerTransactor) SetEhrSubject(opts *bind.TransactOpts, subjectKey [32]byte, ehrId [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Ehrindexer.contract.Transact(opts, "setEhrSubject", subjectKey, ehrId, signer, signature)
}

// SetEhrSubject is a paid mutator transaction binding the contract method 0x9975202f.
//
// Solidity: function setEhrSubject(bytes32 subjectKey, bytes32 ehrId, address signer, bytes signature) returns()
func (_Ehrindexer *EhrindexerSession) SetEhrSubject(subjectKey [32]byte, ehrId [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Ehrindexer.Contract.SetEhrSubject(&_Ehrindexer.TransactOpts, subjectKey, ehrId, signer, signature)
}

// SetEhrSubject is a paid mutator transaction binding the contract method 0x9975202f.
//
// Solidity: function setEhrSubject(bytes32 subjectKey, bytes32 ehrId, address signer, bytes signature) returns()
func (_Ehrindexer *EhrindexerTransactorSession) SetEhrSubject(subjectKey [32]byte, ehrId [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Ehrindexer.Contract.SetEhrSubject(&_Ehrindexer.TransactOpts, subjectKey, ehrId, signer, signature)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0x3f157693.
//
// Solidity: function setEhrUser(bytes32 IDHash, bytes32 ehrId, address signer, bytes signature) returns()
func (_Ehrindexer *EhrindexerTransactor) SetEhrUser(opts *bind.TransactOpts, IDHash [32]byte, ehrId [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Ehrindexer.contract.Transact(opts, "setEhrUser", IDHash, ehrId, signer, signature)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0x3f157693.
//
// Solidity: function setEhrUser(bytes32 IDHash, bytes32 ehrId, address signer, bytes signature) returns()
func (_Ehrindexer *EhrindexerSession) SetEhrUser(IDHash [32]byte, ehrId [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Ehrindexer.Contract.SetEhrUser(&_Ehrindexer.TransactOpts, IDHash, ehrId, signer, signature)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0x3f157693.
//
// Solidity: function setEhrUser(bytes32 IDHash, bytes32 ehrId, address signer, bytes signature) returns()
func (_Ehrindexer *EhrindexerTransactorSession) SetEhrUser(IDHash [32]byte, ehrId [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _Ehrindexer.Contract.SetEhrUser(&_Ehrindexer.TransactOpts, IDHash, ehrId, signer, signature)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ehrindexer *EhrindexerTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _Ehrindexer.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ehrindexer *EhrindexerSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ehrindexer.Contract.TransferOwnership(&_Ehrindexer.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_Ehrindexer *EhrindexerTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _Ehrindexer.Contract.TransferOwnership(&_Ehrindexer.TransactOpts, newOwner)
}
