package cmd

import (
	"context"
	"os"

	"github.com/pastelnetwork/gonode/common/cli"
	"github.com/pastelnetwork/gonode/common/log"
	"github.com/pastelnetwork/gonode/common/sys"
	"github.com/pastelnetwork/gonode/tools/pastel-utility/configs"
)

func setupStopCommand(config *configs.Config) *cli.Command {
	stopCommand := cli.NewCommand("stop")
	stopCommand.SetUsage("usage")
	stopCommandFlags := []*cli.Flag{
		// TODO: Flags will go here
	}
	stopCommand.AddFlags(stopCommandFlags...)

	stopCommand.SetActionFunc(func(ctx context.Context, args []string) error {
		return runStop(ctx, config)
	})
	return stopCommand
}

func runStop(ctx context.Context, config *configs.Config) error {
	log.WithContext(ctx).Info("Stop")
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
