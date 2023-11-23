package core

type L2ScannedBlock struct {
	Number int64 `gorm:"type:integer;primarykey"`
}

type L2ContractEvent struct {
	ID              uint   `gorm:"primarykey"`
	BlockTime       int64  `gorm:"type:integer;not null;index:idx_l2_contract_events_block_time"`
	BlockHash       string `gorm:"type:varchar(256);not null;uniqueIndex:idx_l2_contract_events_block_hash_log_index_key,priority:1;"`
	ContractAddress string `gorm:"type:varchar(256);not null"`
	TransactionHash string `gorm:"type:varchar(256);not null"`
	LogIndex        int    `gorm:"type:integer;not null;uniqueIndex:idx_l2_contract_events_block_hash_log_index_key,priority:2"`
	EventSignature  string `gorm:"type:varchar(256);not null"`
	Proven          bool   `gorm:"type:boolean;not null;default:false"`
	Finalized       bool   `gorm:"type:boolean;not null;default:false"`
	FailureReason   string `gorm:"type:varchar(256);"`
}
