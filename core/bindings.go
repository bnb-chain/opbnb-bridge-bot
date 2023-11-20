package core

import (
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
