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

// AccessStoreAccess is an auto generated low-level Go binding around an user-defined struct.
type AccessStoreAccess struct {
	IdHash  [32]byte
	IdEncr  []byte
	KeyEncr []byte
	Level   uint8
}

// AttributesAttribute is an auto generated low-level Go binding around an user-defined struct.
type AttributesAttribute struct {
	Code  uint8
	Value []byte
}

// DocGroupsDocGroupCreateParams is an auto generated low-level Go binding around an user-defined struct.
type DocGroupsDocGroupCreateParams struct {
	GroupIdHash [32]byte
	GroupIdEncr []byte
	KeyEncr     []byte
	UserIdEncr  []byte
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

// UsersGroupAddUserParams is an auto generated low-level Go binding around an user-defined struct.
type UsersGroupAddUserParams struct {
	GroupIDHash [32]byte
	UserIDHash  [32]byte
	Level       uint8
	UserIDEncr  []byte
	KeyEncr     []byte
	Signer      common.Address
	Signature   []byte
}

// UsersGroupMember is an auto generated low-level Go binding around an user-defined struct.
type UsersGroupMember struct {
	UserIDHash [32]byte
	UserIDEncr []byte
}

// UsersUserGroup is an auto generated low-level Go binding around an user-defined struct.
type UsersUserGroup struct {
	Attrs   []AttributesAttribute
	Members []UsersGroupMember
}

// UsersUserGroupCreateParams is an auto generated low-level Go binding around an user-defined struct.
type UsersUserGroupCreateParams struct {
	GroupIdHash [32]byte
	Attrs       []AttributesAttribute
	Signer      common.Address
	Signature   []byte
}

// EhrIndexerMetaData contains all meta data concerning the EhrIndexer contract.
var EhrIndexerMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"inputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structDocs.AddEhrDocParams\",\"name\":\"p\",\"type\":\"tuple\"}],\"name\":\"addEhrDoc\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"allowedChange\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"deleteDoc\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"groupIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"CIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"CIDEncr\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"docGroupAddDoc\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"groupIdHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"groupIdEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"userIdEncr\",\"type\":\"bytes\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structDocGroups.DocGroupCreateParams\",\"name\":\"p\",\"type\":\"tuple\"}],\"name\":\"docGroupCreate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"groupIdHash\",\"type\":\"bytes32\"}],\"name\":\"docGroupGetDocs\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"\",\"type\":\"bytes[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"ehrSubject\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"name\":\"ehrUsers\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"accessID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"accessIdHash\",\"type\":\"bytes32\"}],\"name\":\"getAccessByIdHash\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumAccessStore.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structAccessStore.Access\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrID\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"}],\"name\":\"getDocByTime\",\"outputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structDocs.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"version\",\"type\":\"bytes32\"}],\"name\":\"getDocByVersion\",\"outputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structDocs.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"docBaseUIDHash\",\"type\":\"bytes32\"}],\"name\":\"getDocLastByBaseID\",\"outputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structDocs.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"}],\"name\":\"getEhrDocs\",\"outputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structDocs.DocumentMeta[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"enumDocs.DocType\",\"name\":\"docType\",\"type\":\"uint8\"}],\"name\":\"getLastEhrDocByType\",\"outputs\":[{\"components\":[{\"internalType\":\"enumDocs.DocStatus\",\"name\":\"status\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"id\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"version\",\"type\":\"bytes\"},{\"internalType\":\"uint32\",\"name\":\"timestamp\",\"type\":\"uint32\"},{\"internalType\":\"bool\",\"name\":\"isLast\",\"type\":\"bool\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"}],\"internalType\":\"structDocs.DocumentMeta\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"accessID\",\"type\":\"bytes32\"}],\"name\":\"getUserAccessList\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumAccessStore.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structAccessStore.Access[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"groupIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"userIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"enumAccessStore.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"userIDEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUsers.GroupAddUserParams\",\"name\":\"p\",\"type\":\"tuple\"}],\"name\":\"groupAddUser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"groupIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"userIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"groupRemoveUser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes[]\",\"name\":\"data\",\"type\":\"bytes[]\"}],\"name\":\"multicall\",\"outputs\":[{\"internalType\":\"bytes[]\",\"name\":\"results\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"nonces\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"addr\",\"type\":\"address\"},{\"internalType\":\"bool\",\"name\":\"allowed\",\"type\":\"bool\"}],\"name\":\"setAllowed\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"CIDHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumAccessStore.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structAccessStore.Access\",\"name\":\"access\",\"type\":\"tuple\"},{\"internalType\":\"address\",\"name\":\"userAddr\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"setDocAccess\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"subjectKey\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"setEhrSubject\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userId\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"ehrId\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"setEhrUser\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"userID\",\"type\":\"bytes32\"},{\"internalType\":\"enumAccessStore.AccessKind\",\"name\":\"kind\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"}],\"name\":\"userAccess\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"idHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"idEncr\",\"type\":\"bytes\"},{\"internalType\":\"bytes\",\"name\":\"keyEncr\",\"type\":\"bytes\"},{\"internalType\":\"enumAccessStore.AccessLevel\",\"name\":\"level\",\"type\":\"uint8\"}],\"internalType\":\"structAccessStore.Access\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"groupIdHash\",\"type\":\"bytes32\"},{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"internalType\":\"structUsers.UserGroupCreateParams\",\"name\":\"p\",\"type\":\"tuple\"}],\"name\":\"userGroupCreate\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"bytes32\",\"name\":\"groupIdHash\",\"type\":\"bytes32\"}],\"name\":\"userGroupGetByID\",\"outputs\":[{\"components\":[{\"components\":[{\"internalType\":\"enumAttributes.Code\",\"name\":\"code\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"value\",\"type\":\"bytes\"}],\"internalType\":\"structAttributes.Attribute[]\",\"name\":\"attrs\",\"type\":\"tuple[]\"},{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"userIDHash\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"userIDEncr\",\"type\":\"bytes\"}],\"internalType\":\"structUsers.GroupMember[]\",\"name\":\"members\",\"type\":\"tuple[]\"}],\"internalType\":\"structUsers.UserGroup\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"userAddr\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"systemID\",\"type\":\"bytes32\"},{\"internalType\":\"enumUsers.Role\",\"name\":\"role\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"pwdHash\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"signer\",\"type\":\"address\"},{\"internalType\":\"bytes\",\"name\":\"signature\",\"type\":\"bytes\"}],\"name\":\"userNew\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"users\",\"outputs\":[{\"internalType\":\"bytes32\",\"name\":\"id\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"systemID\",\"type\":\"bytes32\"},{\"internalType\":\"enumUsers.Role\",\"name\":\"role\",\"type\":\"uint8\"},{\"internalType\":\"bytes\",\"name\":\"pwdHash\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
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
// Solidity: function getAccessByIdHash(bytes32 accessID, bytes32 accessIdHash) view returns((bytes32,bytes,bytes,uint8))
func (_EhrIndexer *EhrIndexerCaller) GetAccessByIdHash(opts *bind.CallOpts, accessID [32]byte, accessIdHash [32]byte) (AccessStoreAccess, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getAccessByIdHash", accessID, accessIdHash)

	if err != nil {
		return *new(AccessStoreAccess), err
	}

	out0 := *abi.ConvertType(out[0], new(AccessStoreAccess)).(*AccessStoreAccess)

	return out0, err

}

