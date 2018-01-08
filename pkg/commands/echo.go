package commands

import (
	"fmt"
	"strings"

	"gopkg.in/urfave/cli.v2"
)

func newEchoCmd() *cli.Command {
	return &cli.Command{
		Name:    "echo",
		Usage:   "Repeats what you say",
		Aliases: []string{"e"},
		Action:  runEcho,
	}
}

func runEcho(c *cli.Context) error {
	text := strings.Join(c.Args().Slice(), " ")
	fmt.Fprint(c.App.Writer, text)
	return nil
}

func init() {
	echo := newEchoCmd()
	Commands = append(Commands, echo)
}
