package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nii236/pond/pkg/bot"
	"github.com/nii236/pond/pkg/commands"
	"github.com/nii236/pond/pkg/slack"
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
				Action:  runSlack,
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
	commands.Commands
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input <- scanner.Text()

		// wse := model.NewWebSocketEvent(model.WEBSOCKET_EVENT_POSTED, "", "hi1asgc5838y9crrcmi1gdbxhh", "63e4xhbaipnf7dtabh6surox8o", map[string]bool{})
		// post := &model.Post{Message: line, ChannelId: "hi1asgc5838y9crrcmi1gdbxhh", UserId: "63e4xhbaipnf7dtabh6surox8o"}
		// wse.Add("post", post.ToJson())
		// ec <- &mattermost.Event{WebSocketEvent: wse}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

type StdOut struct{}

func (s *StdOut) Write(msg string, channelID string, userID string) error {
	fmt.Println(msg)
	return nil
}
func (s *StdOut) WriteError(msg string, channelID string, userID string) error {
	fmt.Println(msg)
	return nil
}

func runWatch(c *cli.Context) error {
	b := bot.New(&StdOut{})
	argChan := b.ArgChan
	go b.Run()

	input := make(chan string)
	go listenToInput(input)

	for {
		select {
		case line := <-input:
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
			argChan <- args
		}
	}

	// return nil
}

func runSlack(c *cli.Context) error {
	b := bot.New(&StdOut{})
	argChan := b.ArgChan
	go b.Run()

	input := make(chan string)
	s := slack.New("xoxb-235722203906-DlcL76LkBPStUXcwqowVXsX0", input)
	go s.Run()

	fmt.Println("Start slack mode")
	for {
		select {
		case line := <-input:
			fmt.Println(line)
			argChan <- strings.Fields(line)
		}
	}

	// return nil
}
