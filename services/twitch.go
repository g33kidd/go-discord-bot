package services

import (
	"fmt"
	"time"

	"github.com/g33kidd/n00b/discord"
	"github.com/g33kidd/n00b/twitch"
)

type twitchStream struct {
	Channel string
	Live    bool
}

var streams map[string]*twitchStream

// TODO: Move these to database
var channels = []string{"g33kidd", "pixelogicdev", "neoplatonist", "naysayer88"}

// With this, need to check whether or not an announcement has already been made...

// TwitchLiveAlerts checks to see if a twitch channel recently went live and then sends a message to the discord channel.
func TwitchLiveAlerts(b *discord.Bot) {
	streams = make(map[string]*twitchStream, 0)
	for _, c := range channels {
		streams[c] = &twitchStream{
			Channel: c,
			Live:    false,
		}
	}

	for {
		for _, s := range streams {
			user := twitch.GetUser(s.Channel)
			if stream := twitch.GetStream(user.ID); stream != nil {
				if !s.Live {
					s.Live = true
					// TODO: move config to a database.
					b.Session.ChannelMessageSend("345085201392730113", fmt.Sprintf("%s is now live!", stream.Channel.DisplayName))
				}
			} else {
				s.Live = false
			}
		}

		time.Sleep(20 * time.Second)
	}
}