// GetAccessByIdHash is a free data retrieval call binding the contract method 0x9ae2da76.
//
// Solidity: function getAccessByIdHash(bytes32 accessID, bytes32 accessIdHash) view returns((bytes32,bytes,bytes,uint8))
func (_EhrIndexer *EhrIndexerSession) GetAccessByIdHash(accessID [32]byte, accessIdHash [32]byte) (AccessStoreAccess, error) {
	return _EhrIndexer.Contract.GetAccessByIdHash(&_EhrIndexer.CallOpts, accessID, accessIdHash)
}

// GetAccessByIdHash is a free data retrieval call binding the contract method 0x9ae2da76.
//
// Solidity: function getAccessByIdHash(bytes32 accessID, bytes32 accessIdHash) view returns((bytes32,bytes,bytes,uint8))
func (_EhrIndexer *EhrIndexerCallerSession) GetAccessByIdHash(accessID [32]byte, accessIdHash [32]byte) (AccessStoreAccess, error) {
	return _EhrIndexer.Contract.GetAccessByIdHash(&_EhrIndexer.CallOpts, accessID, accessIdHash)
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
// Solidity: function getDocLastByBaseID(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_EhrIndexer *EhrIndexerCaller) GetDocLastByBaseID(opts *bind.CallOpts, ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte) (DocsDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getDocLastByBaseID", ehrId, docType, docBaseUIDHash)

	if err != nil {
		return *new(DocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new(DocsDocumentMeta)).(*DocsDocumentMeta)

	return out0, err

}

// GetDocLastByBaseID is a free data retrieval call binding the contract method 0x949c8e77.
//
// Solidity: function getDocLastByBaseID(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_EhrIndexer *EhrIndexerSession) GetDocLastByBaseID(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte) (DocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocLastByBaseID(&_EhrIndexer.CallOpts, ehrId, docType, docBaseUIDHash)
}

