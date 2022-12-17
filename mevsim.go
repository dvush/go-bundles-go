// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package main

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

// MevSimMetaData contains all meta data concerning the MevSim contract.
var MevSimMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[],\"name\":\"BlockMismatch\",\"type\":\"error\"},{\"inputs\":[],\"name\":\"SlotValueMismatch\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"slot\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"value\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"target_block\",\"type\":\"uint256\"}],\"name\":\"auction\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"slot\",\"type\":\"uint256\"}],\"name\":\"getSlot\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// MevSimABI is the input ABI used to generate the binding from.
// Deprecated: Use MevSimMetaData.ABI instead.
var MevSimABI = MevSimMetaData.ABI

// MevSim is an auto generated Go binding around an Ethereum contract.
type MevSim struct {
	MevSimCaller     // Read-only binding to the contract
	MevSimTransactor // Write-only binding to the contract
	MevSimFilterer   // Log filterer for contract events
}

// MevSimCaller is an auto generated read-only Go binding around an Ethereum contract.
type MevSimCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MevSimTransactor is an auto generated write-only Go binding around an Ethereum contract.
type MevSimTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MevSimFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type MevSimFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// MevSimSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type MevSimSession struct {
	Contract     *MevSim           // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MevSimCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type MevSimCallerSession struct {
	Contract *MevSimCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// MevSimTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type MevSimTransactorSession struct {
	Contract     *MevSimTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// MevSimRaw is an auto generated low-level Go binding around an Ethereum contract.
type MevSimRaw struct {
	Contract *MevSim // Generic contract binding to access the raw methods on
}

// MevSimCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type MevSimCallerRaw struct {
	Contract *MevSimCaller // Generic read-only contract binding to access the raw methods on
}

// MevSimTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type MevSimTransactorRaw struct {
	Contract *MevSimTransactor // Generic write-only contract binding to access the raw methods on
}

// NewMevSim creates a new instance of MevSim, bound to a specific deployed contract.
func NewMevSim(address common.Address, backend bind.ContractBackend) (*MevSim, error) {
	contract, err := bindMevSim(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &MevSim{MevSimCaller: MevSimCaller{contract: contract}, MevSimTransactor: MevSimTransactor{contract: contract}, MevSimFilterer: MevSimFilterer{contract: contract}}, nil
}

// NewMevSimCaller creates a new read-only instance of MevSim, bound to a specific deployed contract.
func NewMevSimCaller(address common.Address, caller bind.ContractCaller) (*MevSimCaller, error) {
	contract, err := bindMevSim(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &MevSimCaller{contract: contract}, nil
}

// NewMevSimTransactor creates a new write-only instance of MevSim, bound to a specific deployed contract.
func NewMevSimTransactor(address common.Address, transactor bind.ContractTransactor) (*MevSimTransactor, error) {
	contract, err := bindMevSim(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &MevSimTransactor{contract: contract}, nil
}

// NewMevSimFilterer creates a new log filterer instance of MevSim, bound to a specific deployed contract.
func NewMevSimFilterer(address common.Address, filterer bind.ContractFilterer) (*MevSimFilterer, error) {
	contract, err := bindMevSim(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &MevSimFilterer{contract: contract}, nil
}

// bindMevSim binds a generic wrapper to an already deployed contract.
func bindMevSim(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(MevSimABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MevSim *MevSimRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MevSim.Contract.MevSimCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MevSim *MevSimRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MevSim.Contract.MevSimTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MevSim *MevSimRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MevSim.Contract.MevSimTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_MevSim *MevSimCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _MevSim.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_MevSim *MevSimTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _MevSim.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_MevSim *MevSimTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _MevSim.Contract.contract.Transact(opts, method, params...)
}

// GetSlot is a free data retrieval call binding the contract method 0x7eba7ba6.
//
// Solidity: function getSlot(uint256 slot) view returns(uint256)
func (_MevSim *MevSimCaller) GetSlot(opts *bind.CallOpts, slot *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _MevSim.contract.Call(opts, &out, "getSlot", slot)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetSlot is a free data retrieval call binding the contract method 0x7eba7ba6.
//
// Solidity: function getSlot(uint256 slot) view returns(uint256)
func (_MevSim *MevSimSession) GetSlot(slot *big.Int) (*big.Int, error) {
	return _MevSim.Contract.GetSlot(&_MevSim.CallOpts, slot)
}

// GetSlot is a free data retrieval call binding the contract method 0x7eba7ba6.
//
// Solidity: function getSlot(uint256 slot) view returns(uint256)
func (_MevSim *MevSimCallerSession) GetSlot(slot *big.Int) (*big.Int, error) {
	return _MevSim.Contract.GetSlot(&_MevSim.CallOpts, slot)
}

// Auction is a paid mutator transaction binding the contract method 0xb73e7399.
//
// Solidity: function auction(uint256 slot, uint256 value, uint256 target_block) payable returns()
func (_MevSim *MevSimTransactor) Auction(opts *bind.TransactOpts, slot *big.Int, value *big.Int, target_block *big.Int) (*types.Transaction, error) {
	return _MevSim.contract.Transact(opts, "auction", slot, value, target_block)
}

// Auction is a paid mutator transaction binding the contract method 0xb73e7399.
//
// Solidity: function auction(uint256 slot, uint256 value, uint256 target_block) payable returns()
func (_MevSim *MevSimSession) Auction(slot *big.Int, value *big.Int, target_block *big.Int) (*types.Transaction, error) {
	return _MevSim.Contract.Auction(&_MevSim.TransactOpts, slot, value, target_block)
}

// Auction is a paid mutator transaction binding the contract method 0xb73e7399.
//
// Solidity: function auction(uint256 slot, uint256 value, uint256 target_block) payable returns()
func (_MevSim *MevSimTransactorSession) Auction(slot *big.Int, value *big.Int, target_block *big.Int) (*types.Transaction, error) {
	return _MevSim.Contract.Auction(&_MevSim.TransactOpts, slot, value, target_block)
}
