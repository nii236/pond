package commands

import (
	"fmt"

	"gopkg.in/urfave/cli.v2"
)

type Echo struct {
	*cli.Command
}

func newEchoCmd() *Echo {
	cmd := &cli.Command{
		Name:    "echo",
		Usage:   "copies you",
		Aliases: []string{"e"},
		Action:  runEcho,
	}

	return &Echo{cmd}

}

func runEcho(c *cli.Context) error {
	fmt.Fprintln(c.App.Writer, c.Args().Slice())
	return nil
}

func init() {
	echo := newEchoCmd()
	Commands = append(Commands, echo)
}
