package core

import (
	bindings2 "bnbchain/opbnb-bridge-bot/bindings"
	"context"
	"math/big"
	"time"

	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var L2FeeVaults = []common.Address{
	common.HexToAddress("0x4200000000000000000000000000000000000011"), // SequencerFeeVault
	common.HexToAddress("0x4200000000000000000000000000000000000019"), // BaseFeeVault
	common.HexToAddress("0x420000000000000000000000000000000000001a"), // L1FeeVault
}

type Indexer struct {
	log      log.Logger
	db       *gorm.DB
	l2Client *ClientExt
	cfg      Config

	contracts []common.Address

	// Note: feeVaults indicates the configured list of FeeVault contracts:
	// - [SequencerFeeVault.sol](https://github.com/bnb-chain/opbnb/blob/develop/packages/contracts-bedrock/contracts/L2/SequencerFeeVault.sol)
	// - [BaseFeeVault.sol](https://github.com/bnb-chain/opbnb/blob/develop/packages/contracts-bedrock/contracts/L2/BaseFeeVault.sol)
	// - [L1FeeVault.sol](https://github.com/bnb-chain/opbnb/blob/develop/packages/contracts-bedrock/contracts/L2/L1FeeVault.sol)
	isFeeVaultWithdrawEvent                  func(vlog *types.Log) bool
	isL2StandardBridgeBotWithdrawToEvent     func(vlog *types.Log) bool
	isL2StandardBridgeWithdrawalInitiatedLog func(vlog *types.Log) bool
}

func NewIndexer(log log.Logger, db *gorm.DB, l2Client *ClientExt, cfg Config) *Indexer {
	l2StandardBridgeBots := make(map[common.Address]struct{})
	feeVaults := make(map[common.Address]struct{})
	contracts := make([]common.Address, 0)
	addr := cfg.L2StandardBridgeBot.ContractAddress
	addr_ := common.HexToAddress(addr)
	l2StandardBridgeBots[addr_] = struct{}{}
	contracts = append(contracts, addr_)

	for _, addr := range L2FeeVaults {
		feeVaults[addr] = struct{}{}
		contracts = append(contracts, addr)
	}

	isL2StandardBridgeWithdrawalInitiatedLog := func(vlog *types.Log) bool {
		var (
			L2StandardBridgeAbi, _                 = bindings.L2StandardBridgeMetaData.GetAbi()
			L2StandardBridgeWithdrawalInitiatedSig = L2StandardBridgeAbi.Events["WithdrawalInitiated"].ID
		)
		return len(vlog.Topics) > 1 && vlog.Topics[0] == L2StandardBridgeWithdrawalInitiatedSig
	}
	isL2StandardBridgeBotWithdrawToEvent := func(log *types.Log) bool {
		var (
			L2StandardBridgeBotAbi, _        = bindings2.L2StandardBridgeBotMetaData.GetAbi()
			L2StandardBridgeBotWithdrawToSig = L2StandardBridgeBotAbi.Events["WithdrawTo"].ID
		)
		if len(log.Topics) > 0 && log.Topics[0] == L2StandardBridgeBotWithdrawToSig {
			_, ok := l2StandardBridgeBots[log.Address]
			return ok
		}
		return false
	}
	isFeeVaultWithdrawEvent := func(log *types.Log) bool {
		var (
			FeeVaultAbi, _        = bindings.L1FeeVaultMetaData.GetAbi()
			FeeVaultWithdrawalSig = FeeVaultAbi.Events["Withdrawal"].ID
		)
		if len(log.Topics) > 0 && log.Topics[0] == FeeVaultWithdrawalSig {
			_, ok := feeVaults[log.Address]
			return ok
		}
		return false
	}

	return &Indexer{
		log:                                      log,
		db:                                       db,
		l2Client:                                 l2Client,
		cfg:                                      cfg,
		isL2StandardBridgeWithdrawalInitiatedLog: isL2StandardBridgeWithdrawalInitiatedLog,
		isL2StandardBridgeBotWithdrawToEvent:     isL2StandardBridgeBotWithdrawToEvent,
		isFeeVaultWithdrawEvent:                  isFeeVaultWithdrawEvent,
		contracts:                                contracts,
	}
}

// Start watches for new bot-delegated withdrawals and stores them in the database.
func (i *Indexer) Start(ctx context.Context, l2ScannedBlock *L2ScannedBlock) {
	timer := time.NewTimer(0)
	fromBlockNumber := big.NewInt(l2ScannedBlock.Number)

	for {
		select {
		case <-ctx.Done():
			return
		case <-timer.C:
			timer.Reset(time.Second)
		}

		toBlockNumber := new(big.Int).Add(fromBlockNumber, big.NewInt(i.cfg.L2StandardBridgeBot.LogFilterBlockRange))
		finalizedHeader, err := i.l2Client.GetHeaderByTag(context.Background(), "finalized")
		if err != nil {
			log.Error("call eth_blockNumber", "error", err)
			continue
		}
		if toBlockNumber.Uint64() > finalizedHeader.Number.Uint64() {
			toBlockNumber = finalizedHeader.Number
		}

		if fromBlockNumber.Uint64() > toBlockNumber.Uint64() {
			timer.Reset(5 * time.Second)
			continue
		}

		logs, err := i.getWithdrawalInitiatedLogs(ctx, fromBlockNumber, toBlockNumber)
		if err != nil {
			log.Error("eth_getLogs", "error", err)
			continue
		}

		if len(logs) != 0 {
			for _, vlog := range logs {
				log.Info("fetched bot-delegated withdrawal", "blockNumber", vlog.BlockNumber, "transactionHash", vlog.TxHash.Hex())
			}

			err = i.storeLogs(logs)
			if err != nil {
				log.Error("storeLogs", "error", err)
				continue
			}
		}

		l2ScannedBlock.Number = toBlockNumber.Int64()
		result := i.db.Save(l2ScannedBlock)
		if result.Error != nil {
			log.Error("update l2_scanned_blocks", "error", result.Error)
		}

		fromBlockNumber = new(big.Int).Add(toBlockNumber, big.NewInt(1))
	}
}

func (i *Indexer) getWithdrawalInitiatedLogs(ctx context.Context, fromBlock *big.Int, toBlock *big.Int) ([]types.Log, error) {
	i.log.Debug("Fetching logs from blocks", "fromBlock", fromBlock, "toBlock", toBlock)

	var (
		L2StandardBridgeBotAbi, _        = bindings2.L2StandardBridgeBotMetaData.GetAbi()
		L2StandardBridgeBotWithdrawToSig = L2StandardBridgeBotAbi.Events["WithdrawTo"].ID

		FeeVaultAbi, _        = bindings.L1FeeVaultMetaData.GetAbi()
		FeeVaultWithdrawalSig = FeeVaultAbi.Events["Withdrawal"].ID
	)

	logs, err := i.l2Client.FilterLogs(context.Background(), ethereum.FilterQuery{
		FromBlock: fromBlock,
		ToBlock:   toBlock,
		Addresses: i.contracts,
		Topics:    [][]common.Hash{[]common.Hash{L2StandardBridgeBotWithdrawToSig, FeeVaultWithdrawalSig}},
	})
	if err != nil {
		return nil, err
	}

	withdrawalInitiatedLogs := make([]types.Log, 0)
	for _, vlog := range logs {
		receipt, err := i.l2Client.TransactionReceipt(ctx, vlog.TxHash)
		if err != nil {
			return nil, err
		}

		if i.isL2StandardBridgeBotWithdrawToEvent(&vlog) {
			// Events flow:
			//
			// event[i-5]: WithdrawalInitiated
			// event[i-4]: ETHBridgeInitiated
			// event[i-3]: MessagePassed
			// event[i-2]: SentMessage
			// event[i-1]: SentMessageExtension1
			// event[i]  : L2StandardBridgeBot.WithdrawTo
			withdrawalInitiatedLog := GetLogByLogIndex(receipt, vlog.Index-5)
			if withdrawalInitiatedLog != nil && i.isL2StandardBridgeWithdrawalInitiatedLog(withdrawalInitiatedLog) {
				withdrawalInitiatedLogs = append(withdrawalInitiatedLogs, *withdrawalInitiatedLog)
			} else {
				i.log.Crit("eth_getLogs returned an unexpected event", "L2StandardBridgeBotWithdrawLog", vlog, "WithdrawalInitiatedLog", withdrawalInitiatedLog)
			}
		} else if i.isFeeVaultWithdrawEvent(&vlog) {
			// Events flow:
			//
			// event[i]  : FeeVault.Withdrawal
			// event[i+1]: WithdrawalInitiated
			// event[i+2]: ETHBridgeInitiated
			// event[i+3]: MessagePassed
			// event[i+4]: SentMessage
			// event[i+5]: SentMessageExtension1
			withdrawalInitiatedLog := GetLogByLogIndex(receipt, vlog.Index+1)
			if withdrawalInitiatedLog != nil && i.isL2StandardBridgeWithdrawalInitiatedLog(withdrawalInitiatedLog) {
				withdrawalInitiatedLogs = append(withdrawalInitiatedLogs, *withdrawalInitiatedLog)
			} else {
				i.log.Crit("eth_getLogs returned an unexpected event", "FeeVaultWithdrawLog", vlog, "WithdrawalInitiatedLog", withdrawalInitiatedLog)
			}
		} else {
			i.log.Crit("eth_getLogs returned an unexpected event", "log", vlog)
		}
	}

	return withdrawalInitiatedLogs, nil
}

// storeLogs stores the logs in the database
func (i *Indexer) storeLogs(logs []types.Log) error {
	// save all the logs in this range of blocks
	for _, vLog := range logs {
		header, err := i.l2Client.HeaderByHash(context.Background(), vLog.BlockHash)
		if err != nil {
			return err
		}

		deduped := i.db.Clauses(clause.OnConflict{DoNothing: true})
		result := deduped.Create(&WithdrawalInitiatedLog{
			TransactionHash:      vLog.TxHash.Hex(),
			LogIndex:             int(vLog.Index),
			InitiatedBlockNumber: int64(header.Number.Uint64()),
		})
		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}
