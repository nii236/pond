package commands

import (
	"fmt"

	"gopkg.in/urfave/cli.v2"
)

func newPingCmd() *cli.Command {

	cmd := &cli.Command{
		Name:    "ping",
		Usage:   "Returns pong when called",
		Aliases: []string{"p"},
		Action:  runPing,
	}

	return cmd

}

func runPing(c *cli.Context) error {
	fmt.Fprint(c.App.Writer, "pong")
	return nil
}

func init() {
	ping := newPingCmd()
	Commands = append(Commands, ping)
}
