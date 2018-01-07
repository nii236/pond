package commands

import (
	"fmt"

	"gopkg.in/urfave/cli.v2"
)

func newCryptoCmd() *cli.Command {

	cmd := &cli.Command{
		Name:    "crypto",
		Usage:   "Toolkit for crypto currencies",
		Aliases: []string{"c"},
		Subcommands: []*cli.Command{
			{
				Name:   "all",
				Action: cryptoAll,
			},
		},
	}

	return cmd
}

func cryptoAll(c *cli.Context) error {
	fmt.Fprint(c.App.Writer, c.Args().Slice())
	return nil
}

func init() {
	crypto := newCryptoCmd()
	Commands = append(Commands, crypto)
}
