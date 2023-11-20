package main

import (
	"context"
	"os"

	oplog "github.com/ethereum-optimism/optimism/op-service/log"
	"github.com/ethereum-optimism/optimism/op-service/opio"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/params"
	"github.com/urfave/cli/v2"
)

var (
	ConfigFlag = &cli.StringFlag{
		Name:    "config",
		Value:   "./bot.toml",
		Aliases: []string{"c"},
		Usage:   "path to config file",
		EnvVars: []string{"BOT_CONFIG"},
	}
)

var (
	GitCommit = ""
	GitDate   = ""
)

func main() {
	// This is the most root context, used to propagate
	// cancellations to all spawned application-level goroutines
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		opio.BlockOnInterrupts()
		cancel()
	}()

	oplog.SetupDefaults()
	app := newCli(GitCommit, GitDate)
	if err := app.RunContext(ctx, os.Args); err != nil {
		log.Error("application failed", "err", err)
		os.Exit(1)
	}
}

func newCli(GitCommit string, GitDate string) *cli.App {
	flags := []cli.Flag{ConfigFlag}
	flags = append(flags, oplog.CLIFlags("OPBNB_BRIDGE_BOT")...)
	return &cli.App{
		Name:                 "opbnb-bridge-bot",
		Version:              params.VersionWithCommit(GitCommit, GitDate),
		Description:          "opbnb-bridge-bot",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:        "run",
				Flags:       flags,
				Description: "Runs the indexing service",
				Action:      runCommand,
			},
			{
				Name:        "version",
				Description: "print version",
				Action: func(ctx *cli.Context) error {
					cli.ShowVersion(ctx)
					return nil
				},
			},
		},
	}
}