// GetDocLastByBaseID is a free data retrieval call binding the contract method 0x949c8e77.
//
// Solidity: function getDocLastByBaseID(bytes32 ehrId, uint8 docType, bytes32 docBaseUIDHash) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[]))
func (_EhrIndexer *EhrIndexerCallerSession) GetDocLastByBaseID(ehrId [32]byte, docType uint8, docBaseUIDHash [32]byte) (DocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetDocLastByBaseID(&_EhrIndexer.CallOpts, ehrId, docType, docBaseUIDHash)
}

// GetEhrDocs is a free data retrieval call binding the contract method 0xeec6b331.
//
// Solidity: function getEhrDocs(bytes32 ehrId, uint8 docType) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[])[])
func (_EhrIndexer *EhrIndexerCaller) GetEhrDocs(opts *bind.CallOpts, ehrId [32]byte, docType uint8) ([]DocsDocumentMeta, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getEhrDocs", ehrId, docType)

	if err != nil {
		return *new([]DocsDocumentMeta), err
	}

	out0 := *abi.ConvertType(out[0], new([]DocsDocumentMeta)).(*[]DocsDocumentMeta)

	return out0, err

}

// GetEhrDocs is a free data retrieval call binding the contract method 0xeec6b331.
//
// Solidity: function getEhrDocs(bytes32 ehrId, uint8 docType) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[])[])
func (_EhrIndexer *EhrIndexerSession) GetEhrDocs(ehrId [32]byte, docType uint8) ([]DocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetEhrDocs(&_EhrIndexer.CallOpts, ehrId, docType)
}

// GetEhrDocs is a free data retrieval call binding the contract method 0xeec6b331.
//
// Solidity: function getEhrDocs(bytes32 ehrId, uint8 docType) view returns((uint8,bytes,bytes,uint32,bool,(uint8,bytes)[])[])
func (_EhrIndexer *EhrIndexerCallerSession) GetEhrDocs(ehrId [32]byte, docType uint8) ([]DocsDocumentMeta, error) {
	return _EhrIndexer.Contract.GetEhrDocs(&_EhrIndexer.CallOpts, ehrId, docType)
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

// GetUserAccessList is a free data retrieval call binding the contract method 0xbb059b5c.
//
// Solidity: function getUserAccessList(bytes32 accessID) view returns((bytes32,bytes,bytes,uint8)[])
func (_EhrIndexer *EhrIndexerCaller) GetUserAccessList(opts *bind.CallOpts, accessID [32]byte) ([]AccessStoreAccess, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "getUserAccessList", accessID)

	if err != nil {
		return *new([]AccessStoreAccess), err
	}

	out0 := *abi.ConvertType(out[0], new([]AccessStoreAccess)).(*[]AccessStoreAccess)

	return out0, err

}

// GetUserAccessList is a free data retrieval call binding the contract method 0xbb059b5c.
//
// Solidity: function getUserAccessList(bytes32 accessID) view returns((bytes32,bytes,bytes,uint8)[])
func (_EhrIndexer *EhrIndexerSession) GetUserAccessList(accessID [32]byte) ([]AccessStoreAccess, error) {
	return _EhrIndexer.Contract.GetUserAccessList(&_EhrIndexer.CallOpts, accessID)
}

