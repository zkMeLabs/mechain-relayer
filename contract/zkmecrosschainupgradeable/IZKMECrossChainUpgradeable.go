// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package zkmecrosschainupgradeable

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

// KYCDataLibEventData is an auto generated low-level Go binding around an user-defined struct.
type KYCDataLibEventData struct {
	SrcChainId  uint32
	DestChainId uint32
	Sequence    *big.Int
	ChannelId   *big.Int
	Payload     []byte
}

// IZKMECrossChainUpgradeableMetaData contains all meta data concerning the IZKMECrossChainUpgradeable contract.
var IZKMECrossChainUpgradeableMetaData = &bind.MetaData{
	ABI: "[{\"anonymous\":false,\"inputs\":[{\"components\":[{\"internalType\":\"uint32\",\"name\":\"srcChainId\",\"type\":\"uint32\"},{\"internalType\":\"uint32\",\"name\":\"destChainId\",\"type\":\"uint32\"},{\"internalType\":\"uint256\",\"name\":\"sequence\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"channelId\",\"type\":\"uint256\"},{\"internalType\":\"bytes\",\"name\":\"payload\",\"type\":\"bytes\"}],\"indexed\":true,\"internalType\":\"structKYCDataLib.EventData\",\"name\":\"eventData\",\"type\":\"tuple\"}],\"name\":\"ZkmeSBTCrossChainPackage\",\"type\":\"event\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"chainId\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"tokenId\",\"type\":\"uint256\"}],\"name\":\"ackMinted\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"chainId\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"srcUser\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"destUser\",\"type\":\"address\"}],\"name\":\"forward\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"uint32\",\"name\":\"chainId\",\"type\":\"uint32\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getCrossChainStatus\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"\",\"type\":\"uint256\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]",
}

// IZKMECrossChainUpgradeableABI is the input ABI used to generate the binding from.
// Deprecated: Use IZKMECrossChainUpgradeableMetaData.ABI instead.
var IZKMECrossChainUpgradeableABI = IZKMECrossChainUpgradeableMetaData.ABI

// IZKMECrossChainUpgradeable is an auto generated Go binding around an Ethereum contract.
type IZKMECrossChainUpgradeable struct {
	IZKMECrossChainUpgradeableCaller     // Read-only binding to the contract
	IZKMECrossChainUpgradeableTransactor // Write-only binding to the contract
	IZKMECrossChainUpgradeableFilterer   // Log filterer for contract events
}

// IZKMECrossChainUpgradeableCaller is an auto generated read-only Go binding around an Ethereum contract.
type IZKMECrossChainUpgradeableCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IZKMECrossChainUpgradeableTransactor is an auto generated write-only Go binding around an Ethereum contract.
type IZKMECrossChainUpgradeableTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IZKMECrossChainUpgradeableFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type IZKMECrossChainUpgradeableFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// IZKMECrossChainUpgradeableSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type IZKMECrossChainUpgradeableSession struct {
	Contract     *IZKMECrossChainUpgradeable // Generic contract binding to set the session for
	CallOpts     bind.CallOpts               // Call options to use throughout this session
	TransactOpts bind.TransactOpts           // Transaction auth options to use throughout this session
}

// IZKMECrossChainUpgradeableCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type IZKMECrossChainUpgradeableCallerSession struct {
	Contract *IZKMECrossChainUpgradeableCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts                     // Call options to use throughout this session
}

// IZKMECrossChainUpgradeableTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type IZKMECrossChainUpgradeableTransactorSession struct {
	Contract     *IZKMECrossChainUpgradeableTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts                     // Transaction auth options to use throughout this session
}

// IZKMECrossChainUpgradeableRaw is an auto generated low-level Go binding around an Ethereum contract.
type IZKMECrossChainUpgradeableRaw struct {
	Contract *IZKMECrossChainUpgradeable // Generic contract binding to access the raw methods on
}

// IZKMECrossChainUpgradeableCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type IZKMECrossChainUpgradeableCallerRaw struct {
	Contract *IZKMECrossChainUpgradeableCaller // Generic read-only contract binding to access the raw methods on
}

// IZKMECrossChainUpgradeableTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type IZKMECrossChainUpgradeableTransactorRaw struct {
	Contract *IZKMECrossChainUpgradeableTransactor // Generic write-only contract binding to access the raw methods on
}

// NewIZKMECrossChainUpgradeable creates a new instance of IZKMECrossChainUpgradeable, bound to a specific deployed contract.
func NewIZKMECrossChainUpgradeable(address common.Address, backend bind.ContractBackend) (*IZKMECrossChainUpgradeable, error) {
	contract, err := bindIZKMECrossChainUpgradeable(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &IZKMECrossChainUpgradeable{IZKMECrossChainUpgradeableCaller: IZKMECrossChainUpgradeableCaller{contract: contract}, IZKMECrossChainUpgradeableTransactor: IZKMECrossChainUpgradeableTransactor{contract: contract}, IZKMECrossChainUpgradeableFilterer: IZKMECrossChainUpgradeableFilterer{contract: contract}}, nil
}

// NewIZKMECrossChainUpgradeableCaller creates a new read-only instance of IZKMECrossChainUpgradeable, bound to a specific deployed contract.
func NewIZKMECrossChainUpgradeableCaller(address common.Address, caller bind.ContractCaller) (*IZKMECrossChainUpgradeableCaller, error) {
	contract, err := bindIZKMECrossChainUpgradeable(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &IZKMECrossChainUpgradeableCaller{contract: contract}, nil
}

// NewIZKMECrossChainUpgradeableTransactor creates a new write-only instance of IZKMECrossChainUpgradeable, bound to a specific deployed contract.
func NewIZKMECrossChainUpgradeableTransactor(address common.Address, transactor bind.ContractTransactor) (*IZKMECrossChainUpgradeableTransactor, error) {
	contract, err := bindIZKMECrossChainUpgradeable(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &IZKMECrossChainUpgradeableTransactor{contract: contract}, nil
}

// NewIZKMECrossChainUpgradeableFilterer creates a new log filterer instance of IZKMECrossChainUpgradeable, bound to a specific deployed contract.
func NewIZKMECrossChainUpgradeableFilterer(address common.Address, filterer bind.ContractFilterer) (*IZKMECrossChainUpgradeableFilterer, error) {
	contract, err := bindIZKMECrossChainUpgradeable(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &IZKMECrossChainUpgradeableFilterer{contract: contract}, nil
}

// bindIZKMECrossChainUpgradeable binds a generic wrapper to an already deployed contract.
func bindIZKMECrossChainUpgradeable(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := IZKMECrossChainUpgradeableMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IZKMECrossChainUpgradeable.Contract.IZKMECrossChainUpgradeableCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IZKMECrossChainUpgradeable.Contract.IZKMECrossChainUpgradeableTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IZKMECrossChainUpgradeable.Contract.IZKMECrossChainUpgradeableTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _IZKMECrossChainUpgradeable.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _IZKMECrossChainUpgradeable.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _IZKMECrossChainUpgradeable.Contract.contract.Transact(opts, method, params...)
}

// GetCrossChainStatus is a free data retrieval call binding the contract method 0xfbdafb72.
//
// Solidity: function getCrossChainStatus(uint32 chainId, address user) view returns(uint256)
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableCaller) GetCrossChainStatus(opts *bind.CallOpts, chainId uint32, user common.Address) (*big.Int, error) {
	var out []interface{}
	err := _IZKMECrossChainUpgradeable.contract.Call(opts, &out, "getCrossChainStatus", chainId, user)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetCrossChainStatus is a free data retrieval call binding the contract method 0xfbdafb72.
//
// Solidity: function getCrossChainStatus(uint32 chainId, address user) view returns(uint256)
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableSession) GetCrossChainStatus(chainId uint32, user common.Address) (*big.Int, error) {
	return _IZKMECrossChainUpgradeable.Contract.GetCrossChainStatus(&_IZKMECrossChainUpgradeable.CallOpts, chainId, user)
}

// GetCrossChainStatus is a free data retrieval call binding the contract method 0xfbdafb72.
//
// Solidity: function getCrossChainStatus(uint32 chainId, address user) view returns(uint256)
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableCallerSession) GetCrossChainStatus(chainId uint32, user common.Address) (*big.Int, error) {
	return _IZKMECrossChainUpgradeable.Contract.GetCrossChainStatus(&_IZKMECrossChainUpgradeable.CallOpts, chainId, user)
}

// AckMinted is a paid mutator transaction binding the contract method 0x3ae65e32.
//
// Solidity: function ackMinted(uint32 chainId, address user, uint256 tokenId) returns()
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableTransactor) AckMinted(opts *bind.TransactOpts, chainId uint32, user common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _IZKMECrossChainUpgradeable.contract.Transact(opts, "ackMinted", chainId, user, tokenId)
}

// AckMinted is a paid mutator transaction binding the contract method 0x3ae65e32.
//
// Solidity: function ackMinted(uint32 chainId, address user, uint256 tokenId) returns()
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableSession) AckMinted(chainId uint32, user common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _IZKMECrossChainUpgradeable.Contract.AckMinted(&_IZKMECrossChainUpgradeable.TransactOpts, chainId, user, tokenId)
}

