// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package bindings

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
	_ = abi.ConvertType
)

// L2StandardBridgeBotMetaData contains all meta data concerning the L2StandardBridgeBot contract.
var L2StandardBridgeBotMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"addresspayable\",\"name\":\"_owner\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_delegationFee\",\"type\":\"uint256\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"owner\",\"type\":\"address\"}],\"name\":\"OwnableInvalidOwner\",\"type\":\"error\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"account\",\"type\":\"address\"}],\"name\":\"OwnableUnauthorizedAccount\",\"type\":\"error\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"previousOwner\",\"type\":\"address\"},{\"indexed\":true,\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"OwnershipTransferred\",\"type\":\"event\"},{\"anonymous\":false,\"inputs\":[{\"indexed\":true,\"internalType\":\"address\",\"name\":\"from\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"l2Token\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"address\",\"name\":\"to\",\"type\":\"address\"},{\"indexed\":false,\"internalType\":\"uint256\",\"name\":\"amount\",\"type\":\"uint256\"},{\"indexed\":false,\"internalType\":\"uint32\",\"name\":\"minGasLimit\",\"type\":\"uint32\"},{\"indexed\":false,\"internalType\":\"bytes\",\"name\":\"extraData\",\"type\":\"bytes\"}],\"name\":\"WithdrawTo\",\"type\":\"event\"},{\"stateMutability\":\"payable\",\"type\":\"fallback\"},{\"inputs\":[],\"name\":\"L2_STANDARD_BRIDGE\",\"outputs\":[{\"internalType\":\"contractIL2StandardBridge\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"L2_STANDARD_BRIDGE_ADDRESS\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"delegationFee\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"owner\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"renounceOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint256\",\"name\":\"_delegationFee\",\"type\":\"uint256\"}],\"name\":\"setDelegationFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"newOwner\",\"type\":\"address\"}],\"name\":\"transferOwnership\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_l2Token\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"_minGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"withdraw\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"}],\"name\":\"withdrawFee\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_recipient\",\"type\":\"address\"},{\"internalType\":\"uint32\",\"name\":\"_minGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"withdrawFeeToL1\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"_l2Token\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"_to\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"},{\"internalType\":\"uint32\",\"name\":\"_minGasLimit\",\"type\":\"uint32\"},{\"internalType\":\"bytes\",\"name\":\"_extraData\",\"type\":\"bytes\"}],\"name\":\"withdrawTo\",\"outputs\":[],\"stateMutability\":\"payable\",\"type\":\"function\"},{\"stateMutability\":\"payable\",\"type\":\"receive\"}]",
}

// L2StandardBridgeBotABI is the input ABI used to generate the binding from.
// Deprecated: Use L2StandardBridgeBotMetaData.ABI instead.
var L2StandardBridgeBotABI = L2StandardBridgeBotMetaData.ABI

// L2StandardBridgeBot is an auto generated Go binding around an Ethereum contract.
type L2StandardBridgeBot struct {
	L2StandardBridgeBotCaller     // Read-only binding to the contract
	L2StandardBridgeBotTransactor // Write-only binding to the contract
	L2StandardBridgeBotFilterer   // Log filterer for contract events
}

// L2StandardBridgeBotCaller is an auto generated read-only Go binding around an Ethereum contract.
type L2StandardBridgeBotCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2StandardBridgeBotTransactor is an auto generated write-only Go binding around an Ethereum contract.
type L2StandardBridgeBotTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2StandardBridgeBotFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type L2StandardBridgeBotFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// L2StandardBridgeBotSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type L2StandardBridgeBotSession struct {
	Contract     *L2StandardBridgeBot // Generic contract binding to set the session for
	CallOpts     bind.CallOpts        // Call options to use throughout this session
	TransactOpts bind.TransactOpts    // Transaction auth options to use throughout this session
}

// L2StandardBridgeBotCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type L2StandardBridgeBotCallerSession struct {
	Contract *L2StandardBridgeBotCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts              // Call options to use throughout this session
}

// L2StandardBridgeBotTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type L2StandardBridgeBotTransactorSession struct {
	Contract     *L2StandardBridgeBotTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts              // Transaction auth options to use throughout this session
}

// L2StandardBridgeBotRaw is an auto generated low-level Go binding around an Ethereum contract.
type L2StandardBridgeBotRaw struct {
	Contract *L2StandardBridgeBot // Generic contract binding to access the raw methods on
}

