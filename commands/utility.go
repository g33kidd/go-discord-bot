package commands

import (
	dgo "github.com/bwmarrin/discordgo"
	dc "github.com/g33kidd/n00b/discord"
)

// HelpCommand helps the user
// TODO: Figure out how this can go here
// TODO: Okay this can go here now since we have a full context of this current message and command usage.
func HelpCommand(b *dgo.Session, m *dgo.MessageCreate, c *dc.Command) {

}

// BanCommand allows an administrator to ban a user.
func BanCommand(b *dgo.Session, m *dgo.MessageCreate, c *dc.Command) {
	b.ChannelMessageSend(m.ChannelID, "test")
}

// func HelpCommand(s *dgo.Session, m *dgo.MessageCreate, c *dc.Command) {
// 	cmdHelp, err := c.GetParam(m.Content, "command")
// 	if err != nil {
// 		fmt.Println("Something happened")
// 		fmt.Println(err)
// 	}

// 	if cmdHelp != "" {
// 		foundCmd, err := cmdh.FindCommand(cmdHelp, false)
// 		if err != nil {
// 			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s", err))
// 			return
// 		}

// 		embed := &dgo.MessageEmbed{
// 			Description: foundCmd.Description,
// 			Title:       foundCmd.HelpString(),
// 		}

// 		for _, p := range foundCmd.Parameters {
// 			field := &dgo.MessageEmbedField{
// 				Name:  p.Name,
// 				Value: p.Description,
// 			}

// 			embed.Fields = append(embed.Fields, field)
// 		}

// 		s.ChannelMessageSendEmbed(m.ChannelID, embed)
// 	} else {
// 		embed := &dgo.MessageEmbed{
// 			Description: "HEEEEELPP!",
// 			Title:       "This is a help message!",
// 		}

// 		for _, cmd := range cmdh.Commands {
// 			if cmd.Description != "" && cmd.Signature != "" {
// 				field := &dgo.MessageEmbedField{
// 					Name:   cmd.HelpString(),
// 					Value:  cmd.Description,
// 					Inline: false,
// 				}
// 				embed.Fields = append(embed.Fields, field)
// 			}
// 		}

// 		s.ChannelMessageSendEmbed(m.ChannelID, embed)
// 	}
// }