// AckMinted is a paid mutator transaction binding the contract method 0x3ae65e32.
//
// Solidity: function ackMinted(uint32 chainId, address user, uint256 tokenId) returns()
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableTransactorSession) AckMinted(chainId uint32, user common.Address, tokenId *big.Int) (*types.Transaction, error) {
	return _IZKMECrossChainUpgradeable.Contract.AckMinted(&_IZKMECrossChainUpgradeable.TransactOpts, chainId, user, tokenId)
}

// Forward is a paid mutator transaction binding the contract method 0x04b2c8df.
//
// Solidity: function forward(uint32 chainId, address srcUser, address destUser) returns()
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableTransactor) Forward(opts *bind.TransactOpts, chainId uint32, srcUser common.Address, destUser common.Address) (*types.Transaction, error) {
	return _IZKMECrossChainUpgradeable.contract.Transact(opts, "forward", chainId, srcUser, destUser)
}

// Forward is a paid mutator transaction binding the contract method 0x04b2c8df.
//
// Solidity: function forward(uint32 chainId, address srcUser, address destUser) returns()
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableSession) Forward(chainId uint32, srcUser common.Address, destUser common.Address) (*types.Transaction, error) {
	return _IZKMECrossChainUpgradeable.Contract.Forward(&_IZKMECrossChainUpgradeable.TransactOpts, chainId, srcUser, destUser)
}

// Forward is a paid mutator transaction binding the contract method 0x04b2c8df.
//
// Solidity: function forward(uint32 chainId, address srcUser, address destUser) returns()
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableTransactorSession) Forward(chainId uint32, srcUser common.Address, destUser common.Address) (*types.Transaction, error) {
	return _IZKMECrossChainUpgradeable.Contract.Forward(&_IZKMECrossChainUpgradeable.TransactOpts, chainId, srcUser, destUser)
}

