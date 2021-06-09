package cmd

import (
	"context"
	"os"

	"github.com/pastelnetwork/gonode/common/cli"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/tools/pastel-utility/configs"
)

func setupInstallCommand(config *configs.Config) *cli.Command {
	installCommand := cli.NewCommand("install")
	installCommand.SetUsage("usage")
	installCommandFlags := []*cli.Flag{
		// TODO: Add flags here
	}
	installCommand.AddFlags(installCommandFlags...)

	installCommand.SetActionFunc(func(ctx context.Context, args []string) error {
		return runInstall(ctx, config)
	})
	return installCommand
}

func runInstall(ctx context.Context, config *configs.Config) error {
	log.WithContext(ctx).Info("Install")
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
