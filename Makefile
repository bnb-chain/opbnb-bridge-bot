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

.PHONY: \
	bot  \
	build-go \
	build-solidity