// GetUserAccessList is a free data retrieval call binding the contract method 0xbb059b5c.
//
// Solidity: function getUserAccessList(bytes32 accessID) view returns((bytes32,bytes,bytes,uint8)[])
func (_EhrIndexer *EhrIndexerCallerSession) GetUserAccessList(accessID [32]byte) ([]AccessStoreAccess, error) {
	return _EhrIndexer.Contract.GetUserAccessList(&_EhrIndexer.CallOpts, accessID)
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

// UserAccess is a free data retrieval call binding the contract method 0xa93f7898.
//
// Solidity: function userAccess(bytes32 userID, uint8 kind, bytes32 idHash) view returns((bytes32,bytes,bytes,uint8))
func (_EhrIndexer *EhrIndexerCaller) UserAccess(opts *bind.CallOpts, userID [32]byte, kind uint8, idHash [32]byte) (AccessStoreAccess, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "userAccess", userID, kind, idHash)

	if err != nil {
		return *new(AccessStoreAccess), err
	}

	out0 := *abi.ConvertType(out[0], new(AccessStoreAccess)).(*AccessStoreAccess)

	return out0, err

}

// UserAccess is a free data retrieval call binding the contract method 0xa93f7898.
//
// Solidity: function userAccess(bytes32 userID, uint8 kind, bytes32 idHash) view returns((bytes32,bytes,bytes,uint8))
func (_EhrIndexer *EhrIndexerSession) UserAccess(userID [32]byte, kind uint8, idHash [32]byte) (AccessStoreAccess, error) {
	return _EhrIndexer.Contract.UserAccess(&_EhrIndexer.CallOpts, userID, kind, idHash)
}

// UserAccess is a free data retrieval call binding the contract method 0xa93f7898.
//
// Solidity: function userAccess(bytes32 userID, uint8 kind, bytes32 idHash) view returns((bytes32,bytes,bytes,uint8))
func (_EhrIndexer *EhrIndexerCallerSession) UserAccess(userID [32]byte, kind uint8, idHash [32]byte) (AccessStoreAccess, error) {
	return _EhrIndexer.Contract.UserAccess(&_EhrIndexer.CallOpts, userID, kind, idHash)
}

// UserGroupGetByID is a free data retrieval call binding the contract method 0xc1dfff99.
//
// Solidity: function userGroupGetByID(bytes32 groupIdHash) view returns(((uint8,bytes)[],(bytes32,bytes)[]))
func (_EhrIndexer *EhrIndexerCaller) UserGroupGetByID(opts *bind.CallOpts, groupIdHash [32]byte) (UsersUserGroup, error) {
	var out []interface{}
	err := _EhrIndexer.contract.Call(opts, &out, "userGroupGetByID", groupIdHash)

	if err != nil {
		return *new(UsersUserGroup), err
	}

	out0 := *abi.ConvertType(out[0], new(UsersUserGroup)).(*UsersUserGroup)

	return out0, err

}

// UserGroupGetByID is a free data retrieval call binding the contract method 0xc1dfff99.
//
// Solidity: function userGroupGetByID(bytes32 groupIdHash) view returns(((uint8,bytes)[],(bytes32,bytes)[]))
func (_EhrIndexer *EhrIndexerSession) UserGroupGetByID(groupIdHash [32]byte) (UsersUserGroup, error) {
	return _EhrIndexer.Contract.UserGroupGetByID(&_EhrIndexer.CallOpts, groupIdHash)
}

// UserGroupGetByID is a free data retrieval call binding the contract method 0xc1dfff99.
//
// Solidity: function userGroupGetByID(bytes32 groupIdHash) view returns(((uint8,bytes)[],(bytes32,bytes)[]))
func (_EhrIndexer *EhrIndexerCallerSession) UserGroupGetByID(groupIdHash [32]byte) (UsersUserGroup, error) {
	return _EhrIndexer.Contract.UserGroupGetByID(&_EhrIndexer.CallOpts, groupIdHash)
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
// Solidity: function docGroupAddDoc(bytes32 groupIdHash, bytes32 CIDHash, bytes CIDEncr, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactor) DocGroupAddDoc(opts *bind.TransactOpts, groupIdHash [32]byte, CIDHash [32]byte, CIDEncr []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "docGroupAddDoc", groupIdHash, CIDHash, CIDEncr, signer, signature)
}

// DocGroupAddDoc is a paid mutator transaction binding the contract method 0x14ce75d6.
//
// Solidity: function docGroupAddDoc(bytes32 groupIdHash, bytes32 CIDHash, bytes CIDEncr, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerSession) DocGroupAddDoc(groupIdHash [32]byte, CIDHash [32]byte, CIDEncr []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.DocGroupAddDoc(&_EhrIndexer.TransactOpts, groupIdHash, CIDHash, CIDEncr, signer, signature)
}

