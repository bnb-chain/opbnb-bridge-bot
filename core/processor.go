package core

import (
	bindings2 "bnbchain/opbnb-bridge-bot/bindings"
	"context"
	"errors"
	"fmt"
	"math/big"

	"github.com/ethereum-optimism/optimism/indexer/config"
	"github.com/ethereum-optimism/optimism/op-bindings/bindings"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/log"
)

type Processor struct {
	log log.Logger

	L1Client *ClientExt
	L2Client *ClientExt

	cfg         Config
	L2Contracts config.L2Contracts

	whitelistL2TokenMap map[common.Address]struct{}
}

func NewProcessor(
	log log.Logger,
	l1Client *ClientExt,
	l2Client *ClientExt,
	cfg Config,
) *Processor {
	l2Contracts := config.L2ContractsFromPredeploys()

	var whitelistL2TokenMap map[common.Address]struct{} = nil
	if cfg.L2StandardBridgeBot.WhitelistL2TokenList != nil {
		whitelistL2TokenMap = make(map[common.Address]struct{})
		for _, l2Token := range *cfg.L2StandardBridgeBot.WhitelistL2TokenList {
			whitelistL2TokenMap[common.HexToAddress(l2Token)] = struct{}{}
		}
	}

	return &Processor{log, l1Client, l2Client, cfg, l2Contracts, whitelistL2TokenMap}
}

func (b *Processor) toWithdrawal(botDelegatedWithdrawToEvent *L2ContractEvent, receipt *types.Receipt) (*bindings.TypesWithdrawalTransaction, error) {
	// Events flow:
	//
	// event[i-5]: WithdrawalInitiated
	// event[i-4]: ETHBridgeInitiated
	// event[i-3]: MessagePassed
	// event[i-2]: SentMessage
	// event[i-1]: SentMessageExtension1
	// event[i]  : L2StandardBridgeBot.WithdrawTo
	if botDelegatedWithdrawToEvent.LogIndex < 5 || len(receipt.Logs) < 6 {
		return nil, fmt.Errorf("invalid botDelegatedWithdrawToEvent: %v", botDelegatedWithdrawToEvent)
	}

	messagePassedLog := receipt.Logs[botDelegatedWithdrawToEvent.LogIndex-3]
	sentMessageLog := receipt.Logs[botDelegatedWithdrawToEvent.LogIndex-2]
	sentMessageExtension1Log := receipt.Logs[botDelegatedWithdrawToEvent.LogIndex-1]

	sentMessageEvent, err := b.toL2CrossDomainMessengerSentMessageExtension1(sentMessageLog, sentMessageExtension1Log)
	if err != nil {
		return nil, err
	}
	messagePassedEvent, err := b.toMessagePassed(messagePassedLog)
	if err != nil {
		return nil, err
	}

	withdrawalTx, err := b.toLowLevelMessage(sentMessageEvent, messagePassedEvent)
	if err != nil {
		return nil, fmt.Errorf("toLowLevelMessage err: %v", err)
	}

	return withdrawalTx, nil
}

