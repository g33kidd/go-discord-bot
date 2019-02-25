package commands

import (
	"fmt"
	"strconv"

	dgo "github.com/bwmarrin/discordgo"

	"github.com/g33kidd/n00b/discord"
	"github.com/g33kidd/n00b/twitch"
)

// TwitchChannelEditCommand edits the twitch channel!
func TwitchChannelEditCommand(ctx *discord.MessageContext) {
	_, m, c, s := ctx.GetVal()
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

// TwitchChannelInfoCommand edits the twitch channel!
func TwitchChannelInfoCommand(ctx *discord.MessageContext) {
	_, m, c, s := ctx.GetVal()
	channel, err := c.GetParam(m.Content, "channel")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "I need the `channel` param. Try `$help twitch` for more info.")
		return
	}

	user := twitch.GetUser(channel)
	if user == nil {
		s.ChannelMessageSend(m.ChannelID, "could not find that channel.")
		return
	}

	cinfo := twitch.GetChannel(user.ID)
	if cinfo == nil {
		s.ChannelMessageSend(m.ChannelID, "Something's not right.")
		return
	}

	fmt.Println(cinfo)

	embed := &dgo.MessageEmbed{}
	embed.Title = fmt.Sprintf("%s's Channel", cinfo.DisplayName)

	color, _ := strconv.ParseInt("228B22", 16, 64)

	embed.Color = int(color)

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
