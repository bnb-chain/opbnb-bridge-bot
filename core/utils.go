package core

import "github.com/ethereum/go-ethereum/core/types"

// GetLogByLogIndex searches through receipt.Logs using the provided logIndex and return the corresponding log.
//
// Be aware that:
// - Log.Index is accumulated for the entire block.
// - There is no guarantee that receipt.Logs will be ordered by index.
func GetLogByLogIndex(receipt *types.Receipt, logIndex uint) *types.Log {
	if receipt == nil {
		return nil
	}

	for _, vlog := range receipt.Logs {
		if vlog.Index == logIndex {
			return vlog
		}
	}
	return nil
}