// L2StandardBridgeBotCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type L2StandardBridgeBotCallerRaw struct {
	Contract *L2StandardBridgeBotCaller // Generic read-only contract binding to access the raw methods on
}

// L2StandardBridgeBotTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type L2StandardBridgeBotTransactorRaw struct {
	Contract *L2StandardBridgeBotTransactor // Generic write-only contract binding to access the raw methods on
}

// NewL2StandardBridgeBot creates a new instance of L2StandardBridgeBot, bound to a specific deployed contract.
func NewL2StandardBridgeBot(address common.Address, backend bind.ContractBackend) (*L2StandardBridgeBot, error) {
	contract, err := bindL2StandardBridgeBot(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &L2StandardBridgeBot{L2StandardBridgeBotCaller: L2StandardBridgeBotCaller{contract: contract}, L2StandardBridgeBotTransactor: L2StandardBridgeBotTransactor{contract: contract}, L2StandardBridgeBotFilterer: L2StandardBridgeBotFilterer{contract: contract}}, nil
}

// NewL2StandardBridgeBotCaller creates a new read-only instance of L2StandardBridgeBot, bound to a specific deployed contract.
func NewL2StandardBridgeBotCaller(address common.Address, caller bind.ContractCaller) (*L2StandardBridgeBotCaller, error) {
	contract, err := bindL2StandardBridgeBot(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &L2StandardBridgeBotCaller{contract: contract}, nil
}

// NewL2StandardBridgeBotTransactor creates a new write-only instance of L2StandardBridgeBot, bound to a specific deployed contract.
func NewL2StandardBridgeBotTransactor(address common.Address, transactor bind.ContractTransactor) (*L2StandardBridgeBotTransactor, error) {
	contract, err := bindL2StandardBridgeBot(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &L2StandardBridgeBotTransactor{contract: contract}, nil
}

// NewL2StandardBridgeBotFilterer creates a new log filterer instance of L2StandardBridgeBot, bound to a specific deployed contract.
func NewL2StandardBridgeBotFilterer(address common.Address, filterer bind.ContractFilterer) (*L2StandardBridgeBotFilterer, error) {
	contract, err := bindL2StandardBridgeBot(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &L2StandardBridgeBotFilterer{contract: contract}, nil
}

// bindL2StandardBridgeBot binds a generic wrapper to an already deployed contract.
func bindL2StandardBridgeBot(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := L2StandardBridgeBotMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L2StandardBridgeBot *L2StandardBridgeBotRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2StandardBridgeBot.Contract.L2StandardBridgeBotCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L2StandardBridgeBot *L2StandardBridgeBotRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.L2StandardBridgeBotTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L2StandardBridgeBot *L2StandardBridgeBotRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.L2StandardBridgeBotTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_L2StandardBridgeBot *L2StandardBridgeBotCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _L2StandardBridgeBot.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.contract.Transact(opts, method, params...)
}

// L2STANDARDBRIDGE is a free data retrieval call binding the contract method 0x21d12763.
//
// Solidity: function L2_STANDARD_BRIDGE() view returns(address)
func (_L2StandardBridgeBot *L2StandardBridgeBotCaller) L2STANDARDBRIDGE(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2StandardBridgeBot.contract.Call(opts, &out, "L2_STANDARD_BRIDGE")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2STANDARDBRIDGE is a free data retrieval call binding the contract method 0x21d12763.
//
// Solidity: function L2_STANDARD_BRIDGE() view returns(address)
func (_L2StandardBridgeBot *L2StandardBridgeBotSession) L2STANDARDBRIDGE() (common.Address, error) {
	return _L2StandardBridgeBot.Contract.L2STANDARDBRIDGE(&_L2StandardBridgeBot.CallOpts)
}

// L2STANDARDBRIDGE is a free data retrieval call binding the contract method 0x21d12763.
//
// Solidity: function L2_STANDARD_BRIDGE() view returns(address)
func (_L2StandardBridgeBot *L2StandardBridgeBotCallerSession) L2STANDARDBRIDGE() (common.Address, error) {
	return _L2StandardBridgeBot.Contract.L2STANDARDBRIDGE(&_L2StandardBridgeBot.CallOpts)
}

// L2STANDARDBRIDGEADDRESS is a free data retrieval call binding the contract method 0x2cb7cb06.
//
// Solidity: function L2_STANDARD_BRIDGE_ADDRESS() view returns(address)
func (_L2StandardBridgeBot *L2StandardBridgeBotCaller) L2STANDARDBRIDGEADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2StandardBridgeBot.contract.Call(opts, &out, "L2_STANDARD_BRIDGE_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// L2STANDARDBRIDGEADDRESS is a free data retrieval call binding the contract method 0x2cb7cb06.
//
// Solidity: function L2_STANDARD_BRIDGE_ADDRESS() view returns(address)
func (_L2StandardBridgeBot *L2StandardBridgeBotSession) L2STANDARDBRIDGEADDRESS() (common.Address, error) {
	return _L2StandardBridgeBot.Contract.L2STANDARDBRIDGEADDRESS(&_L2StandardBridgeBot.CallOpts)
}

// L2STANDARDBRIDGEADDRESS is a free data retrieval call binding the contract method 0x2cb7cb06.
//
// Solidity: function L2_STANDARD_BRIDGE_ADDRESS() view returns(address)
func (_L2StandardBridgeBot *L2StandardBridgeBotCallerSession) L2STANDARDBRIDGEADDRESS() (common.Address, error) {
	return _L2StandardBridgeBot.Contract.L2STANDARDBRIDGEADDRESS(&_L2StandardBridgeBot.CallOpts)
}

// DelegationFee is a free data retrieval call binding the contract method 0xc5f0a58f.
//
// Solidity: function delegationFee() view returns(uint256)
func (_L2StandardBridgeBot *L2StandardBridgeBotCaller) DelegationFee(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _L2StandardBridgeBot.contract.Call(opts, &out, "delegationFee")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// DelegationFee is a free data retrieval call binding the contract method 0xc5f0a58f.
//
// Solidity: function delegationFee() view returns(uint256)
func (_L2StandardBridgeBot *L2StandardBridgeBotSession) DelegationFee() (*big.Int, error) {
	return _L2StandardBridgeBot.Contract.DelegationFee(&_L2StandardBridgeBot.CallOpts)
}

// DelegationFee is a free data retrieval call binding the contract method 0xc5f0a58f.
//
// Solidity: function delegationFee() view returns(uint256)
func (_L2StandardBridgeBot *L2StandardBridgeBotCallerSession) DelegationFee() (*big.Int, error) {
	return _L2StandardBridgeBot.Contract.DelegationFee(&_L2StandardBridgeBot.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2StandardBridgeBot *L2StandardBridgeBotCaller) Owner(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _L2StandardBridgeBot.contract.Call(opts, &out, "owner")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2StandardBridgeBot *L2StandardBridgeBotSession) Owner() (common.Address, error) {
	return _L2StandardBridgeBot.Contract.Owner(&_L2StandardBridgeBot.CallOpts)
}

// Owner is a free data retrieval call binding the contract method 0x8da5cb5b.
//
// Solidity: function owner() view returns(address)
func (_L2StandardBridgeBot *L2StandardBridgeBotCallerSession) Owner() (common.Address, error) {
	return _L2StandardBridgeBot.Contract.Owner(&_L2StandardBridgeBot.CallOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactor) RenounceOwnership(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2StandardBridgeBot.contract.Transact(opts, "renounceOwnership")
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotSession) RenounceOwnership() (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.RenounceOwnership(&_L2StandardBridgeBot.TransactOpts)
}

// RenounceOwnership is a paid mutator transaction binding the contract method 0x715018a6.
//
// Solidity: function renounceOwnership() returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactorSession) RenounceOwnership() (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.RenounceOwnership(&_L2StandardBridgeBot.TransactOpts)
}

// SetDelegationFee is a paid mutator transaction binding the contract method 0x55bfc81c.
//
// Solidity: function setDelegationFee(uint256 _delegationFee) returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactor) SetDelegationFee(opts *bind.TransactOpts, _delegationFee *big.Int) (*types.Transaction, error) {
	return _L2StandardBridgeBot.contract.Transact(opts, "setDelegationFee", _delegationFee)
}

// SetDelegationFee is a paid mutator transaction binding the contract method 0x55bfc81c.
//
// Solidity: function setDelegationFee(uint256 _delegationFee) returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotSession) SetDelegationFee(_delegationFee *big.Int) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.SetDelegationFee(&_L2StandardBridgeBot.TransactOpts, _delegationFee)
}

// SetDelegationFee is a paid mutator transaction binding the contract method 0x55bfc81c.
//
// Solidity: function setDelegationFee(uint256 _delegationFee) returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactorSession) SetDelegationFee(_delegationFee *big.Int) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.SetDelegationFee(&_L2StandardBridgeBot.TransactOpts, _delegationFee)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactor) TransferOwnership(opts *bind.TransactOpts, newOwner common.Address) (*types.Transaction, error) {
	return _L2StandardBridgeBot.contract.Transact(opts, "transferOwnership", newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.TransferOwnership(&_L2StandardBridgeBot.TransactOpts, newOwner)
}

// TransferOwnership is a paid mutator transaction binding the contract method 0xf2fde38b.
//
// Solidity: function transferOwnership(address newOwner) returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactorSession) TransferOwnership(newOwner common.Address) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.TransferOwnership(&_L2StandardBridgeBot.TransactOpts, newOwner)
}

// Withdraw is a paid mutator transaction binding the contract method 0x32b7006d.
//
// Solidity: function withdraw(address _l2Token, uint256 _amount, uint32 _minGasLimit, bytes _extraData) payable returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactor) Withdraw(opts *bind.TransactOpts, _l2Token common.Address, _amount *big.Int, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L2StandardBridgeBot.contract.Transact(opts, "withdraw", _l2Token, _amount, _minGasLimit, _extraData)
}

// Withdraw is a paid mutator transaction binding the contract method 0x32b7006d.
//
// Solidity: function withdraw(address _l2Token, uint256 _amount, uint32 _minGasLimit, bytes _extraData) payable returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotSession) Withdraw(_l2Token common.Address, _amount *big.Int, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.Withdraw(&_L2StandardBridgeBot.TransactOpts, _l2Token, _amount, _minGasLimit, _extraData)
}

// Withdraw is a paid mutator transaction binding the contract method 0x32b7006d.
//
// Solidity: function withdraw(address _l2Token, uint256 _amount, uint32 _minGasLimit, bytes _extraData) payable returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactorSession) Withdraw(_l2Token common.Address, _amount *big.Int, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.Withdraw(&_L2StandardBridgeBot.TransactOpts, _l2Token, _amount, _minGasLimit, _extraData)
}

// WithdrawFee is a paid mutator transaction binding the contract method 0x1ac3ddeb.
//
// Solidity: function withdrawFee(address _recipient) returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactor) WithdrawFee(opts *bind.TransactOpts, _recipient common.Address) (*types.Transaction, error) {
	return _L2StandardBridgeBot.contract.Transact(opts, "withdrawFee", _recipient)
}

// WithdrawFee is a paid mutator transaction binding the contract method 0x1ac3ddeb.
//
// Solidity: function withdrawFee(address _recipient) returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotSession) WithdrawFee(_recipient common.Address) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.WithdrawFee(&_L2StandardBridgeBot.TransactOpts, _recipient)
}

