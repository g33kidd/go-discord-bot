package commands

import (
	"fmt"

	dgo "github.com/bwmarrin/discordgo"

	"github.com/g33kidd/n00b/discord"
	"github.com/g33kidd/n00b/twitch"
)

// TwitchChannelEditCommand edits the twitch channel!
func TwitchChannelEditCommand(s *dgo.Session, m *dgo.MessageCreate, c *discord.Command) {
	ch := twitch.GetMyChannel()

	game, err := c.GetParam(m.Content, "game")
	if err != nil {
		fmt.Println(err)
		return
	}

	title, err := c.GetParam(m.Content, "title")
	if err != nil {
		fmt.Println(err)
		return
	}

	chUpdate := &twitch.TwitchChannelEditData{
		Game:   game,
		Status: title,
	}

	if update := twitch.UpdateChannel(ch.ID, chUpdate); update != nil {
		s.ChannelMessageSend(m.ChannelID, "Something did not go as planned!")
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I set the game to **%s** and the title to **%s**", game, title))
	}
}
