package core

type L2ScannedBlock struct {
	Number int64 `gorm:"type:integer"`
}

type L2ContractEvent struct {
	ID              uint   `gorm:"primarykey"`
	BlockTime       int64  `gorm:"type:integer;not null;index"`
	BlockHash       string `gorm:"type:varchar(256);not null;index"`
	ContractAddress string `gorm:"type:varchar(256);not nul"`
	TransactionHash string `gorm:"type:varchar(256);not null;index"`
	LogIndex        int    `gorm:"type:integer;not null"`
	EventSignature  string `gorm:"type:varchar(256);not null;index"`
	Proven          bool   `gorm:"type:boolean;not null;default:false"`
	Finalized       bool   `gorm:"type:boolean;not null;default:false"`
}
