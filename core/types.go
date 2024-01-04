package core

import "time"

type L2ScannedBlock struct {
	Number int64 `gorm:"type:integer;primarykey"`
}

type BotDelegatedWithdrawal struct {
	// ID is the incrementing primary key.
	ID uint `gorm:"primarykey"`

	// TransactionHash and LogIndex are the L2 transaction hash and log index of the withdrawal event.
	TransactionHash string `gorm:"type:varchar(256);not null;uniqueIndex:idx_bot_delegated_withdrawals_transaction_hash_log_index_key,priority:1"`
	LogIndex        int    `gorm:"type:integer;not null;uniqueIndex:idx_bot_delegated_withdrawals_transaction_hash_log_index_key,priority:2"`

	// InitiatedBlockNumber is the l2 block number at which the withdrawal was initiated on L2.
	InitiatedBlockNumber int64 `gorm:"type:integer;not null;index:idx_withdrawals_initiated_block_number"`

	// ProvenTime is the local time at which the withdrawal was proven on L1. NULL if not yet proven.
	ProvenTime *time.Time `gorm:"type:datetime;index:idx_withdrawals_proven_time"`

	// FinalizedTime is the local time at which the withdrawal was finalized on L1. NULL if not yet finalized.
	FinalizedTime *time.Time `gorm:"type:datetime;index:idx_withdrawals_finalized_time"`

	// FailureReason is the reason for the withdrawal failure, including sending transaction error and off-chain configured filter error. NULL if not yet failed.
	FailureReason *string `gorm:"type:text"`
}
