package main

import (
	"bnbchain/opbnb-bridge-bot/core"
	"context"
	"errors"
	"fmt"
	"math/big"
	"strings"
	"time"

	"github.com/ethereum-optimism/optimism/indexer/config"
	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum-optimism/optimism/op-service/retry"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/urfave/cli/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func runCommand(ctx *cli.Context) error {
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
	err = db.AutoMigrate(&core.L2ContractEvent{})
	if err != nil {
		return fmt.Errorf("failed to migrate l2_contract_events: %w", err)
	}

	l2ScannedBlock, err := queryL2ScannedBlock(db, &cfg)
	if err != nil {
		return err
	}

	go WatchBotDelegatedWithdrawals(ctx.Context, logger, db, l2Client, l2ScannedBlock, cfg)
	go ProcessBotDelegatedWithdrawals(ctx.Context, logger, db, l1Client, l2Client, cfg)

	<-ctx.Context.Done()
	return nil
}

// ProcessBotDelegatedWithdrawals processes the indexed bot-delegated withdrawals.
// It will prove the withdrawal transaction when the proposal time window has passed;
// and it will finalize the withdrawal when the challenge time window has passed.
func ProcessBotDelegatedWithdrawals(ctx context.Context, log log.Logger, db *gorm.DB, l1Client *core.ClientExt, l2Client *core.ClientExt, cfg core.Config) {
	ticker := time.NewTicker(3 * time.Second)
	processor := core.NewProcessor(log, l1Client, l2Client, cfg)
	limit := 1000

	for {
		select {
		case <-ticker.C:
		case <-ctx.Done():
			return
		}

		unprovens := make([]core.L2ContractEvent, 0)
		if result := db.Order("id asc").Where("proven = false").Limit(limit).Find(&unprovens); result.Error != nil {
			log.Error("querying l2_contract_events", "error", result.Error)
			continue
		} else if len(unprovens) > 0 {
			for _, unproven := range unprovens {
				if unproven.BlockTime+cfg.Misc.ProposeTimeWindow < time.Now().Unix() {
					err := processor.ProveWithdrawalTransaction(&unproven)
					if err != nil {
						if strings.Contains(err.Error(), "OptimismPortal: withdrawal hash has already been proven") {
							// The withdrawal has already proven, mark it
							unproven.Proven = true
							result := db.Save(&unproven)
							if result.Error != nil {
								log.Crit("update proven l2_contract_events", "error", result.Error)
							}
						} else if strings.Contains(err.Error(), "L2OutputOracle: cannot get output for a block that has not been proposed") {
							break
						} else {
							log.Crit("ProveWithdrawalTransaction", "error", err.Error())
						}
					}
				}
			}
		}

		unfinalizeds := make([]core.L2ContractEvent, 0)
		if result := db.Order("block_time asc").Where("proven = true AND finalized = false").Limit(limit).Find(&unfinalizeds); result.Error != nil {
			log.Error("querying l2_contract_events", "error", result.Error)
			continue
		} else if len(unfinalizeds) > 0 {
			for _, unfinalized := range unfinalizeds {
				if unfinalized.BlockTime+cfg.Misc.ChallengeTimeWindow < time.Now().Unix() {
					err := processor.FinalizeMessage(&unfinalized)
					if err != nil {
						if strings.Contains(err.Error(), "OptimismPortal: withdrawal has already been finalized") {
							// The withdrawal has already finalized, mark it
							unfinalized.Finalized = true
							result := db.Save(&unfinalized)
							if result.Error != nil {
								log.Crit("update finalized l2_contract_events", "error", result.Error)
							}
						} else if strings.Contains(err.Error(), "OptimismPortal: withdrawal has not been proven yet") || strings.Contains(err.Error(), "OptimismPortal: proven withdrawal finalization period has not elapsed") {
							break
						} else {
							log.Crit("FinalizeMessage", "error", err.Error())
						}
					}
				}
			}
		}
	}
}

