package core

import "time"

type L2ScannedBlockV2 struct {
	Number int64 `gorm:"type:integer;primarykey"`
}

// WithdrawalInitiatedLog is parsed record that represent a withdrawal.
//
// See also [L2StandardBridge.WithdrawalInitiated](https://github.com/bnb-chain/opbnb/blob/develop/packages/contracts-bedrock/contracts/L2/L2StandardBridge.sol#L21-L39)
type WithdrawalInitiatedLog struct {
	// ID is the incrementing primary key.
	ID uint `gorm:"primarykey"`

	// TransactionHash and LogIndex are the L2 transaction hash and log index of the withdrawal event.
	TransactionHash string `gorm:"type:varchar(256);not null;uniqueIndex:idx_withdrawal_initiated_logs_transaction_hash_log_index_key,priority:1"`
	LogIndex        int    `gorm:"type:integer;not null;uniqueIndex:idx_withdrawal_initiated_logs_transaction_hash_log_index_key,priority:2"`

	// InitiatedBlockNumber is the l2 block number at which the withdrawal was initiated on L2.
	InitiatedBlockNumber int64 `gorm:"type:integer;not null;index:idx_withdrawal_initiated_logs_initiated_block_number"`

	// ProvenTime is the local time at which the withdrawal was proven on L1. NULL if not yet proven.
	ProvenTime *time.Time `gorm:"type:datetime;index:idx_withdrawal_initiated_logs_proven_time"`

	// FinalizedTime is the local time at which the withdrawal was finalized on L1. NULL if not yet finalized.
	FinalizedTime *time.Time `gorm:"type:datetime;index:idx_withdrawal_initiated_logs_finalized_time"`

	// FailureReason is the reason for the withdrawal failure, including sending transaction error and off-chain configured filter error. NULL if not yet failed.
	FailureReason *string `gorm:"type:text"`
}
