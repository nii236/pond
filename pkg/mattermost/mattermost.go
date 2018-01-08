package mattermost

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/mattermost/platform/model"
	"github.com/nii236/pond/pkg/logger"
	"github.com/nii236/pond/pkg/pond"
	"github.com/pkg/errors"
)

var log *logger.Log

// Client is the slack adaptor
type Client struct {
	botIn  chan *pond.Message
	botOut chan *pond.Message
	// *model.WebSocketClient
	// *model.Client
	*sync.RWMutex
	*bot
}

// New creates a new mattermost client
func New() (*Client, error) {
	log = logger.Get()
	httpClient := initHTTPClient()
	startPing(httpClient)
	bot := &bot{
		httpClient: httpClient,
	}

	bot.makeSureServerIsRunning()
	bot.loginAsTheBotUser()
	model.NewClient("")
	webSocketClient, err := model.NewWebSocketClient("wss://chat.theninja.life", bot.httpClient.AuthToken)
	if err != nil {
		return nil, err
	}

	bot.websocketClient = webSocketClient

	bot.updateTheBotUserIfNeeded()
	bot.init()
	bot.findBotTeam()
	bot.httpClient.SetTeamId(bot.team.Id)
	bot.createBotDebuggingChannelIfNeeded()
	c := &Client{
		// Client:          httpClient,
		// WebSocketClient: webSocketClient,
		bot:     bot,
		RWMutex: &sync.RWMutex{},
	}
	return c, nil
}

// Init will add the bot channels to the adaptor
func (c *Client) Init(in chan *pond.Message, out chan *pond.Message) {
	c.botIn = in
	c.botOut = out

}

// handleBot will listen for outputs from the bot
func (c *Client) handleBot() {
	for {
		select {
		case ev := <-c.botOut:
			post := &model.Post{}
			channelID, ok := ev.Meta["channelID"]
			if !ok {
				post.ChannelId = debuggingChannel.Id
			}
			rootID := ev.Meta["rootID"]

			post.ChannelId = channelID
			post.Message = ev.Message
			post.RootId = rootID
			post.MakeNonNil()

			_, err := c.httpClient.CreatePost(post)
			if err != nil {
				fmt.Println(errors.Wrap(err, "could not create post"))
				continue
			}
		}
	}
}

func (c *Client) handleMM() {
	for {
		select {
		case ev := <-c.bot.websocketClient.EventChannel:
			if ev.EventType() != "posted" {
				continue
			}

			post := &model.Post{}
			err := json.Unmarshal([]byte(ev.Data["post"].(string)), post)
			if err != nil {
				fmt.Println(errors.Wrap(err, "Could not unmarshal JSON"))
				continue
			}
			if strings.HasPrefix(post.Message, "![") {
				continue
			}
			if !strings.HasPrefix(post.Message, "!") {
				continue
			}

			fmt.Println(post.RootId)
			msg := "temp " + post.Message[1:]
			c.botIn <- &pond.Message{
				Message: msg,
				Meta: map[string]string{
					"channelID": post.ChannelId,
					"rootID":    post.RootId,
				},
			}

		}
	}
}

// Run will run the client
func (c *Client) Run() {
	err := c.bot.websocketClient.Connect()
	if err != nil {
		log.Fatalln(err)
	}
	c.bot.websocketClient.Listen()
	log.Infoln("Websocket is listening")

	go c.handleBot()
	go c.handleMM()

}
