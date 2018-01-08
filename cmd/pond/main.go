package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/nii236/pond/pkg/bot"
	"github.com/nii236/pond/pkg/logger"
	"github.com/nii236/pond/pkg/mattermost"
	"github.com/nii236/pond/pkg/pond"
	"github.com/nii236/pond/pkg/slack"
	"github.com/nii236/pond/pkg/stdout"
	"github.com/pkg/errors"
	"gopkg.in/urfave/cli.v2"
)

var log *logger.Log

func main() {
	app := cli.App{
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "production",
				Aliases: []string{"p"},
			},
			&cli.BoolFlag{
				Name:    "debug",
				Aliases: []string{"d"},
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "watch",
				Aliases: []string{"w"},
				Action:  runWatch,
			},
			{
				Name:    "mattermost",
				Aliases: []string{"m"},
				Action:  runMattermost,
			},
			{
				Name:    "slack",
				Aliases: []string{"s"},
				Action:  runSlack,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "slack-token",
						Aliases: []string{"s"},
					},
				},
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

func runMattermost(c *cli.Context) error {
	logger.New(c.Bool("production"), c.Bool("debug"))
	mm, err := mattermost.New()
	if err != nil {
		return errors.Wrap(err, "Could not create new Mattermost bot")
	}

	b := bot.New()
	go b.Run()

	inChan := b.Input()
	outChan := b.Output()
	mm.Init(inChan, outChan)
	mm.Run()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()
	return nil
}
func runSlack(c *cli.Context) error {
	logger.New(c.Bool("production"), c.Bool("debug"))
	log = logger.Get()
	log.Infoln("Running in Slack mode")
	b := bot.New()
	go b.Run()

	inChan := b.Input()
	outChan := b.Output()
	token := c.String("slack-token")
	s := slack.New(token)

	s.Init(inChan, outChan)
	s.Run()

	wg := sync.WaitGroup{}
	wg.Add(1)
	wg.Wait()

	return nil
}

func runWatch(c *cli.Context) error {
	logger.New(c.Bool("production"), c.Bool("debug"))
	log = logger.Get()
	b := bot.New()
	inChan := b.Input()
	outChan := b.Output()
	adaptorIn := make(chan *pond.Message)
	so := &stdout.Adaptor{}
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
