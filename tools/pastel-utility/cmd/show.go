package cmd

import (
	"context"
	"os"

	"github.com/pastelnetwork/gonode/common/cli"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/tools/pastel-utility/configs"
)

func setupShowCommand(config *configs.Config) *cli.Command {
	showCommand := cli.NewCommand("show")
	showCommand.SetUsage("usage")
	showCommandFlags := []*cli.Flag{
		// TODO: Flags will go here
	}
	showCommand.AddFlags(showCommandFlags...)

	showCommand.SetActionFunc(func(ctx context.Context, args []string) error {

		return runShow(ctx, config)
	})
	return showCommand
}

func runShow(ctx context.Context, config *configs.Config) error {
	log.WithContext(ctx).Info("Show")
	defer log.WithContext(ctx).Info("End")

	log.WithContext(ctx).Infof("Config: %s", config)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sys.RegisterInterruptHandler(cancel, func() {
		log.WithContext(ctx).Info("Interrupt signal received. Gracefully shutting down...")
		os.Exit(0)
	})

	// actions to run goes here
	return nil
}