func (b *Processor) ProveWithdrawalTransaction(ctx context.Context, botDelegatedWithdrawToEvent *L2ContractEvent) error {
	receipt, err := b.L2Client.TransactionReceipt(ctx, common.HexToHash(botDelegatedWithdrawToEvent.TransactionHash))
	if err != nil {
		return err
	}

	err = b.CheckByFilterOptions(botDelegatedWithdrawToEvent, receipt)
	if err != nil {
		return err
	}

	l2BlockNumber := receipt.BlockNumber
	withdrawalTx, err := b.toWithdrawal(botDelegatedWithdrawToEvent, receipt)
	if err != nil {
		return fmt.Errorf("toWithdrawal err: %v", err)
	}

	hash, err := b.hashWithdrawal(withdrawalTx)
	if err != nil {
		return fmt.Errorf("hashWithdrawal err: %v", err)
	}

	messageSlot, err := b.hashMessageHash(hash)
	if err != nil {
		return fmt.Errorf("hashMesaageHash err: %v", err)
	}

	l2OutputIndex, l2OutputProposal, err := b.getLatestL2OutputProposal()
	if err != nil {
		return err
	}
	if l2OutputProposal.L2BlockNumber.Uint64() < l2BlockNumber.Uint64() {
		return errors.New("L2OutputOracle: cannot get output for a block that has not been proposed")
	}

	accountResult, err := b.L2Client.GetProof(
		b.L2Contracts.L2ToL1MessagePasser,
		[]string{"0x" + messageSlot},
		l2OutputProposal.L2BlockNumber,
	)
	if err != nil {
		return fmt.Errorf("GetProof err: %v", err)
	}

	outputProposalBlock, err := b.L2Client.HeaderByNumber(ctx, l2OutputProposal.L2BlockNumber)
	if err != nil {
		return fmt.Errorf("get output proposal block error: %v", err)
	}

	if len(accountResult.StorageProof) == 0 {
		return fmt.Errorf("no storage proof")
	}

	withdrawalProof := accountResult.StorageProof[0]
	withdrawalProof2Bytes := make([][]byte, 0)
	for _, p1 := range withdrawalProof.Proof {
		withdrawalProof2Bytes = append(withdrawalProof2Bytes, p1)
	}

	outputRootProof := bindings.TypesOutputRootProof{
		Version:                  common.HexToHash("0x"),
		StateRoot:                outputProposalBlock.Root,
		MessagePasserStorageRoot: accountResult.StorageHash,
		LatestBlockhash:          outputProposalBlock.Hash(),
	}

	l1ChainId, err := b.L1Client.ChainID(context.Background())
	if err != nil {
		return err
	}

	gasPrice := big.NewInt(b.cfg.TxSigner.GasPrice)
	signerPrivkey, signerAddress, err := b.cfg.SignerKeyPair()
	if err != nil {
		return err
	}

	optimismPortalTransactor, _ := bindings.NewOptimismPortalTransactor(
		b.cfg.L1Contracts.OptimismPortalProxy,
		b.L1Client,
	)
	signedTx, err := optimismPortalTransactor.ProveWithdrawalTransaction(
		&bind.TransactOpts{
			From:     *signerAddress,
			GasPrice: gasPrice,
			Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
				return types.SignTx(tx, types.NewEIP155Signer(l1ChainId), signerPrivkey)
			},
		},
		*withdrawalTx,
		l2OutputIndex,
		outputRootProof,
		withdrawalProof2Bytes,
	)
	if err != nil {
		return err
	}

	b.log.Info("ProveWithdrawalTransaction", "tx_hash", signedTx.Hash().Hex())
	return nil
}

// FinalizeMessage https://github.com/ethereum-optimism/optimism/blob/d90e7818de894f0bc93ae7b449b9049416bda370/packages/sdk/src/cross-chain-messenger.ts#L1611
func (b *Processor) FinalizeMessage(ctx context.Context, botDelegatedWithdrawToEvent *L2ContractEvent) error {
	receipt, err := b.L2Client.TransactionReceipt(ctx, common.HexToHash(botDelegatedWithdrawToEvent.TransactionHash))
	if err != nil {
		return err
	}

	err = b.CheckByFilterOptions(botDelegatedWithdrawToEvent, receipt)
	if err != nil {
		return err
	}

	withdrawalTx, err := b.toWithdrawal(botDelegatedWithdrawToEvent, receipt)
	if err != nil {
		return fmt.Errorf("toWithdrawal err: %v", err)
	}

	l1ChainId, err := b.L1Client.ChainID(ctx)
	if err != nil {
		return err
	}

	gasPrice := big.NewInt(b.cfg.TxSigner.GasPrice)
	signerPrivkey, signerAddress, err := b.cfg.SignerKeyPair()
	if err != nil {
		return err
	}

	optimismPortalTransactor, _ := bindings.NewOptimismPortalTransactor(
		b.cfg.L1Contracts.OptimismPortalProxy,
		b.L1Client,
	)
	signedTx, err := optimismPortalTransactor.FinalizeWithdrawalTransaction(
		&bind.TransactOpts{
			From:     *signerAddress,
			GasPrice: gasPrice,
			Signer: func(address common.Address, tx *types.Transaction) (*types.Transaction, error) {
				return types.SignTx(tx, types.NewEIP155Signer(l1ChainId), signerPrivkey)
			},
		},
		*withdrawalTx,
	)
	if err != nil {
		return fmt.Errorf("optimismPortalTransactor.FinalizeWithdrawalTransaction: %w", err)
	}

	b.log.Info("FinalizeWithdrawalTransaction", "tx_hash", signedTx.Hash().Hex())
	return nil
}