// WithdrawFee is a paid mutator transaction binding the contract method 0x1ac3ddeb.
//
// Solidity: function withdrawFee(address _recipient) returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactorSession) WithdrawFee(_recipient common.Address) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.WithdrawFee(&_L2StandardBridgeBot.TransactOpts, _recipient)
}

// WithdrawFeeToL1 is a paid mutator transaction binding the contract method 0x244cafe0.
//
// Solidity: function withdrawFeeToL1(address _recipient, uint32 _minGasLimit, bytes _extraData) returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactor) WithdrawFeeToL1(opts *bind.TransactOpts, _recipient common.Address, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L2StandardBridgeBot.contract.Transact(opts, "withdrawFeeToL1", _recipient, _minGasLimit, _extraData)
}

// WithdrawFeeToL1 is a paid mutator transaction binding the contract method 0x244cafe0.
//
// Solidity: function withdrawFeeToL1(address _recipient, uint32 _minGasLimit, bytes _extraData) returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotSession) WithdrawFeeToL1(_recipient common.Address, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.WithdrawFeeToL1(&_L2StandardBridgeBot.TransactOpts, _recipient, _minGasLimit, _extraData)
}

// WithdrawFeeToL1 is a paid mutator transaction binding the contract method 0x244cafe0.
//
// Solidity: function withdrawFeeToL1(address _recipient, uint32 _minGasLimit, bytes _extraData) returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactorSession) WithdrawFeeToL1(_recipient common.Address, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.WithdrawFeeToL1(&_L2StandardBridgeBot.TransactOpts, _recipient, _minGasLimit, _extraData)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0xa3a79548.
//
// Solidity: function withdrawTo(address _l2Token, address _to, uint256 _amount, uint32 _minGasLimit, bytes _extraData) payable returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactor) WithdrawTo(opts *bind.TransactOpts, _l2Token common.Address, _to common.Address, _amount *big.Int, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L2StandardBridgeBot.contract.Transact(opts, "withdrawTo", _l2Token, _to, _amount, _minGasLimit, _extraData)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0xa3a79548.
//
// Solidity: function withdrawTo(address _l2Token, address _to, uint256 _amount, uint32 _minGasLimit, bytes _extraData) payable returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotSession) WithdrawTo(_l2Token common.Address, _to common.Address, _amount *big.Int, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.WithdrawTo(&_L2StandardBridgeBot.TransactOpts, _l2Token, _to, _amount, _minGasLimit, _extraData)
}

