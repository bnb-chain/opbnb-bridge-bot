package core_test

import (
	"bnbchain/opbnb-bridge-bot/core"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
)

func TestIsL2TokenWhitelisted(t *testing.T) {
	// Test case 1: whitelistL2TokenMap is nil
	var whitelistL2TokenMap map[common.Address]struct{} = nil
	l2Token := common.HexToAddress("0x1234abcdef")
	assert.True(t, core.IsL2TokenWhitelisted(whitelistL2TokenMap, &l2Token))

	// Test case 2: l2Token is whitelisted
	whitelistL2TokenMap = make(map[common.Address]struct{})
	whitelistL2TokenMap[l2Token] = struct{}{}
	assert.True(t, core.IsL2TokenWhitelisted(whitelistL2TokenMap, &l2Token))

	// Test case 3: l2Token is not whitelisted
	whitelistL2TokenMap = make(map[common.Address]struct{})
	anotherL2Token := &common.Address{}
	assert.False(t, core.IsL2TokenWhitelisted(whitelistL2TokenMap, anotherL2Token))

	// Test case 4: checksum formated l2Token  is whitelisted
	whitelistL2TokenMap = make(map[common.Address]struct{})
	whitelistL2TokenMap[l2Token] = struct{}{}
	checksumL2Token := common.HexToAddress("0x1234ABCDEF")
	assert.True(t, core.IsL2TokenWhitelisted(whitelistL2TokenMap, &checksumL2Token))
}
func TestIsMinGasLimitValid(t *testing.T) {
	// Test case 1: upperMinGasLimit is nil
	var upperMinGasLimit *uint32 = nil
	minGasLimit := uint32(100)
	assert.True(t, core.IsMinGasLimitValid(upperMinGasLimit, minGasLimit))

	// Test case 2: minGasLimit is less than or equal to upperMinGasLimit
	upperMinGasLimit = uint32Ptr(200)
	minGasLimit = uint32(200)
	assert.True(t, core.IsMinGasLimitValid(upperMinGasLimit, minGasLimit))

	// Test case 3: minGasLimit is greater than upperMinGasLimit
	upperMinGasLimit = uint32Ptr(200)
	minGasLimit = uint32(250)
	assert.False(t, core.IsMinGasLimitValid(upperMinGasLimit, minGasLimit))
}

func uint32Ptr(n uint32) *uint32 {
	return &n
}