func (b *Processor) hashWithdrawal(w *bindings.TypesWithdrawalTransaction) (string, error) {
	uint256Type, _ := abi.NewType("uint256", "", nil)
	addressType, _ := abi.NewType("address", "", nil)
	bytesType, _ := abi.NewType("bytes", "", nil)
	types_ := abi.Arguments{
		{Type: uint256Type},
		{Type: addressType},
		{Type: addressType},
		{Type: uint256Type},
		{Type: uint256Type},
		{Type: bytesType},
	}
	encoded, err := types_.Pack(w.Nonce, w.Sender, w.Target, w.Value, w.GasLimit, w.Data)
	if err != nil {
		return "", fmt.Errorf("pack hashWithdrawal arguments: %w", err)
	}
	result := crypto.Keccak256(encoded)
	return common.Bytes2Hex(result), nil
}

func (b *Processor) hashMessageHash(messageHash string) (string, error) {
	uint256Type, _ := abi.NewType("uint256", "", nil)
	bytes32Type, _ := abi.NewType("bytes32", "", nil)
	types_ := abi.Arguments{
		{
			Type: bytes32Type,
		},
		{
			Type: uint256Type,
		},
	}

	encoded, err := types_.Pack(common.HexToHash(messageHash), big.NewInt(0))
	if err != nil {
		return "", err
	}

	return common.Bytes2Hex(crypto.Keccak256(encoded)), nil
}

type L2CrossDomainMessengerSentMessageExtension1 struct {
	bindings.L2CrossDomainMessengerSentMessage
	Value *big.Int
}

func (b *Processor) toL2CrossDomainMessengerSentMessageExtension1(sentMessageLog, sentMessageExtension1Log *types.Log) (*L2CrossDomainMessengerSentMessageExtension1, error) {
	addressType, _ := abi.NewType("address", "", nil)
	L2CrossDomainMessengerAbi, _ := bindings.L2CrossDomainMessengerMetaData.GetAbi()

	if !(sentMessageLog.Address == b.L2Contracts.L2CrossDomainMessenger &&
		len(sentMessageLog.Topics) > 1 &&
		sentMessageLog.Topics[0] == L2CrossDomainMessengerAbi.Events["SentMessage"].ID) {
		return nil, errors.New("invalid log: not SentMessage event")
	}

	sentMessageEvent := bindings.L2CrossDomainMessengerSentMessage{}
	err := abi.ParseTopics(
		&sentMessageEvent,
		[]abi.Argument{
			{
				Name:    "target",
				Type:    addressType,
				Indexed: true,
			},
		},
		sentMessageLog.Topics[1:],
	)
	if err != nil {
		return nil, fmt.Errorf("parse indexed event arguments from log.topics of SentMessage event, err: %v", err)
	}

	// NOTE: log.Data only contains the non-indexed arguments
	err = L2CrossDomainMessengerAbi.UnpackIntoInterface(&sentMessageEvent, "SentMessage", sentMessageLog.Data)
	if err != nil {
		return nil, fmt.Errorf("parse non-indexed event arguments from log.data of SentMessage event, err: %v", err)
	}

	if !(sentMessageExtension1Log.Address == b.L2Contracts.L2CrossDomainMessenger &&
		len(sentMessageExtension1Log.Topics) > 1 &&
		sentMessageExtension1Log.Topics[0] == L2CrossDomainMessengerAbi.Events["SentMessageExtension1"].ID) {
		return nil, errors.New("invalid log: not SentMessageExtension1 event")
	}

	sentMessageExtension1 := bindings.L2CrossDomainMessengerSentMessageExtension1{}
	err = L2CrossDomainMessengerAbi.UnpackIntoInterface(&sentMessageExtension1, "SentMessageExtension1", sentMessageExtension1Log.Data)
	if err != nil {
		return nil, fmt.Errorf("UnpackIntoInterface SentMessageExtension1: %w", err)
	}

	return &L2CrossDomainMessengerSentMessageExtension1{
		L2CrossDomainMessengerSentMessage: sentMessageEvent,
		Value:                             sentMessageExtension1.Value,
	}, nil
}

