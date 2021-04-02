package cli

import (
	"github.com/urfave/cli/v2"
)

// App is a wrapper of cli.App
type App struct {
	*cli.App
}

// AddCommands adds subcommands
func (app *App) AddCommands(commands ...*Command) {
	for _, command := range commands {
		app.Commands = append(app.Commands, &command.Command)
	}
}

// AddFlags adds flags
func (app *App) AddFlags(flags ...*Flag) {
	for _, flag := range flags {
		app.Flags = append(app.Flags, flag)
	}
}

// SetBefore sets the Before fucntion for the cli.App
func (app *App) SetBefore(before func() error) {
	app.Before = func(c *cli.Context) error {
		return before()
	}
}

// SetAction sets the Action function for the cli.App
func (app *App) SetAction(run func(args []string) error) {
	app.Action = func(c *cli.Context) error {
		args := []string(c.Args().Tail())
		return run(args)
	}
}

// NewApp create a new instance of the App struct
func NewApp() *App {
	app := cli.NewApp()
	app.OnUsageError = func(c *cli.Context, err error, isSubcommand bool) error {
		return err
	}

	return &App{
		App: app,
	}
}

func init() {
	cli.OsExiter = func(exitCode int) {
		// Do nothing. We just need to override this function, as the default value calls os.Exit, which
		// kills the app (or any automated test) dead in its tracks.
	}
}
