package cmd

import (
	"github.com/urfave/cli/v2"
)

var (
	cmds []*cli.Command
)

func Run(args []string) error {
	app := &cli.App{
		Commands: cmds,
	}

	return app.Run(args)
}

func registerCommand(newCmd *cli.Command) {
	cmds = append(cmds, newCmd)
}