func (b *Processor) toMessagePassed(messagePassedLog *types.Log) (*bindings.L2ToL1MessagePasserMessagePassed, error) {
	uint256Type, _ := abi.NewType("uint256", "", nil)
	addressType, _ := abi.NewType("address", "", nil)
	L2ToL1MessagePasserAbi, _ := bindings.L2ToL1MessagePasserMetaData.GetAbi()

	if messagePassedLog.Address == b.L2Contracts.L2ToL1MessagePasser &&
		len(messagePassedLog.Topics) > 0 &&
		messagePassedLog.Topics[0] == L2ToL1MessagePasserAbi.Events["MessagePassed"].ID {
	}

	messagePassedEvent := bindings.L2ToL1MessagePasserMessagePassed{}
	err := abi.ParseTopics(
		&messagePassedEvent,
		[]abi.Argument{
			{Name: "nonce", Type: uint256Type, Indexed: true},
			{Name: "sender", Type: addressType, Indexed: true},
			{Name: "target", Type: addressType, Indexed: true},
		},
		messagePassedLog.Topics[1:],
	)
	if err != nil {
		return nil, fmt.Errorf("parse indexed event arguments from log.topics of MessagePassed event, err: %v", err)
	}

	// NOTE: log.Data only contains the non-indexed arguments
	err = L2ToL1MessagePasserAbi.UnpackIntoInterface(&messagePassedEvent, "MessagePassed", messagePassedLog.Data)
	if err != nil {
		return nil, fmt.Errorf("parse non-indexed event arguments from log.data of SentMessage event, err: %v", err)
	}

	// NOTE: log.Data only contains the non-indexed arguments
	err = L2ToL1MessagePasserAbi.UnpackIntoInterface(&messagePassedEvent, "MessagePassed", messagePassedLog.Data)
	if err != nil {
		return nil, fmt.Errorf("parse non-indexed event arguments from log.data of MessagePassed event, err: %v", err)
	}

	return &messagePassedEvent, nil
}

// getSentMessagesByReceipt retrieves all cross chain messages sent within a given transaction.
func (b *Processor) getSentMessagesByReceipt(receipt *types.Receipt) ([]L2CrossDomainMessengerSentMessageExtension1, error) {
	L2CrossDomainMessengerAbi, _ := bindings.L2CrossDomainMessengerMetaData.GetAbi()
	addressType, _ := abi.NewType("address", "", nil)

	// Filter SentMessage(address indexed target, address sender, bytes message, uint256 messageNonce, uint256 gasLimit)
	sentMessageEvents := make([]L2CrossDomainMessengerSentMessageExtension1, 0)
	for i, l := range receipt.Logs {
		if l.Address == b.L2Contracts.L2CrossDomainMessenger &&
			len(l.Topics) > 0 &&
			l.Topics[0] == L2CrossDomainMessengerAbi.Events["SentMessage"].ID {

			sentMessageEvent := bindings.L2CrossDomainMessengerSentMessage{}
			err := abi.ParseTopics(
				&sentMessageEvent,
				[]abi.Argument{
					{
						Name:    "target",
						Type:    addressType,
						Indexed: true,
					},
				},
				l.Topics[1:],
			)
			if err != nil {
				return nil, fmt.Errorf("parse indexed event arguments from log.topics of SentMessage event, err: %v", err)
			}

			// NOTE: log.Data only contains the non-indexed arguments
			err = L2CrossDomainMessengerAbi.UnpackIntoInterface(&sentMessageEvent, "SentMessage", l.Data)
			if err != nil {
				return nil, fmt.Errorf("parse non-indexed event arguments from log.data of SentMessage event, err: %v", err)
			}

			if i+1 < len(receipt.Logs) &&
				receipt.Logs[i+1].Address == b.L2Contracts.L2CrossDomainMessenger &&
				len(receipt.Logs[i+1].Topics) > 1 &&
				receipt.Logs[i+1].Topics[0] == L2CrossDomainMessengerAbi.Events["SentMessageExtension1"].ID {

				sentMessageExtension1 := bindings.L2CrossDomainMessengerSentMessageExtension1{}
				err := L2CrossDomainMessengerAbi.UnpackIntoInterface(&sentMessageExtension1, "SentMessageExtension1", receipt.Logs[i+1].Data)
				if err != nil {
					return nil, fmt.Errorf("UnpackIntoInterface SentMessageExtension1: %w", err)
				}

				sentMessageEvents = append(sentMessageEvents, L2CrossDomainMessengerSentMessageExtension1{
					L2CrossDomainMessengerSentMessage: sentMessageEvent,
					Value:                             sentMessageExtension1.Value,
				})
			}
		}
	}

	return sentMessageEvents, nil
}

