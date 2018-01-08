package mattermost

import (
	"crypto/tls"
	"net/http"
	"os"
	"time"

	"github.com/blockninja/tobikage"
	"github.com/mattermost/platform/model"
)

var debuggingChannel *model.Channel

type bot struct {
	initialLoad     *model.InitialLoad
	commands        map[string]tobikage.Command
	user            *model.User
	team            *model.Team
	httpClient      *model.Client
	websocketClient *model.WebSocketClient
}

func startPing(client *model.Client) {
	t := time.NewTicker(1 * time.Minute)
	go func() {
		for {
			select {
			case <-t.C:
				log.Println("Pinging server...")
				_, err := client.GetPing()
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
}
func initHTTPClient() *model.Client {
	result := model.NewClient(mattermostURL)

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	result.HttpClient.Transport = tr
	return result
}

func (b *bot) makeSureServerIsRunning() {
	if props, err := b.httpClient.GetPing(); err != nil {
		log.Infoln("There was a problem pinging the Mattermost server.  Are you sure it's running?")
		os.Exit(1)
	} else {
		log.Infoln("Server detected and is running version " + props["version"])
	}
}

func (b *bot) loginAsTheBotUser() {
	if loginResult, err := b.httpClient.Login(UserEmail, UserPassword); err != nil {
		log.Println("There was a problem logging into the Mattermost server.  Are you sure ran the setup steps from the README.md?")
		log.Errorln(err)
		os.Exit(1)
	} else {
		b.user = loginResult.Data.(*model.User)
	}
}

func (b *bot) updateTheBotUserIfNeeded() {
	if b.user.FirstName != UserFirst || b.user.LastName != UserLast || b.user.Username != UserName {
		b.user.FirstName = UserFirst
		b.user.LastName = UserLast
		b.user.Username = UserName

		if updateUserResult, err := b.httpClient.UpdateUser(b.user); err != nil {
			log.Infoln("We failed to update the Sample Bot user")
			log.Errorln(err)
			os.Exit(1)
		} else {
			b.user = updateUserResult.Data.(*model.User)
			log.Infoln("Looks like this might be the first run so we've updated the bots account settings")
		}
	}
}

func (b *bot) init() {
	if initialLoadResults, err := b.httpClient.GetInitialLoad(); err != nil {
		log.Infoln("We failed to get the initial load")
		log.Errorln(err)
		os.Exit(1)
	} else {
		b.initialLoad = initialLoadResults.Data.(*model.InitialLoad)
	}
}

func (b *bot) findBotTeam() {
	for _, team := range b.initialLoad.Teams {
		if team.Name == TeamName {
			b.team = team
			break
		}
	}

	if b.team == nil {
		log.Infoln("We do not appear to be a member of the team '" + TeamName + "'")
		os.Exit(1)
	}
}

func (b *bot) createBotDebuggingChannelIfNeeded() {
	if channelsResult, err := b.httpClient.GetChannels(""); err != nil {
		log.Infoln("We failed to get the channels")
		log.Errorln(err)
	} else {
		channelList := channelsResult.Data.(*model.ChannelList)
		for _, channel := range *channelList {
			// The logging channel has alredy been created, lets just use it
			if channel.Name == ChannelLogName {
				debuggingChannel = channel
				return
			}
		}
	}

	// Looks like we need to create the logging channel
	channel := &model.Channel{}
	channel.Name = ChannelLogName
	channel.DisplayName = "Debugging For Sample Bot"
	channel.Purpose = "This is used as a test channel for logging bot debug messages"
	channel.Type = model.CHANNEL_OPEN
	if channelResult, err := b.httpClient.CreateChannel(channel); err != nil {
		log.Infoln("We failed to create the channel " + ChannelLogName)
		log.Errorln(err)
	} else {
		debuggingChannel = channelResult.Data.(*model.Channel)
		log.Infoln("Looks like this might be the first run so we've created the channel " + ChannelLogName)
	}

}
