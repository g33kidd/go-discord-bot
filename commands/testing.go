package commands

import (
	"fmt"

	dgo "github.com/bwmarrin/discordgo"
	dc "github.com/g33kidd/n00b/discord"
)

// MacroCommand does a macro thing
func MacroCommand(s *dgo.Session, m *dgo.MessageCreate, c *dc.Command) {
	content := m.Content

	// TODO: Re-work this so that I can return nil for name or cmd1. Just requires returning it as a pointer.
	name, err := c.GetParam(content, "name")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "whoops... you are missing a few params. Try `$help macro` for more information!")
		return
	}

	cmd1, _ := c.GetParam(content, "cmd1")
	// cmd2, _ := c.GetParam(content, "cmd2")
	// cmd3, _ := c.GetParam(content, "cmd3")

	if cmd1 == "" && name != "" {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Okay, running this macro called `%s`.", name))
		return
	}

	fmt.Println(cmd1)

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Okay creating this macro `%s` doing `%s`", name, cmd1))
}