func (b *Processor) getMessagePassedMessagesFromReceipt(receipt *types.Receipt) ([]bindings.L2ToL1MessagePasserMessagePassed, error) {
	L2ToL1MessagePasserAbi, _ := bindings.L2ToL1MessagePasserMetaData.GetAbi()
	uint256Type, _ := abi.NewType("uint256", "", nil)
	addressType, _ := abi.NewType("address", "", nil)

	messagePassedLogs := make([]*types.Log, 0)
	for _, l := range receipt.Logs {
		if l.Address == b.L2Contracts.L2ToL1MessagePasser &&
			len(l.Topics) > 0 &&
			l.Topics[0] == L2ToL1MessagePasserAbi.Events["MessagePassed"].ID {
			messagePassedLogs = append(messagePassedLogs, l)
		}
	}
	if len(messagePassedLogs) == 0 {
		return nil, errors.New("no MessagePassed event")
	}

	// Parse SentMessage events
	messagePassedEvents := make([]bindings.L2ToL1MessagePasserMessagePassed, len(messagePassedLogs))
	for i, l := range messagePassedLogs {
		messagePassedEvent := bindings.L2ToL1MessagePasserMessagePassed{}
		err := abi.ParseTopics(
			&messagePassedEvent,
			[]abi.Argument{
				{Name: "nonce", Type: uint256Type, Indexed: true},
				{Name: "sender", Type: addressType, Indexed: true},
				{Name: "target", Type: addressType, Indexed: true},
			},
			l.Topics[1:],
		)
		if err != nil {
			return nil, fmt.Errorf("parse indexed event arguments from log.topics of MessagePassed event, err: %v", err)
		}

		// NOTE: log.Data only contains the non-indexed arguments
		err = L2ToL1MessagePasserAbi.UnpackIntoInterface(&messagePassedEvent, "MessagePassed", l.Data)
		if err != nil {
			return nil, fmt.Errorf("parse non-indexed event arguments from log.data of SentMessage event, err: %v", err)
		}

		// NOTE: log.Data only contains the non-indexed arguments
		err = L2ToL1MessagePasserAbi.UnpackIntoInterface(&messagePassedEvent, "MessagePassed", l.Data)
		if err != nil {
			return nil, fmt.Errorf("parse non-indexed event arguments from log.data of MessagePassed event, err: %v", err)
		}

		messagePassedEvents[i] = messagePassedEvent
	}

	return messagePassedEvents, nil
}

func (b *Processor) getLatestL2OutputProposal() (*big.Int, *bindings.TypesOutputProposal, error) {
	l2OutputOracleCaller, err := bindings.NewL2OutputOracleCaller(
		b.cfg.L1Contracts.L2OutputOracleProxy,
		b.L1Client,
	)
	if err != nil {
		return nil, nil, fmt.Errorf("NewL2OutputOracleCaller err: %v", err)
	}

	// [getBedrockMessageProof](https://github.com/ethereum-optimism/optimism/blob/d90e7818de894f0bc93ae7b449b9049416bda370/packages/sdk/src/cross-chain-messenger.ts#L1916)
	l2OutputIndex, err := l2OutputOracleCaller.LatestOutputIndex(&bind.CallOpts{})
	if err != nil {
		return nil, nil, fmt.Errorf("GetL2OutputIndexAfter err: %v", err)
	}

	outputProposal, err := l2OutputOracleCaller.GetL2Output(&bind.CallOpts{}, l2OutputIndex)
	if err != nil {
		return nil, nil, fmt.Errorf("GetL2Output err: %v", err)
	}

	return l2OutputIndex, &outputProposal, nil
}

