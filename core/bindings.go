package core

import (
	"math/big"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/common"
	"golang.org/x/crypto/sha3"
)

// WithdrawToEventSig is the signature for the WithdrawTo event:
//
// ```
// WithdrawTo(address indexed from, address l2Token, address to, uint256 amount, uint32 minGasLimit, bytes extraData);
// ```
func WithdrawToEventSig() common.Hash {
	eventSignature := "WithdrawTo(address,address,address,uint256,uint32,bytes)"
	keccak256 := sha3.NewLegacyKeccak256()
	keccak256.Write([]byte(eventSignature))
	eventSignatureHash := keccak256.Sum(nil)
	return common.BytesToHash(eventSignatureHash)
}

// L2OutputOracleLatestBlockNumber calls the "latestBlockNumber" function on the L2OutputOracle contract at the given address.
func L2OutputOracleLatestBlockNumber(address common.Address, l1Client *ClientExt) (*big.Int, error) {
	caller, err := bindings.NewL2OutputOracleCaller(address, l1Client)
	if err != nil {
		return nil, err
	}

	return caller.LatestBlockNumber(nil)
}
