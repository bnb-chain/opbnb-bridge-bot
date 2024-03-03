package core

import (
	"context"
	"errors"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"gorm.io/gorm"
	"time"
)

var (
	TxSignerBalance = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "opbnb_bridge_bot_tx_signer_balance",
		Help: "The balance of the tx signer",
	})

	ScannedBlockNumber = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "opbnb_bridge_bot_scanned_block_number",
		Help: "The block number that has been scanned",
	})

	UnprovenWithdrawals = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "opbnb_bridge_bot_unproven_withdrawals",
		Help: "The number of unproven withdrawals",
	})
	UnfinalizedWithdrawals = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "opbnb_bridge_bot_unfinalized_withdrawals",
		Help: "The number of unfinalized withdrawals",
	})
	EarliestUnProvenWithdrawalBlockNumber = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "opbnb_bridge_bot_earliest_unproven_withdrawal_block_number",
		Help: "The earliest block number of unproven withdrawals",
	})
	EarliestUnfinalizedWithdrawalBlockNumber = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "opbnb_bridge_bot_earliest_unfinalized_withdrawal_block_number",
		Help: "The earliest block number of unfinalized withdrawals",
	})

	FailedWithdrawals = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "opbnb_bridge_bot_failed_withdrawals",
		Help: "The number of failed withdrawals",
	})
)

func StartMetrics(ctx context.Context, cfg *Config, l1Client *ethclient.Client, db *gorm.DB, logger log.Logger) {
	ticker := time.NewTicker(5 * time.Second)
	for range ticker.C {
		_, signerAddress, _ := cfg.SignerKeyPair()
		balance, err := l1Client.BalanceAt(ctx, *signerAddress, nil)
		if err != nil {
			logger.Error("failed to get signer balance", "error", err)
		}
		TxSignerBalance.Set(float64(balance.Int64()))

		var scannedBlock L2ScannedBlock
		result := db.Last(&scannedBlock)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Error("failed to query scanned block", "error", result.Error)
		}
		ScannedBlockNumber.Set(float64(scannedBlock.Number))

		var unprovenCnt int64
		result = db.Table("withdrawal_initiated_logs").Where("proven_time IS NULL AND failure_reason IS NULL").Count(&unprovenCnt)
		if result.Error != nil {
			logger.Error("failed to count withdrawals", "error", result.Error)
		}
		UnprovenWithdrawals.Set(float64(unprovenCnt))

		var unfinalizedCnt int64
		result = db.Table("withdrawal_initiated_logs").Where("finalized_time IS NULL AND proven_time IS NOT NULL AND failure_reason IS NULL").Count(&unfinalizedCnt)
		if result.Error != nil {
			logger.Error("failed to count withdrawals", "error", result.Error)
		}
		UnfinalizedWithdrawals.Set(float64(unfinalizedCnt))

		var failedCnt int64
		result = db.Table("withdrawal_initiated_logs").Where("failure_reason IS NOT NULL").Count(&failedCnt)
		if result.Error != nil {
			logger.Error("failed to count withdrawals", "error", result.Error)
		}
		FailedWithdrawals.Set(float64(failedCnt))

		firstUnproven := WithdrawalInitiatedLog{}
		result = db.Table("withdrawal_initiated_logs").Order("id asc").Where("proven_time IS NULL AND failure_reason IS NULL").First(&firstUnproven)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Error("failed to query withdrawals", "error", result.Error)
		}
		EarliestUnProvenWithdrawalBlockNumber.Set(float64(firstUnproven.InitiatedBlockNumber))

		firstUnfinalized := WithdrawalInitiatedLog{}
		result = db.Table("withdrawal_initiated_logs").Order("id asc").Where("finalized_time IS NULL AND proven_time IS NOT NULL AND failure_reason IS NULL").First(&firstUnfinalized)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logger.Error("failed to query withdrawals", "error", result.Error)
		}
		EarliestUnfinalizedWithdrawalBlockNumber.Set(float64(firstUnfinalized.InitiatedBlockNumber))
	}
}