func (b *Processor) toLowLevelMessage(
	sentMessageEvent *L2CrossDomainMessengerSentMessageExtension1,
	messagePassedEvent *bindings.L2ToL1MessagePasserMessagePassed,
) (*bindings.TypesWithdrawalTransaction, error) {
	// Encode "relayMessage" with signature, the result will be attached to [WithdrawalTransaction.Data](https://github.com/ethereum-optimism/optimism/blob/f54a2234f2f350795552011f35f704a3feb56a08/packages/contracts-bedrock/src/libraries/Types.sol#L68)
	L2CrossDomainMessengerAbi, _ := bindings.L2CrossDomainMessengerMetaData.GetAbi()
	relayMessageCalldata, err := L2CrossDomainMessengerAbi.Pack(
		"relayMessage",
		sentMessageEvent.MessageNonce,
		sentMessageEvent.Sender,
		sentMessageEvent.Target,
		sentMessageEvent.Value,
		sentMessageEvent.GasLimit,
		sentMessageEvent.Message,
	)
	if err != nil {
		return nil, fmt.Errorf("encode relayMessage calldata, err: %v", err)
	}

	withdrawalTx := bindings.TypesWithdrawalTransaction{
		Nonce:    messagePassedEvent.Nonce,
		Sender:   b.L2Contracts.L2CrossDomainMessenger,
		Target:   b.cfg.L1Contracts.L1CrossDomainMessengerProxy,
		Value:    sentMessageEvent.Value,
		GasLimit: messagePassedEvent.GasLimit,
		Data:     relayMessageCalldata,
	}
	return &withdrawalTx, nil
}

func (b *Processor) CheckByFilterOptions(botDelegatedWithdrawToEvent *L2ContractEvent, receipt *types.Receipt) error {
	L2StandardBridgeBotAbi, _ := bindings2.L2StandardBridgeBotMetaData.GetAbi()
	withdrawToEvent := bindings2.L2StandardBridgeBotWithdrawTo{}
	indexedArgs := func(arguments abi.Arguments) abi.Arguments {
		indexedArgs := abi.Arguments{}
		for _, arg := range arguments {
			if arg.Indexed {
				indexedArgs = append(indexedArgs, arg)
			}
		}
		return indexedArgs
	}
	err := abi.ParseTopics(&withdrawToEvent, indexedArgs(L2StandardBridgeBotAbi.Events["WithdrawTo"].Inputs), receipt.Logs[botDelegatedWithdrawToEvent.LogIndex].Topics[1:])
	if err != nil {
		return fmt.Errorf("parse indexed event arguments from log.topics of L2StandardBridgeBotWithdrawTo event, err: %v", err)
	}

	err = L2StandardBridgeBotAbi.UnpackIntoInterface(&withdrawToEvent, "WithdrawTo", receipt.Logs[botDelegatedWithdrawToEvent.LogIndex].Data)
	if err != nil {
		return fmt.Errorf("parse non-indexed event arguments from log.data of L2StandardBridgeBotWithdrawTo event, err: %v", err)
	}

	if !IsL2TokenWhitelisted(b.whitelistL2TokenMap, &withdrawToEvent.L2Token) {
		return fmt.Errorf("filtered: token is not whitelisted, l2-token: %s", withdrawToEvent.L2Token)
	}
	if !IsMinGasLimitValid(b.cfg.L2StandardBridgeBot.UpperMinGasLimit, withdrawToEvent.MinGasLimit) {
		return fmt.Errorf("filtered: minGasLimit is too large, minGasLimit: %d", withdrawToEvent.MinGasLimit)
	}
	if !IsExtraDataValid(b.cfg.L2StandardBridgeBot.UpperMinGasLimit, &withdrawToEvent.ExtraData) {
		return fmt.Errorf("filtered: extraData is too large, extraDataSize: %d", len(withdrawToEvent.ExtraData))
	}

	return nil
}

func IsL2TokenWhitelisted(whitelistL2TokenMap map[common.Address]struct{}, l2Token *common.Address) bool {
	// nil means all L2 tokens are whitelisted
	if whitelistL2TokenMap == nil {
		return true
	}

	_, exists := whitelistL2TokenMap[*l2Token]
	return exists
}

func IsMinGasLimitValid(upperMinGasLimit *uint32, minGasLimit uint32) bool {
	// nil means no limit
	if upperMinGasLimit == nil {
		return true
	}

	return minGasLimit <= *upperMinGasLimit
}

func IsExtraDataValid(upperExtraDataSize *uint32, extraData *[]byte) bool {
	// nil means no limit
	if upperExtraDataSize == nil {
		return true
	}

	return len(*extraData) <= int(*upperExtraDataSize)
}
