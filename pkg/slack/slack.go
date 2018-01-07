package slack

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nlopes/slack"
)

type Client struct {
	*slack.Client
	*slack.RTM
	channelID string
}

func (c *Client) updateChannelID(id string) {
	c.channelID = id
}

func (c *Client) Write(b []byte) (int, error) {
	msg := c.NewOutgoingMessage(string(b), c.channelID)
	c.SendMessage(msg)
	// channelID, timestamp, err := c.PostMessage(c.channelID, string(b), params)
	return len(b), nil
}

func (c *Client) Run() {

	go c.RTM.ManageConnection()

	for msg := range c.RTM.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.HelloEvent:
			fmt.Println("######### Connected to Slack #########")
		case *slack.MessageEvent:
			c.updateChannelID(ev.Channel)
			j, _ := json.Marshal(ev.Msg)
			log.Printf("Message: %v\n", string(j))
		case *slack.RTMError:
			fmt.Printf("Error: %s\n", ev.Error())
		}
	}

}

func New(token string, argChan chan string) *Client {
	api := slack.New(token)
	rtm := api.NewRTM()
	c := &Client{api, rtm, ""}
	return c
}