// WithdrawTo is a paid mutator transaction binding the contract method 0xa3a79548.
//
// Solidity: function withdrawTo(address _l2Token, address _to, uint256 _amount, uint32 _minGasLimit, bytes _extraData) payable returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactorSession) WithdrawTo(_l2Token common.Address, _to common.Address, _amount *big.Int, _minGasLimit uint32, _extraData []byte) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.WithdrawTo(&_L2StandardBridgeBot.TransactOpts, _l2Token, _to, _amount, _minGasLimit, _extraData)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactor) Fallback(opts *bind.TransactOpts, calldata []byte) (*types.Transaction, error) {
	return _L2StandardBridgeBot.contract.RawTransact(opts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.Fallback(&_L2StandardBridgeBot.TransactOpts, calldata)
}

// Fallback is a paid mutator transaction binding the contract fallback function.
//
// Solidity: fallback() payable returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactorSession) Fallback(calldata []byte) (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.Fallback(&_L2StandardBridgeBot.TransactOpts, calldata)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactor) Receive(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _L2StandardBridgeBot.contract.RawTransact(opts, nil) // calldata is disallowed for receive function
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotSession) Receive() (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.Receive(&_L2StandardBridgeBot.TransactOpts)
}

// Receive is a paid mutator transaction binding the contract receive function.
//
// Solidity: receive() payable returns()
func (_L2StandardBridgeBot *L2StandardBridgeBotTransactorSession) Receive() (*types.Transaction, error) {
	return _L2StandardBridgeBot.Contract.Receive(&_L2StandardBridgeBot.TransactOpts)
}

