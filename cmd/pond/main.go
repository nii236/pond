package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nii236/pond/pkg/bot"
	"github.com/nii236/pond/pkg/pond"
	"gopkg.in/urfave/cli.v2"
)

func main() {
	app := cli.App{
		Commands: []*cli.Command{
			{
				Name:    "watch",
				Aliases: []string{"w"},
				Action:  runWatch,
			},
			{
				Name:    "slack",
				Aliases: []string{"s"},
			},
			{
				Name:    "mattermost",
				Aliases: []string{"m"},
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
	}
	// wg := &sync.WaitGroup{}
	// wg.Add(1)
	// wg.Wait()
}

func listenToInput(input chan string) {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input <- scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}

type BotAdaptor interface {
	Init(chan *pond.Message, chan *pond.Message) error
	Run()
}

type StdOut struct {
	AdaptorIn chan *pond.Message

	BotIn  chan *pond.Message
	BotOut chan *pond.Message
}

func (s *StdOut) Init(botIn chan *pond.Message, botOut chan *pond.Message) error {
	s.BotIn = botIn
	s.BotOut = botOut

	return nil
}

func (s *StdOut) Run() {
	for {
		select {
		case out := <-s.BotOut:
			fmt.Println(out.Message)
			// fmt.Printf("%+v\n", out.Meta)
		case in := <-s.AdaptorIn:
			s.BotIn <- in
		}
	}
}

func runWatch(c *cli.Context) error {
	b := bot.New()
	inChan := b.Input()
	outChan := b.Output()
	adaptorIn := make(chan *pond.Message)
	so := &StdOut{}
	so.AdaptorIn = adaptorIn
	so.Init(inChan, outChan)

	go b.Run()
	go so.Run()

	keyInput := make(chan string)
	go listenToInput(keyInput)

	for {
		select {
		case line := <-keyInput:
			if len(line) < 1 {
				continue
			}
			if line[:1] != ";" {
				continue
			}

			line = line[1:]
			args := []string{"temp"}
			args = append(args, strings.Fields(line)...)

			if len(args) < 1 {
				fmt.Println("Not enough args")
				continue
			}
			adaptorIn <- &pond.Message{
				Message: strings.Join(args, " "),
			}
		}
	}
}
