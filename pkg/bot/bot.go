package bot

import (
	"fmt"

	"github.com/nii236/pond/pkg/commands"
	"gopkg.in/urfave/cli.v2" // imports as package "cli"
)

type Writer interface {
	Write(msg string, channelID string, userID string) error
	WriteError(msg string, channelID string, userID string) error
}
type Bot struct {
	Writer Writer
	*cli.App
	ArgChan chan []string
}

func notFound(c *cli.Context, cmd string) {
	fmt.Fprintln(c.App.Writer, "Command not found:", cmd)

}

func New(w Writer) *Bot {
	app := &cli.App{
		Name:  "Pond",
		Usage: "HI",
		Authors: []*cli.Author{
			{
				Name:  "John Nguyen",
				Email: "jtnguyen236@gmail.com",
			},
		},
		UsageText:   "Its like a puddle but better",
		Description: "Useful bot stuffs",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:   "user",
				Value:  "",
				Hidden: true,
			},
			&cli.StringFlag{
				Name:   "channel",
				Value:  "",
				Hidden: true,
			},
		},
		CommandNotFound: notFound,
		Commands:        commands.Commands,
	}

	argChan := make(chan []string)
	return &Bot{
		App:     app,
		ArgChan: argChan,
		Writer:  w,
	}
}

func (b *Bot) Run() {
	for {
		select {
		case args := <-b.ArgChan:
			err := b.App.Run(args)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}
