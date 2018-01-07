package commands

import (
	"fmt"

	"github.com/nii236/pond/pkg/bot"
	"gopkg.in/urfave/cli.v2"
)

func newCryptoCmd(w bot.Writer) *cli.Command {

	cmd := &cli.Command{
		Name:    "crypto",
		Usage:   "Toolkit for crypto currencies",
		Aliases: []string{"c"},
		Subcommands: []*cli.Command{
			{
				Name:   "all",
				Action: cryptoAll(w),
			},
		},
	}

	return cmd
}

func cryptoAll(w Writer) func(c *cli.Context) error {
	return func(c *cli.Context) error {
		fmt.Fprintln(c.App.Writer, c.Args().Slice())
		return nil
	}
}

func init() {
	crypto := newCryptoCmd()
	Commands = append(Commands, crypto)
}