// storeLogs stores the logs in the database
func storeLogs(db *gorm.DB, client *core.ClientExt, logs []types.Log) error {
	// save all the logs in this range of blocks
	for _, vLog := range logs {
		header, err := client.HeaderByHash(context.Background(), vLog.BlockHash)
		if err != nil {
			return err
		}

		event := core.L2ContractEvent{
			BlockTime:       int64(header.Time),
			BlockHash:       vLog.BlockHash.Hex(),
			ContractAddress: vLog.Address.Hex(),
			TransactionHash: vLog.TxHash.Hex(),
			LogIndex:        int(vLog.Index),
			EventSignature:  vLog.Topics[0].Hex(),
		}

		deduped := db.Clauses(
			clause.OnConflict{DoNothing: true},
		)
		result := deduped.Create(&event)
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// WatchBotDelegatedWithdrawals watches for new bot-delegated withdrawals and stores them in the database.
func WatchBotDelegatedWithdrawals(ctx context.Context, log log.Logger, db *gorm.DB, client *core.ClientExt, l2ScannedBlock *core.L2ScannedBlock, cfg core.Config) {
	timer := time.NewTimer(0)
	fromBlockNumber := big.NewInt(l2ScannedBlock.Number)

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			timer.Reset(time.Second)
		}

		toBlockNumber := new(big.Int).Add(fromBlockNumber, big.NewInt(cfg.Misc.LogFilterBlockRange))
		latestNumber, err := client.BlockNumber(context.Background())
		if err != nil {
			log.Error("call eth_blockNumber", "error", err)
			continue
		}

		if latestNumber < uint64(cfg.Misc.ConfirmBlocks) {
			toBlockNumber = big.NewInt(0)
		} else if latestNumber-uint64(cfg.Misc.ConfirmBlocks) < toBlockNumber.Uint64() {
			toBlockNumber = big.NewInt(int64(latestNumber - uint64(cfg.Misc.ConfirmBlocks)))
		}

		if fromBlockNumber.Uint64() > toBlockNumber.Uint64() {
			timer.Reset(5 * time.Second)
			continue
		}

		log.Info("Fetching logs from blocks", "fromBlock", fromBlockNumber, "toBlock", toBlockNumber)
		logs, err := getLogs(client, fromBlockNumber, toBlockNumber, common.HexToAddress(cfg.Misc.L2StandardBridgeBot), core.WithdrawToEventSig())
		if err != nil {
			log.Error("eth_getLogs", "error", err)
			continue
		}

		if len(logs) != 0 {
			for _, vlog := range logs {
				log.Info("fetched bot-delegated withdrawal", "blockNumber", vlog.BlockNumber, "transactionHash", vlog.TxHash.Hex())
			}

			err = storeLogs(db, client, logs)
			if err != nil {
				log.Error("storeLogs", "error", err)
				continue
			}
		}

		l2ScannedBlock.Number = toBlockNumber.Int64()
		result := db.Where("number >= 0").Updates(l2ScannedBlock)
		if result.Error != nil {
			log.Error("update l2_scanned_blocks", "error", result.Error)
		}

		fromBlockNumber = new(big.Int).Add(toBlockNumber, big.NewInt(1))
	}
}

// getLogs returns the logs for a given contract address and block range
func getLogs(client *core.ClientExt, fromBlock *big.Int, toBlock *big.Int, contractAddress common.Address, eventSig common.Hash) ([]types.Log, error) {
	query := ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: []common.Address{
			contractAddress,
		},
		Topics: [][]common.Hash{[]common.Hash{eventSig}},
	}
	return client.FilterLogs(context.Background(), query)
}

// connect connects to the database
func connect(log log.Logger, dbConfig config.DBConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s dbname=%s sslmode=disable", dbConfig.Host, dbConfig.Name)
	if dbConfig.Port != 0 {
		dsn += fmt.Sprintf(" port=%d", dbConfig.Port)
	}
	if dbConfig.User != "" {
		dsn += fmt.Sprintf(" user=%s", dbConfig.User)
	}
	if dbConfig.Password != "" {
		dsn += fmt.Sprintf(" password=%s", dbConfig.Password)
	}

	gormConfig := gorm.Config{
		Logger:                 core.NewGormLogger(log),
		SkipDefaultTransaction: true,

		// The postgres parameter counter for a given query is represented with uint16,
		// resulting in a parameter limit of 65535. In order to avoid reaching this limit
		// we'll utilize a batch size of 3k for inserts, well below the limit as long as
		// the number of columns < 20.
		CreateBatchSize: 3_000,
	}

	retryStrategy := &retry.ExponentialStrategy{Min: 1000, Max: 20_000, MaxJitter: 250}
	gorm_, err := retry.Do[*gorm.DB](context.Background(), 10, retryStrategy, func() (*gorm.DB, error) {
		gorm_, err := gorm.Open(postgres.Open(dsn), &gormConfig)
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
func queryL2ScannedBlock(db *gorm.DB, cfg *core.Config) (*core.L2ScannedBlock, error) {
	l2ScannedBlock := core.L2ScannedBlock{Number: 0}
	result := db.Order("number desc").Last(&l2ScannedBlock)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			db.Create(&l2ScannedBlock)
		} else {
			return nil, fmt.Errorf("failed to query l2_scanned_blocks: %w", result.Error)
		}
	} else {
		if l2ScannedBlock.Number < cfg.Misc.ConfirmBlocks {
			l2ScannedBlock.Number = 0
		} else {
			l2ScannedBlock.Number -= cfg.Misc.ConfirmBlocks
		}
	}
	return &l2ScannedBlock, nil
}
