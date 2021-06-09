package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/pastelnetwork/gonode/common/cli"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/tools/pastel-utility/configs"
)

var (
	errSubCommandRequired = fmt.Errorf("subcommand is required")
)

func setupStartCommand(config *configs.Config) *cli.Command {
	startCommand := cli.NewCommand("start")
	startCommand.SetUsage("usage")

	superNodeSubcommand := cli.NewCommand("supernode")
	superNodeSubcommand.SetUsage("Starts supernode")
	superNodeSubcommand.SetActionFunc(func(ctx context.Context, args []string) error {
		return runStartSuperNodeSubCommand(ctx, config)
	})
	superNodeFlags := []*cli.Flag{
		cli.NewFlag("i", &config.Start.InteractiveMode),
		cli.NewFlag("r", &config.Start.Restart),
	}
	superNodeSubcommand.AddFlags(superNodeFlags...)

	masterNodeSubCommand := cli.NewCommand("masternode")
	masterNodeSubCommand.SetUsage("Starts master node")
	masterNodeSubCommand.SetActionFunc(func(ctx context.Context, args []string) error {
		return runStartMasterNodeSubCommand(ctx, config)
	})
	masterNodeFlags := []*cli.Flag{
		cli.NewFlag("i", &config.Start.InteractiveMode),
		cli.NewFlag("r", &config.Start.Restart),
		cli.NewFlag("name", &config.Start.Name).SetUsage("name of the Master node").SetRequired(),
		cli.NewFlag("test-net", &config.Start.IsTestNet),
		cli.NewFlag("create", &config.Start.IsCreate),
		cli.NewFlag("update", &config.Start.IsUpdate),
		cli.NewFlag("txid", &config.Start.TxID),
		cli.NewFlag("ind", &config.Start.IND),
		cli.NewFlag("ip", &config.Start.IP),
		cli.NewFlag("port", &config.Start.Port),
		cli.NewFlag("pkey", &config.Start.PrivateKey),
		cli.NewFlag("pastelid", &config.Start.PastelID),
		cli.NewFlag("passphrase", &config.Start.PassPhrase),
		cli.NewFlag("rpc-ip", &config.Start.RPCIP),
		cli.NewFlag("rpc-port", &config.Start.RPCPort),
		cli.NewFlag("node-ip", &config.Start.NodeIP),
		cli.NewFlag("node-port", &config.Start.NodePort),
	}
	masterNodeSubCommand.AddFlags(masterNodeFlags...)

	nodeSubCommand := cli.NewCommand("node")
	nodeSubCommand.SetUsage("Starts specified node")
	nodeSubCommand.SetActionFunc(func(ctx context.Context, args []string) error {
		return runStartNodeSubCommand(ctx, config)
	})
	nodeFlags := []*cli.Flag{
		cli.NewFlag("i", &config.Start.InteractiveMode),
		cli.NewFlag("r", &config.Start.Restart),
	}
	nodeSubCommand.AddFlags(nodeFlags...)

	walletSubCommand := cli.NewCommand("wallet")
	walletSubCommand.SetUsage("Starts wallet")
	walletSubCommand.SetActionFunc(func(ctx context.Context, args []string) error {
		return runStartWalletSubCommand(ctx, config)
	})
	walletFlags := []*cli.Flag{
		cli.NewFlag("i", &config.Start.InteractiveMode),
		cli.NewFlag("r", &config.Start.Restart),
	}
	walletSubCommand.AddFlags(walletFlags...)

	startCommand.AddSubcommands(
		superNodeSubcommand,
		masterNodeSubCommand,
		nodeSubCommand,
		walletSubCommand,
	)

	startCommand.SetActionFunc(func(ctx context.Context, args []string) error {
		if len(args) == 0 {
			return errSubCommandRequired
		}
		return runStart(ctx, config)
	})
	return startCommand
}

func runStart(ctx context.Context, config *configs.Config) error {
	log.WithContext(ctx).Info("Start")
	defer log.WithContext(ctx).Info("End")

	log.WithContext(ctx).Infof("Config: ")
	log.WithContext(ctx).Infof(config.String())

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sys.RegisterInterruptHandler(cancel, func() {
		log.WithContext(ctx).Info("Interrupt signal received. Gracefully shutting down...")
		os.Exit(0)
	})

	// actions to run goes here

	return nil
}

func runStartNodeSubCommand(ctx context.Context, config *configs.Config) error {
	log.WithContext(ctx).Infof("Config: ")
	log.WithContext(ctx).Infof(config.String())
	return nil
}

func runStartSuperNodeSubCommand(ctx context.Context, config *configs.Config) error {
	log.WithContext(ctx).Infof("Config: ")
	log.WithContext(ctx).Infof(config.String())
	return nil
}

func runStartMasterNodeSubCommand(ctx context.Context, config *configs.Config) error {
	log.WithContext(ctx).Infof("Config: ")
	log.WithContext(ctx).Infof(config.String())
	return nil
}

func runStartWalletSubCommand(ctx context.Context, config *configs.Config) error {
	log.WithContext(ctx).Infof("Config: ")
	log.WithContext(ctx).Infof(config.String())
	return nil
}
