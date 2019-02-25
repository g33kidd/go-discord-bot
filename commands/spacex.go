package commands

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/g33kidd/n00b/discord"
	spx "github.com/orcaman/spacex"
)

// Maybe I want to do a general error message function like this?
func sendErrorMessage(s *dgo.Session, c string, m string) {
	s.ChannelMessageSend(c, m)
}

// NextLaunchCommand notifies the user when the next SpaceX Launch is going to be
// using the SpaceX Data API.
func NextLaunchCommand(ctx *discord.MessageContext) {
	_, m, _, s := ctx.GetVal()
	spc := spx.New()
	launch, err := spc.GetNextLaunch()
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Couldnt get next launch sorry!")
		return
	}

	embed := &dgo.MessageEmbed{
		Title:       "/r/SpaceX API",
		Description: launch.Details,
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}

// RocketInformationCommand does this...
func RocketInformationCommand(ctx *discord.MessageContext) {
	_, m, c, s := ctx.GetVal()
	rocketName, err := c.GetParam(m.Content, "name")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	spc := spx.New()
	rocket, err := spc.GetRocket(rocketName)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "Couldnt get next launch sorry!")
		return
	}

	embed := &dgo.MessageEmbed{
		Title:       rocket.Name,
		Description: rocket.Description,
	}

	s.ChannelMessageSendEmbed(m.ChannelID, embed)
}