// IZKMECrossChainUpgradeableZkmeSBTCrossChainPackageIterator is returned from FilterZkmeSBTCrossChainPackage and is used to iterate over the raw logs and unpacked data for ZkmeSBTCrossChainPackage events raised by the IZKMECrossChainUpgradeable contract.
type IZKMECrossChainUpgradeableZkmeSBTCrossChainPackageIterator struct {
	Event *IZKMECrossChainUpgradeableZkmeSBTCrossChainPackage // Event containing the contract specifics and raw log

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
func (it *IZKMECrossChainUpgradeableZkmeSBTCrossChainPackageIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(IZKMECrossChainUpgradeableZkmeSBTCrossChainPackage)
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
		it.Event = new(IZKMECrossChainUpgradeableZkmeSBTCrossChainPackage)
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
func (it *IZKMECrossChainUpgradeableZkmeSBTCrossChainPackageIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *IZKMECrossChainUpgradeableZkmeSBTCrossChainPackageIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// IZKMECrossChainUpgradeableZkmeSBTCrossChainPackage represents a ZkmeSBTCrossChainPackage event raised by the IZKMECrossChainUpgradeable contract.
type IZKMECrossChainUpgradeableZkmeSBTCrossChainPackage struct {
	EventData KYCDataLibEventData
	Raw       types.Log // Blockchain specific contextual infos
}

// FilterZkmeSBTCrossChainPackage is a free log retrieval operation binding the contract event 0x1e2dd7094825f10dc568deef8da4b8efff58a93840fc76f4ead206f6c8c5cb82.
//
// Solidity: event ZkmeSBTCrossChainPackage((uint32,uint32,uint256,uint256,bytes) indexed eventData)
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableFilterer) FilterZkmeSBTCrossChainPackage(opts *bind.FilterOpts, eventData []KYCDataLibEventData) (*IZKMECrossChainUpgradeableZkmeSBTCrossChainPackageIterator, error) {

	var eventDataRule []interface{}
	for _, eventDataItem := range eventData {
		eventDataRule = append(eventDataRule, eventDataItem)
	}

	logs, sub, err := _IZKMECrossChainUpgradeable.contract.FilterLogs(opts, "ZkmeSBTCrossChainPackage", eventDataRule)
	if err != nil {
		return nil, err
	}
	return &IZKMECrossChainUpgradeableZkmeSBTCrossChainPackageIterator{contract: _IZKMECrossChainUpgradeable.contract, event: "ZkmeSBTCrossChainPackage", logs: logs, sub: sub}, nil
}

// WatchZkmeSBTCrossChainPackage is a free log subscription operation binding the contract event 0x1e2dd7094825f10dc568deef8da4b8efff58a93840fc76f4ead206f6c8c5cb82.
//
// Solidity: event ZkmeSBTCrossChainPackage((uint32,uint32,uint256,uint256,bytes) indexed eventData)
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableFilterer) WatchZkmeSBTCrossChainPackage(opts *bind.WatchOpts, sink chan<- *IZKMECrossChainUpgradeableZkmeSBTCrossChainPackage, eventData []KYCDataLibEventData) (event.Subscription, error) {

	var eventDataRule []interface{}
	for _, eventDataItem := range eventData {
		eventDataRule = append(eventDataRule, eventDataItem)
	}

	logs, sub, err := _IZKMECrossChainUpgradeable.contract.WatchLogs(opts, "ZkmeSBTCrossChainPackage", eventDataRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(IZKMECrossChainUpgradeableZkmeSBTCrossChainPackage)
				if err := _IZKMECrossChainUpgradeable.contract.UnpackLog(event, "ZkmeSBTCrossChainPackage", log); err != nil {
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

// ParseZkmeSBTCrossChainPackage is a log parse operation binding the contract event 0x1e2dd7094825f10dc568deef8da4b8efff58a93840fc76f4ead206f6c8c5cb82.
//
// Solidity: event ZkmeSBTCrossChainPackage((uint32,uint32,uint256,uint256,bytes) indexed eventData)
func (_IZKMECrossChainUpgradeable *IZKMECrossChainUpgradeableFilterer) ParseZkmeSBTCrossChainPackage(log types.Log) (*IZKMECrossChainUpgradeableZkmeSBTCrossChainPackage, error) {
	event := new(IZKMECrossChainUpgradeableZkmeSBTCrossChainPackage)
	if err := _IZKMECrossChainUpgradeable.contract.UnpackLog(event, "ZkmeSBTCrossChainPackage", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
