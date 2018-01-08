package bot

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"sync"
	"text/tabwriter"

	"github.com/nii236/pond/pkg/pond"

	"github.com/alecthomas/template"
	"github.com/nii236/pond/pkg/commands"
	"gopkg.in/urfave/cli.v2" // imports as package "cli"
)

// Bot contains the state for the bot
type Bot struct {
	inChan  chan *pond.Message
	outChan chan *pond.Message
	*cli.App
	*sync.RWMutex
}

// Input returns the input channel for the bot
func (bot *Bot) Input() chan *pond.Message {
	return bot.inChan

}

// Output returns the output channel for the bot
func (bot *Bot) Output() chan *pond.Message {
	return bot.outChan
}

// Write will write to the output channel
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
	// 	templ = `
	// {{range .VisibleCommands}}
	// {{join .Names ", "}}{{"\t"}}{{.Usage}}
	// {{end}}
	// `
	funcMap := template.FuncMap{
		"join": strings.Join,
	}

	bw := bufio.NewWriter(out)
	w := tabwriter.NewWriter(bw, 1, 8, 2, ' ', 0)
	t := template.Must(template.New("help").Funcs(funcMap).Parse(templ))

	// app := data.(*cli.App)
	// fmt.Println(app.VisibleCommands()[0].Name)

	err := t.Execute(w, data.(*cli.App))
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

func notFound(c *cli.Context, cmd string) {
	fmt.Fprint(c.App.Writer, "Command not found:", cmd)
}

// New will return a new Bot
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
	result.RWMutex = &sync.RWMutex{}
	return result

}

// Run will start the bot
func (bot *Bot) Run() {
	for {
		select {
		case in := <-bot.inChan:
			bot.Lock()

			meta := in.Meta
			metaInterface := map[string]interface{}{}
			for key, value := range meta {
				metaInterface[key] = value
			}

			bot.App.Metadata = metaInterface

			err := bot.App.Run(strings.Fields(in.Message))
			if err != nil {
				bot.Write([]byte(err.Error()))
				fmt.Println(err)
			}

			bot.Unlock()
		}
	}
}
