package main

import (
	"bnbchain/opbnb-bridge-bot/core"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ethereum-optimism/optimism/indexer/config"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum-optimism/optimism/op-service/retry"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/log"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func RunCommand(ctx *cli.Context) error {
	// Config logger
	logger := oplog.NewLogger(oplog.AppOut(ctx), oplog.ReadCLIConfig(ctx)).New("role", "bot")
	oplog.SetGlobalLogHandler(logger.GetHandler())
	logger.Info("opbnb-bridge-bot is starting")

	// Load config
	cfg, err := core.LoadConfig(logger, ctx.String(ConfigFlag.Name))
	if err != nil {
		logger.Error("failed to load config", "err", err)
		return err
	}

	// Connect to L1 and L2 RPCs
	l1Client, err := core.Dial(cfg.RPCs.L1RPC)
	if err != nil {
		return fmt.Errorf("dial endpoint %s: %w", cfg.RPCs.L1RPC, err)
	}
	l2Client, err := core.Dial(cfg.RPCs.L2RPC)
	if err != nil {
		return fmt.Errorf("dial endpoint %s: %w", cfg.RPCs.L2RPC, err)
	}

	// Connect to database and ensure schemas initialized
	db, err := connect(logger, cfg.DB)
	if err != nil {
		return fmt.Errorf("failed to connect database: %w", err)
	}

	err = db.AutoMigrate(&core.L2ScannedBlock{})
	if err != nil {
		return fmt.Errorf("failed to migrate l2_scanned_blocks: %w", err)
	}
	err = db.AutoMigrate(&core.WithdrawalInitiatedLog{})
	if err != nil {
		return fmt.Errorf("failed to migrate withdrawals: %w", err)
	}

	l2ScannedBlock, err := queryL2ScannedBlock(db, cfg.L2StartingNumber)
	if err != nil {
		return fmt.Errorf("failed to query l2_scanned_blocks: %w", err)
	} else {
		logger.Info("starting from block", "blockNumber", l2ScannedBlock.Number)
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.Handle("/debug/metrics/prometheus", promhttp.Handler())
		err := http.ListenAndServe(":6060", nil)
		if err != nil {
			logger.Error("failed to start prometheus server", "error", err)
		}
	}()

	go core.StartMetrics(ctx.Context, &cfg, &l1Client.Client, db, logger)
	go func() {
		indexer := core.NewIndexer(logger, db, l2Client, cfg)
		indexer.Start(ctx.Context, l2ScannedBlock)
	}()
	go ProcessBotDelegatedWithdrawals(ctx.Context, logger, db, l1Client, l2Client, cfg)

	<-ctx.Context.Done()
	return nil
}

// ProcessBotDelegatedWithdrawals processes the indexed bot-delegated withdrawals.
// It will prove the withdrawal transaction when the proposal time window has passed;
// and it will finalize the withdrawal when the challenge time window has passed.
func ProcessBotDelegatedWithdrawals(ctx context.Context, log log.Logger, db *gorm.DB, l1Client *core.ClientExt, l2Client *core.ClientExt, cfg core.Config) {
	ticker := time.NewTicker(3 * time.Second)
	pending := core.NewPendingTxsManager()
	for {
		select {
		case <-ticker.C:
			// In order to avoid re-processing the same withdrawal, we need to check if the pending nonce is
			// the chain nonce. If they are not equal, it means that there are some pending transactions that
			// been confirmed yet.
			_, signerAddress, _ := cfg.SignerKeyPair()
			if equal, err := isPendingAndChainNonceEqual(l1Client, signerAddress); err != nil {
				log.Error("failed to check pending and chain nonce", "error", err)
				continue
			} else if !equal {
				log.Info("pending nonce is not equal to chain nonce, skip processing")
				continue
			}

			currentNonce, err := l1Client.NonceAt(ctx, *signerAddress, nil)
			if err != nil {
				log.Error("failed to get chain nonce", "error", err)
				continue
			}

			ProcessUnprovenBotDelegatedWithdrawals(ctx, log, db, l1Client, l2Client, cfg, pending, &currentNonce)
			ProcessUnfinalizedBotDelegatedWithdrawals(ctx, log, db, l1Client, l2Client, cfg, pending, &currentNonce)
		case <-ctx.Done():
			return
		}
	}
}

func ProcessUnprovenBotDelegatedWithdrawals(ctx context.Context, log log.Logger, db *gorm.DB, l1Client *core.ClientExt, l2Client *core.ClientExt, cfg core.Config, pending *core.PendingTxnCheck, currentNonce *uint64) {
	latestProposedNumber, err := core.L2OutputOracleLatestBlockNumber(cfg.L1Contracts.L2OutputOracleProxy, l1Client)
	if err != nil {
		log.Error("failed to get latest proposed block number", "error", err)
		return
	}

	processor := core.NewProcessor(log, l1Client, l2Client, cfg)
	limit := 1000

	unprovens := make([]core.WithdrawalInitiatedLog, 0)
	result := db.Order("id asc").Where("proven_time IS NULL AND initiated_block_number <= ? AND failure_reason IS NULL", latestProposedNumber.Uint64()).Limit(limit).Find(&unprovens)
	if result.Error != nil {
		log.Error("failed to query withdrawals", "error", result.Error)
		return
	}

	pending.Prune(*currentNonce)
	for _, unproven := range unprovens {
		// Avoid re-processing the same withdrawal
		if pending.IsPendingTxn(unproven.ID) {
			continue
		}

		now := time.Now()
		err := processor.ProveWithdrawalTransaction(ctx, &unproven, *currentNonce)
		if err != nil {
			if strings.Contains(err.Error(), "OptimismPortal: withdrawal hash has already been proven") {
				// The withdrawal has already proven, mark it
				result := db.Model(&unproven).Update("proven_time", now)
				if result.Error != nil {
					log.Error("failed to update proven withdrawals", "error", result.Error)
				}
			} else if strings.Contains(err.Error(), "L2OutputOracle: cannot get output for a block that has not been proposed") {
				// Since the unproven withdrawals are sorted by the on-chain order, we can break here because we know
				// that the subsequent of the withdrawals are not ready to be proven yet.
				return
			} else if strings.Contains(err.Error(), "execution reverted") || strings.Contains(err.Error(), "filtered") {
				// Proven transaction reverted, mark it with the failure reason
				result := db.Model(&unproven).Update("failure_reason", err.Error())
				if result.Error != nil {
					log.Error("failed to update failure reason of withdrawals", "error", result.Error)
				}
			} else {
				// non-revert error, stop processing the subsequent withdrawals
				log.Error("ProveWithdrawalTransaction", "non-revert error", err.Error())
				return
			}
		} else {
			pending.AddPendingTxn(unproven.ID, *currentNonce)
			*currentNonce = *currentNonce + 1
		}
	}
}

