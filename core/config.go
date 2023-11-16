package core

import (
	"errors"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/ethereum-optimism/optimism/indexer/config"
	"github.com/ethereum/go-ethereum/log"
)

const (
	defaultConfirmBlocks       = 15
	defaultLogFilterBlockRange = 100
)

type Config struct {
	Misc        MiscConfig         `toml:"misc"`
	RPCs        config.RPCsConfig  `toml:"rpcs"`
	DB          config.DBConfig    `toml:"db"`
	L1Contracts config.L1Contracts `toml:"l1-contracts"`
}

type MiscConfig struct {
	L2StandardBridgeBot string `toml:"l2-standard-bridge-bot"`
	ProposeTimeWindow   int64  `toml:"propose-time-window"`
	ChallengeTimeWindow int64  `toml:"challenge-time-window"`
	ConfirmBlocks       int64  `toml:"confirm-blocks"`
	LogFilterBlockRange int64  `toml:"log-filter-block-range"`
}

// LoadConfig loads the `bot.toml` config file from a given path
func LoadConfig(log log.Logger, path string) (Config, error) {
	log.Debug("loading config", "path", path)

	var conf Config
	data, err := os.ReadFile(path)
	if err != nil {
		return conf, err
	}

	data = []byte(os.ExpandEnv(string(data)))
	log.Debug("parsed config file", "data", string(data))
	if _, err := toml.Decode(string(data), &conf); err != nil {
		log.Info("failed to decode config file", "err", err)
		return conf, err
	}

	if conf.Misc.ConfirmBlocks == 0 {
		log.Info("setting default confirm blocks", "confirm-blocks", defaultConfirmBlocks)
		conf.Misc.ConfirmBlocks = defaultConfirmBlocks
	}
	if conf.Misc.LogFilterBlockRange == 0 {
		log.Info("setting default log filter block range", "log-filter-block-range", defaultLogFilterBlockRange)
		conf.Misc.LogFilterBlockRange = defaultLogFilterBlockRange
	}
	if conf.Misc.ProposeTimeWindow == 0 {
		return conf, errors.New("propose-time-window must be set")
	}
	if conf.Misc.ChallengeTimeWindow == 0 {
		return conf, errors.New("challenge-time-window must be set")
	}

	log.Info("loaded config")
	return conf, nil
}