// L2StandardBridgeBotOwnershipTransferredIterator is returned from FilterOwnershipTransferred and is used to iterate over the raw logs and unpacked data for OwnershipTransferred events raised by the L2StandardBridgeBot contract.
type L2StandardBridgeBotOwnershipTransferredIterator struct {
	Event *L2StandardBridgeBotOwnershipTransferred // Event containing the contract specifics and raw log

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
func (it *L2StandardBridgeBotOwnershipTransferredIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2StandardBridgeBotOwnershipTransferred)
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
		it.Event = new(L2StandardBridgeBotOwnershipTransferred)
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
func (it *L2StandardBridgeBotOwnershipTransferredIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2StandardBridgeBotOwnershipTransferredIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2StandardBridgeBotOwnershipTransferred represents a OwnershipTransferred event raised by the L2StandardBridgeBot contract.
type L2StandardBridgeBotOwnershipTransferred struct {
	PreviousOwner common.Address
	NewOwner      common.Address
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterOwnershipTransferred is a free log retrieval operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L2StandardBridgeBot *L2StandardBridgeBotFilterer) FilterOwnershipTransferred(opts *bind.FilterOpts, previousOwner []common.Address, newOwner []common.Address) (*L2StandardBridgeBotOwnershipTransferredIterator, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L2StandardBridgeBot.contract.FilterLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return &L2StandardBridgeBotOwnershipTransferredIterator{contract: _L2StandardBridgeBot.contract, event: "OwnershipTransferred", logs: logs, sub: sub}, nil
}

// WatchOwnershipTransferred is a free log subscription operation binding the contract event 0x8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e0.
//
// Solidity: event OwnershipTransferred(address indexed previousOwner, address indexed newOwner)
func (_L2StandardBridgeBot *L2StandardBridgeBotFilterer) WatchOwnershipTransferred(opts *bind.WatchOpts, sink chan<- *L2StandardBridgeBotOwnershipTransferred, previousOwner []common.Address, newOwner []common.Address) (event.Subscription, error) {

	var previousOwnerRule []interface{}
	for _, previousOwnerItem := range previousOwner {
		previousOwnerRule = append(previousOwnerRule, previousOwnerItem)
	}
	var newOwnerRule []interface{}
	for _, newOwnerItem := range newOwner {
		newOwnerRule = append(newOwnerRule, newOwnerItem)
	}

	logs, sub, err := _L2StandardBridgeBot.contract.WatchLogs(opts, "OwnershipTransferred", previousOwnerRule, newOwnerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2StandardBridgeBotOwnershipTransferred)
				if err := _L2StandardBridgeBot.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
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
func (_L2StandardBridgeBot *L2StandardBridgeBotFilterer) ParseOwnershipTransferred(log types.Log) (*L2StandardBridgeBotOwnershipTransferred, error) {
	event := new(L2StandardBridgeBotOwnershipTransferred)
	if err := _L2StandardBridgeBot.contract.UnpackLog(event, "OwnershipTransferred", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// L2StandardBridgeBotWithdrawToIterator is returned from FilterWithdrawTo and is used to iterate over the raw logs and unpacked data for WithdrawTo events raised by the L2StandardBridgeBot contract.
type L2StandardBridgeBotWithdrawToIterator struct {
	Event *L2StandardBridgeBotWithdrawTo // Event containing the contract specifics and raw log

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
func (it *L2StandardBridgeBotWithdrawToIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(L2StandardBridgeBotWithdrawTo)
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
		it.Event = new(L2StandardBridgeBotWithdrawTo)
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
func (it *L2StandardBridgeBotWithdrawToIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *L2StandardBridgeBotWithdrawToIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// L2StandardBridgeBotWithdrawTo represents a WithdrawTo event raised by the L2StandardBridgeBot contract.
type L2StandardBridgeBotWithdrawTo struct {
	From        common.Address
	L2Token     common.Address
	To          common.Address
	Amount      *big.Int
	MinGasLimit uint32
	ExtraData   []byte
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterWithdrawTo is a free log retrieval operation binding the contract event 0x56f66275d9ebc94b7d6895aa0d96a3783550d0183ba106408d387d19f2e877f1.
//
// Solidity: event WithdrawTo(address indexed from, address l2Token, address to, uint256 amount, uint32 minGasLimit, bytes extraData)
func (_L2StandardBridgeBot *L2StandardBridgeBotFilterer) FilterWithdrawTo(opts *bind.FilterOpts, from []common.Address) (*L2StandardBridgeBotWithdrawToIterator, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _L2StandardBridgeBot.contract.FilterLogs(opts, "WithdrawTo", fromRule)
	if err != nil {
		return nil, err
	}
	return &L2StandardBridgeBotWithdrawToIterator{contract: _L2StandardBridgeBot.contract, event: "WithdrawTo", logs: logs, sub: sub}, nil
}

// WatchWithdrawTo is a free log subscription operation binding the contract event 0x56f66275d9ebc94b7d6895aa0d96a3783550d0183ba106408d387d19f2e877f1.
//
// Solidity: event WithdrawTo(address indexed from, address l2Token, address to, uint256 amount, uint32 minGasLimit, bytes extraData)
func (_L2StandardBridgeBot *L2StandardBridgeBotFilterer) WatchWithdrawTo(opts *bind.WatchOpts, sink chan<- *L2StandardBridgeBotWithdrawTo, from []common.Address) (event.Subscription, error) {

	var fromRule []interface{}
	for _, fromItem := range from {
		fromRule = append(fromRule, fromItem)
	}

	logs, sub, err := _L2StandardBridgeBot.contract.WatchLogs(opts, "WithdrawTo", fromRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(L2StandardBridgeBotWithdrawTo)
				if err := _L2StandardBridgeBot.contract.UnpackLog(event, "WithdrawTo", log); err != nil {
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

// ParseWithdrawTo is a log parse operation binding the contract event 0x56f66275d9ebc94b7d6895aa0d96a3783550d0183ba106408d387d19f2e877f1.
//
// Solidity: event WithdrawTo(address indexed from, address l2Token, address to, uint256 amount, uint32 minGasLimit, bytes extraData)
func (_L2StandardBridgeBot *L2StandardBridgeBotFilterer) ParseWithdrawTo(log types.Log) (*L2StandardBridgeBotWithdrawTo, error) {
	event := new(L2StandardBridgeBotWithdrawTo)
	if err := _L2StandardBridgeBot.contract.UnpackLog(event, "WithdrawTo", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
