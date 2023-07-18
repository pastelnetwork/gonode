package cli

import (
	"io"

	"github.com/urfave/cli/v2"
)

const (
	defaultAuthor = "Pastel Network <pastel.network>"
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
		app.Flags = append(app.Flags, flag.Flag)
	}
}

// SetBeforeFunc sets the Before function for the cli.App
// An action to execute before any subcommands are run, but after the context is ready.
func (app *App) SetBeforeFunc(beforeFn func() error) {
	app.Before = func(c *cli.Context) error {
		return beforeFn()
	}
}

// SetActionFunc sets the Action function for the cli.App
// The action to execute when no subcommands are specified.
func (app *App) SetActionFunc(actionFn ActionFn) {
	app.Action = func(c *cli.Context) error {
		args := []string(c.Args().Tail())
		return actionFn(c.Context, args)
	}
}

// SetUsage sets description of the program.
func (app *App) SetUsage(usage string) {
	app.Usage = usage
}

// SetVersion sets version of the program.
func (app *App) SetVersion(version string) {
	app.Version = version
}

// SetOutput sets writer to write output to.
func (app *App) SetOutput(write io.Writer) {
	app.Writer = write
}

// SetError sets writer to write error output to.
func (app *App) SetError(write io.Writer) {
	app.ErrWriter = write
}

// SetCustomAppHelpTemplate sets a custom help template
func (app *App) SetCustomAppHelpTemplate(appHelperTemplate string) {
	app.CustomAppHelpTemplate = appHelperTemplate
}

// NewApp create a new instance of the App struct
func NewApp(name string) *App {
	app := cli.NewApp()
	app.Name = name
	app.Authors = []*cli.Author{&cli.Author{Name: defaultAuthor}}
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