// DocGroupAddDoc is a paid mutator transaction binding the contract method 0x14ce75d6.
//
// Solidity: function docGroupAddDoc(bytes32 groupIdHash, bytes32 CIDHash, bytes CIDEncr, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) DocGroupAddDoc(groupIdHash [32]byte, CIDHash [32]byte, CIDEncr []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.DocGroupAddDoc(&_EhrIndexer.TransactOpts, groupIdHash, CIDHash, CIDEncr, signer, signature)
}

// DocGroupCreate is a paid mutator transaction binding the contract method 0x3a6a5628.
//
// Solidity: function docGroupCreate((bytes32,bytes,bytes,bytes,(uint8,bytes)[],address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerTransactor) DocGroupCreate(opts *bind.TransactOpts, p DocGroupsDocGroupCreateParams) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "docGroupCreate", p)
}

// DocGroupCreate is a paid mutator transaction binding the contract method 0x3a6a5628.
//
// Solidity: function docGroupCreate((bytes32,bytes,bytes,bytes,(uint8,bytes)[],address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerSession) DocGroupCreate(p DocGroupsDocGroupCreateParams) (*types.Transaction, error) {
	return _EhrIndexer.Contract.DocGroupCreate(&_EhrIndexer.TransactOpts, p)
}

// DocGroupCreate is a paid mutator transaction binding the contract method 0x3a6a5628.
//
// Solidity: function docGroupCreate((bytes32,bytes,bytes,bytes,(uint8,bytes)[],address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) DocGroupCreate(p DocGroupsDocGroupCreateParams) (*types.Transaction, error) {
	return _EhrIndexer.Contract.DocGroupCreate(&_EhrIndexer.TransactOpts, p)
}

// GroupAddUser is a paid mutator transaction binding the contract method 0xa544050a.
//
// Solidity: function groupAddUser((bytes32,bytes32,uint8,bytes,bytes,address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerTransactor) GroupAddUser(opts *bind.TransactOpts, p UsersGroupAddUserParams) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "groupAddUser", p)
}

// GroupAddUser is a paid mutator transaction binding the contract method 0xa544050a.
//
// Solidity: function groupAddUser((bytes32,bytes32,uint8,bytes,bytes,address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerSession) GroupAddUser(p UsersGroupAddUserParams) (*types.Transaction, error) {
	return _EhrIndexer.Contract.GroupAddUser(&_EhrIndexer.TransactOpts, p)
}

// GroupAddUser is a paid mutator transaction binding the contract method 0xa544050a.
//
// Solidity: function groupAddUser((bytes32,bytes32,uint8,bytes,bytes,address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) GroupAddUser(p UsersGroupAddUserParams) (*types.Transaction, error) {
	return _EhrIndexer.Contract.GroupAddUser(&_EhrIndexer.TransactOpts, p)
}

// GroupRemoveUser is a paid mutator transaction binding the contract method 0xb94e5cee.
//
// Solidity: function groupRemoveUser(bytes32 groupIDHash, bytes32 userIDHash, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactor) GroupRemoveUser(opts *bind.TransactOpts, groupIDHash [32]byte, userIDHash [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "groupRemoveUser", groupIDHash, userIDHash, signer, signature)
}

// GroupRemoveUser is a paid mutator transaction binding the contract method 0xb94e5cee.
//
// Solidity: function groupRemoveUser(bytes32 groupIDHash, bytes32 userIDHash, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerSession) GroupRemoveUser(groupIDHash [32]byte, userIDHash [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.GroupRemoveUser(&_EhrIndexer.TransactOpts, groupIDHash, userIDHash, signer, signature)
}

// GroupRemoveUser is a paid mutator transaction binding the contract method 0xb94e5cee.
//
// Solidity: function groupRemoveUser(bytes32 groupIDHash, bytes32 userIDHash, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) GroupRemoveUser(groupIDHash [32]byte, userIDHash [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.GroupRemoveUser(&_EhrIndexer.TransactOpts, groupIDHash, userIDHash, signer, signature)
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

// SetDocAccess is a paid mutator transaction binding the contract method 0x3591ac33.
//
// Solidity: function setDocAccess(bytes32 CIDHash, (bytes32,bytes,bytes,uint8) access, address userAddr, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactor) SetDocAccess(opts *bind.TransactOpts, CIDHash [32]byte, access AccessStoreAccess, userAddr common.Address, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setDocAccess", CIDHash, access, userAddr, signer, signature)
}