func ProcessUnfinalizedBotDelegatedWithdrawals(ctx context.Context, log log.Logger, db *gorm.DB, l1Client *core.ClientExt, l2Client *core.ClientExt, cfg core.Config, pending *core.PendingTxnCheck, currentNonce *uint64) {
	processor := core.NewProcessor(log, l1Client, l2Client, cfg)
	limit := 1000

	now := time.Now()
	maxProvenTime := now.Add(-time.Duration(cfg.ChallengeTimeWindow) * time.Second)

	unfinalizeds := make([]core.WithdrawalInitiatedLog, 0)
	result := db.Order("id asc").Where("finalized_time IS NULL AND proven_time IS NOT NULL AND proven_time < ? AND failure_reason IS NULL", maxProvenTime).Limit(limit).Find(&unfinalizeds)
	if result.Error != nil {
		log.Error("failed to query withdrawals", "error", result.Error)
		return
	}

	pending.Prune(*currentNonce)
	for _, unfinalized := range unfinalizeds {
		// In order to avoid re-processing the same withdrawal
		if pending.IsPendingTxn(unfinalized.ID) {
			continue
		}

		err := processor.FinalizeMessage(ctx, &unfinalized)
		if err != nil {
			if strings.Contains(err.Error(), "OptimismPortal: withdrawal has already been finalized") {
				// The withdrawal has already finalized, mark it
				result := db.Model(&unfinalized).Update("finalized_time", now)
				if result.Error != nil {
					log.Error("failed to update finalized withdrawals", "error", result.Error)
				}
			} else if strings.Contains(err.Error(), "OptimismPortal: withdrawal has not been proven yet") {
				log.Error("detected a unproven withdrawal when send finalized transaction", "withdrawal", unfinalized)
				continue
			} else if strings.Contains(err.Error(), "OptimismPortal: proven withdrawal finalization period has not elapsed") {
				// Continue to handle the subsequent unfinalized withdrawals
				continue
			} else if strings.Contains(err.Error(), "execution reverted") {
				// Finalized transaction reverted, mark it with the failure reason
				result := db.Model(&unfinalized).Update("failure_reason", err.Error())
				if result.Error != nil {
					log.Error("failed to update failure reason of withdrawals", "error", result.Error)
				}
			} else {
				// non-revert error, stop processing the subsequent withdrawals
				log.Error("FinalizedMessage", "non-revert error", err.Error(), "wi", unfinalized.ID)
				return
			}
		} else {
			pending.AddPendingTxn(unfinalized.ID, *currentNonce)
			*currentNonce = *currentNonce + 1
		}
	}
}

// connect connects to the database
func connect(log log.Logger, dbConfig config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local", dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.Name)
	gormConfig := gorm.Config{
		Logger:                 core.NewGormLogger(log),
		SkipDefaultTransaction: true,
		CreateBatchSize:        3_000,
	}

	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	gorm_, err := retry.Do[*gorm.DB](context.Background(), 10, retryStrategy, func() (*gorm.DB, error) {
		gorm_, err := gorm.Open(mysql.Open(dsn), &gormConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}

		return gorm_, nil
	})

	if err != nil {
		return nil, err
	}

	return gorm_, nil
}

// queryL2ScannedBlock queries the l2_scanned_blocks table for the last scanned block
func queryL2ScannedBlock(db *gorm.DB, l2StartingNumber int64) (*core.L2ScannedBlock, error) {
	l2ScannedBlock := core.L2ScannedBlock{}
	if result := db.Order("number desc").Last(&l2ScannedBlock); result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			l2ScannedBlock.Number = l2StartingNumber
			db.Create(&l2ScannedBlock)
			return &l2ScannedBlock, nil
		} else {
			return nil, fmt.Errorf("failed to query l2_scanned_blocks: %w", result.Error)
		}
	}
	return &l2ScannedBlock, nil
}

// isPendingAndChainNonceEqual checks if the pending nonce and the chain nonce are equal.
func isPendingAndChainNonceEqual(l1Client *core.ClientExt, address *common.Address) (bool, error) {
	pendingNonce, err := l1Client.PendingNonceAt(context.Background(), *address)
	if err != nil {
		return false, fmt.Errorf("failed to get pending nonce: %w", err)
	}

	latestNonce, err := l1Client.NonceAt(context.Background(), *address, nil)
	if err != nil {
		return false, fmt.Errorf("failed to get latest nonce: %w", err)
	}

	return pendingNonce == latestNonce, nil
}
