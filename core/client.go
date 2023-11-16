package core

import (
	"context"
	"math/big"
	"time"

	"github.com/ethereum-optimism/optimism/op-service/eth"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	// defaultDialTimeout is default duration the processor will wait on
	// startup to make a connection to the backend
	defaultDialTimeout = 5 * time.Second

	// defaultDialAttempts is the default attempts a connection will be made
	// before failing
	defaultDialAttempts = 5

	// defaultRequestTimeout is the default duration the processor will
	// wait for a request to be fulfilled
	defaultRequestTimeout = 10 * time.Second
)

type ClientExt struct {
	ethclient.Client
}

var _ bind.ContractCaller = (*ClientExt)(nil)

func Dial(rawUrl string) (*ClientExt, error) {
	client, err := ethclient.DialContext(context.Background(), rawUrl)
	if err != nil {
		return nil, err
	}

	return &ClientExt{*client}, nil
}

func (c *ClientExt) GetProof(address common.Address, storageKeys []string, blockNumber *big.Int) (*eth.AccountResult, error) {
	ctxwt, cancel := context.WithTimeout(context.Background(), defaultRequestTimeout)
	defer cancel()

	var result eth.AccountResult
	err := c.Client.Client().CallContext(ctxwt, &result, "eth_getProof", address, storageKeys, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// Impl bind.ContractCaller

func (c *ClientExt) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	var result string
	err := c.Client.Client().CallContext(ctx, &result, "eth_getCode", contract, toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}

	return hexutil.Decode(result)
}

func (c *ClientExt) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var result string
	err := c.Client.Client().CallContext(ctx, &result, "eth_call", toCallArg(call), toBlockNumArg(blockNumber))
	if err != nil {
		return nil, err
	}

	return hexutil.Decode(result)
}

// Needed private utils from geth

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	if number.Sign() >= 0 {
		return hexutil.EncodeBig(number)
	}
	// It's negative.
	return rpc.BlockNumber(number.Int64()).String()
}

func toCallArg(msg ethereum.CallMsg) interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["input"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}