// SetDocAccess is a paid mutator transaction binding the contract method 0x3591ac33.
//
// Solidity: function setDocAccess(bytes32 CIDHash, (bytes32,bytes,bytes,uint8) access, address userAddr, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerSession) SetDocAccess(CIDHash [32]byte, access AccessStoreAccess, userAddr common.Address, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetDocAccess(&_EhrIndexer.TransactOpts, CIDHash, access, userAddr, signer, signature)
}

// SetDocAccess is a paid mutator transaction binding the contract method 0x3591ac33.
//
// Solidity: function setDocAccess(bytes32 CIDHash, (bytes32,bytes,bytes,uint8) access, address userAddr, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) SetDocAccess(CIDHash [32]byte, access AccessStoreAccess, userAddr common.Address, signer common.Address, signature []byte) (*types.Transaction, error) {
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
// Solidity: function setEhrUser(bytes32 userId, bytes32 ehrId, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactor) SetEhrUser(opts *bind.TransactOpts, userId [32]byte, ehrId [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "setEhrUser", userId, ehrId, signer, signature)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0x3f157693.
//
// Solidity: function setEhrUser(bytes32 userId, bytes32 ehrId, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerSession) SetEhrUser(userId [32]byte, ehrId [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrUser(&_EhrIndexer.TransactOpts, userId, ehrId, signer, signature)
}

// SetEhrUser is a paid mutator transaction binding the contract method 0x3f157693.
//
// Solidity: function setEhrUser(bytes32 userId, bytes32 ehrId, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) SetEhrUser(userId [32]byte, ehrId [32]byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.SetEhrUser(&_EhrIndexer.TransactOpts, userId, ehrId, signer, signature)
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

// UserGroupCreate is a paid mutator transaction binding the contract method 0x5e6a6500.
//
// Solidity: function userGroupCreate((bytes32,(uint8,bytes)[],address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerTransactor) UserGroupCreate(opts *bind.TransactOpts, p UsersUserGroupCreateParams) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "userGroupCreate", p)
}

// UserGroupCreate is a paid mutator transaction binding the contract method 0x5e6a6500.
//
// Solidity: function userGroupCreate((bytes32,(uint8,bytes)[],address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerSession) UserGroupCreate(p UsersUserGroupCreateParams) (*types.Transaction, error) {
	return _EhrIndexer.Contract.UserGroupCreate(&_EhrIndexer.TransactOpts, p)
}

// UserGroupCreate is a paid mutator transaction binding the contract method 0x5e6a6500.
//
// Solidity: function userGroupCreate((bytes32,(uint8,bytes)[],address,bytes) p) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) UserGroupCreate(p UsersUserGroupCreateParams) (*types.Transaction, error) {
	return _EhrIndexer.Contract.UserGroupCreate(&_EhrIndexer.TransactOpts, p)
}

// UserNew is a paid mutator transaction binding the contract method 0x60cf0a8b.
//
// Solidity: function userNew(address userAddr, bytes32 id, bytes32 systemID, uint8 role, bytes pwdHash, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactor) UserNew(opts *bind.TransactOpts, userAddr common.Address, id [32]byte, systemID [32]byte, role uint8, pwdHash []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.contract.Transact(opts, "userNew", userAddr, id, systemID, role, pwdHash, signer, signature)
}

// UserNew is a paid mutator transaction binding the contract method 0x60cf0a8b.
//
// Solidity: function userNew(address userAddr, bytes32 id, bytes32 systemID, uint8 role, bytes pwdHash, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerSession) UserNew(userAddr common.Address, id [32]byte, systemID [32]byte, role uint8, pwdHash []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.UserNew(&_EhrIndexer.TransactOpts, userAddr, id, systemID, role, pwdHash, signer, signature)
}

// UserNew is a paid mutator transaction binding the contract method 0x60cf0a8b.
//
// Solidity: function userNew(address userAddr, bytes32 id, bytes32 systemID, uint8 role, bytes pwdHash, address signer, bytes signature) returns()
func (_EhrIndexer *EhrIndexerTransactorSession) UserNew(userAddr common.Address, id [32]byte, systemID [32]byte, role uint8, pwdHash []byte, signer common.Address, signature []byte) (*types.Transaction, error) {
	return _EhrIndexer.Contract.UserNew(&_EhrIndexer.TransactOpts, userAddr, id, systemID, role, pwdHash, signer, signature)
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
