package commands

import (
	"fmt"

	"gopkg.in/urfave/cli.v2"
)

func newEchoCmd() *cli.Command {
	return &cli.Command{
		Name:    "echo",
		Usage:   "copies you",
		Aliases: []string{"e"},
		Action:  runEcho,
	}
}

func runEcho(c *cli.Context) error {
	fmt.Fprint(c.App.Writer, c.Args().Slice())
	return nil
}

func init() {
	echo := newEchoCmd()
	Commands = append(Commands, echo)
}
