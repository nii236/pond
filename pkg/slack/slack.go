package slack

import (
	"os"
	"strings"
	"sync"

	"github.com/nii236/pond/pkg/logger"
	"github.com/nii236/pond/pkg/pond"

	"github.com/nlopes/slack"
)

var log *logger.Log

// Client is the slack adaptor
type Client struct {
	*slack.Client
	*slack.RTM

	botIn  chan *pond.Message
	botOut chan *pond.Message

	*sync.RWMutex
}

// Init will add the bot channels to the adaptor
func (c *Client) Init(in chan *pond.Message, out chan *pond.Message) {
	c.botIn = in
	c.botOut = out
}

func (c *Client) handleBot() {
	for {
		select {
		case msg := <-c.botOut:
			log.Debugln(msg)
			c.RTM.PostMessage(msg.Meta["channelID"], msg.Message, slack.NewPostMessageParameters())
		}
	}
}

func (c *Client) handleRTM() {
	for msg := range c.RTM.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			log.Infoln("Connected to Slack")
		case *slack.InvalidAuthEvent:
			log.Debugln("Invalid credentials:", ev)
			os.Exit(1)
		case *slack.ConnectedEvent:
			log.WithField("name", ev.Info.User.Name).Debugln("User connected")
		case *slack.MessageEvent:
			if !strings.HasPrefix(ev.Msg.Text, "!") {
				continue
			}
			msg := "temp " + ev.Msg.Text[1:]
			payload := &pond.Message{
				Message: msg,
				Meta: map[string]string{
					"channelID": ev.Channel,
				},
			}

			c.botIn <- payload
		case *slack.RTMError:
			log.Errorf("Error: %s\n", ev.Error())
		case *slack.ConnectionErrorEvent:
			log.Errorf("Error: %s\n", ev.Error())
		default:
			log.Debugln(ev)
		}
	}
}

// Run will run the client
func (c *Client) Run() {
	go c.handleBot()
	go c.handleRTM()
}

// New creates a new slack client
func New(token string) *Client {
	log = logger.Get()
	api := slack.New(token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()

	c := &Client{Client: api, RTM: rtm, RWMutex: &sync.RWMutex{}}
	return c
}
