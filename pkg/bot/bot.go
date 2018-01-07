package bot

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/nii236/pond/pkg/pond"

	"github.com/alecthomas/template"
	"github.com/nii236/pond/pkg/commands"
	"gopkg.in/urfave/cli.v2" // imports as package "cli"
)

func (b *Bot) Input() chan *pond.Message {
	return b.inChan

}

func (b *Bot) Output() chan *pond.Message {
	return b.outChan
}

func (bot *Bot) Write(b []byte) (int, error) {
	meta := bot.App.Metadata
	metaString := map[string]string{}
	for key, value := range meta {
		switch value := value.(type) {
		case string:
			metaString[key] = value
		}
	}

	bot.outChan <- &pond.Message{
		Message: string(b),
		Meta:    metaString,
	}
	return len(b), nil
}

func helpPrinter(out io.Writer, templ string, data interface{}) {
	funcMap := template.FuncMap{
		"join": strings.Join,
	}
	cli.HelpPrinter = helpPrinter
	bw := bufio.NewWriter(out)
	w := tabwriter.NewWriter(bw, 1, 8, 2, ' ', 0)
	t := template.Must(template.New("help").Funcs(funcMap).Parse(templ))

	err := t.Execute(w, data)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = bw.Flush()
	if err != nil {
		fmt.Println(err)
		return
	}
}

type Bot struct {
	inChan  chan *pond.Message
	outChan chan *pond.Message
	*cli.App
}

func notFound(c *cli.Context, cmd string) {
	fmt.Fprint(c.App.Writer, "Command not found:", cmd)
}

func New() *Bot {
	result := &Bot{}
	result.App = &cli.App{
		Name:  "Pond",
		Usage: "HI",
		Authors: []*cli.Author{
			{
				Name:  "John Nguyen",
				Email: "jtnguyen236@gmail.com",
			},
		},
		UsageText:       "Its like a puddle but better",
		Description:     "Useful bot stuffs",
		Writer:          result,
		ErrWriter:       result,
		CommandNotFound: notFound,
		Commands:        commands.Commands,
	}

	cli.HelpPrinter = helpPrinter

	result.inChan = make(chan *pond.Message)
	result.outChan = make(chan *pond.Message)

	return result

}

func (b *Bot) Run() {
	for {
		select {
		case in := <-b.inChan:
			meta := in.Meta
			metaInterface := map[string]interface{}{}
			for key, value := range meta {
				metaInterface[key] = value
			}

			b.App.Metadata = metaInterface

			err := b.App.Run(strings.Fields(in.Message))
			if err != nil {
				b.Write([]byte(err.Error()))
				fmt.Println(err)
			}
		}
	}
}
