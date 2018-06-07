package main

import (
	"fmt"

	"github.com/fatih/color"

	dgo "github.com/bwmarrin/discordgo"
)

// TODO: Colors!
func logMessageCreate(s *dgo.Session, m *dgo.MessageCreate) {
	ch, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("didnt find the channel!")
		return
	}

	g, err := s.State.Guild(ch.GuildID)
	if err != nil {
		fmt.Println("didnt find the guild!")
		return
	}

	fmt.Printf(
		"[%s #%s] @%s : %s\n",
		color.GreenString(g.Name),
		color.BlueString(ch.Name),
		color.MagentaString(m.Author.Username),
		m.ContentWithMentionsReplaced(),
	)
}
