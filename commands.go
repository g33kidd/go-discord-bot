package main

import (
	"fmt"

	dgo "github.com/bwmarrin/discordgo"
)

func testCommandHandler(s *dgo.Session, m *dgo.MessageCreate, c *Command) {
	fmt.Println("testCommandHandler called!")
	fmt.Printf("\nGot: %s -> %s\n", m.ChannelID, m.Content)
	s.ChannelMessageSend(m.ChannelID, "test ran!")

	nameP, err := c.GetParam(m.Content, "name")
	if err != nil {
		fmt.Println(err)
	}
	s.ChannelMessageSend(m.ChannelID, nameP)

	somethingP, err := c.GetParam(m.Content, "something")
	if err != nil {
		fmt.Println(err)
	}
	s.ChannelMessageSend(m.ChannelID, somethingP)

	/// All this gets the Member of who sent the message.
	// channel, err := s.State.Channel(m.ChannelID)
	// if err != nil {
	// 	fmt.Println("couldn't find the channel!")
	// 	return
	// }

	// guild, err := s.State.Guild(channel.GuildID)
	// if err != nil {
	// 	fmt.Println("couldn't find the guild!")
	// 	return
	// }

	// gm, err := s.State.Member(guild.ID, m.Author.ID)
	// if err != nil {
	// 	fmt.Println("couldn't find the user!")
	// 	return
	// }

	// for _, role := range gm.Roles {
	// 	fmt.Println(role)
	// }
}

func twitchChannelEditCommand(s *dgo.Session, m *dgo.MessageCreate, c *Command) {
	game, err := c.GetParam(m.Content, "game")
	if err != nil {
		fmt.Println("could not find the game param")
		return
	}

	title, err := c.GetParam(m.Content, "title")
	if err != nil {
		fmt.Println("could not find the title parameter.")
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Setting the twitch game to **%s** and the title to **%s**", game, title))
}
