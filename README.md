# opbnb-bridge-bot

## Description

The opbnb-bridge-bot project primarily maintains a Layer2 contract named `L2StandardBridgeBot`. This contract, when called, initiates withdrawal transactions and collects a fixed fee for every withdrawals. The project promises that after the withdrawal transaction has can be proven and finalized after the required time window, a third-party account would complete the corresponding L1 proven withdrawal transactions and L1 finalized withdrawal transactions, thus completing the entire withdrawal process.

## Design Principle and Working Mechanism

This project consists of an on-chain contract `contracts/src/L2StandardBridgeBot.sol` and an off-chain bot.

The `L2StandardBridgeBot` contract provides a `withdrawTo` function, which charges a fixed fee for every execution and then forwards the call to the `L2StandardBridge.withdraw`.

The off-chain bot watches the `L2StandardBridgeBot.WithdrawTo` events and based on these events, re-constructs the corresponding withdrawals. We name these withdrawals as **bot-delegated withdrawals**. As time go out of the bot-delegated withdrawal's proven and finalized time window, our bot will send proven and finalized transactions to complete the entire withdrawal process.

## User Guide

### Getting Started at opBNB testnet

1. Prepare a PostgreSQL database

```
docker-compose up -d
```

2. Compile the off-chain bot and output the artifact to `./bot`

```
make build-go
```

3. Run the off-chain bot

```
OPBNB_BRIDGE_BOT_PRIVKEY=<bot privkey> ./bot --config ./bot.toml
```

### Deploy and Use Contracts

1. Compile the contract using `forge`

```
make build-solidity
```

2. Deploy contract

```
cd contracts

export DELEGATE_FEE=1000000000000000
forge create \
    --rpc-url $OPBNB_TESTNET \
    --private-key $OP_DEPLOYER_PRIVKEY \
    src/L2StandardBridgeBot.sol:L2StandardBridgeBot --constructor-args $OP_DEPLOYER_ADDRESS $DELEGATE_FEE
```

3. Withdraw via the deployed contract

```
export DELEGATE_FEE=1000000000000000
export amount=2000000000000001
export contract_addr=<deployed contract address>

cast send --rpc-url $OPBNB_TESTNET \
          --private-key $OP_DEPLOYER_PRIVKEY \
          --value $(($DELEGATE_FEE + $amount)) \
          $contract_addr \
          $(cast calldata 'withdrawTo( address _l2Token, address _to, uint256 _amount, uint32 _minGasLimit, bytes calldata _extraData)' 0xDeadDeAddeAddEAddeadDEaDDEAdDeaDDeAD0000 $OP_DEPLOYER_ADDRESS $amount 150469 "")
```

## License

MIT
