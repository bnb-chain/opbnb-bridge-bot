GITCOMMIT := $(shell git rev-parse HEAD)
GITDATE := $(shell git show -s --format='%ct')

LDFLAGSSTRING +=-X main.GitCommit=$(GITCOMMIT)
LDFLAGSSTRING +=-X main.GitDate=$(GITDATE)
LDFLAGS := -ldflags "$(LDFLAGSSTRING)"

build-go:
	env GO111MODULE=on go build -v $(LDFLAGS) ./cmd/bot

build-solidity:
	pushd contracts; \
	forge build; \
	popd;

bindings: build-solidity
	jq '.abi' contracts/out/L2StandardBridgeBot.sol/L2StandardBridgeBot.json  > contracts/out/L2StandardBridgeBot.sol/L2StandardBridgeBot.abi; \
	abigen --abi contracts/out/L2StandardBridgeBot.sol/L2StandardBridgeBot.abi --pkg bindings --type L2StandardBridgeBot --out bindings/L2StandardBridgeBot.go

.PHONY: \
	bot  \
	build-go \
	build-solidity
