package commands

import (
	"fmt"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/g33kidd/n00b/discord"
)

// HelpCommand displays some help information
func HelpCommand(ctx *discord.MessageContext) {
	b, m, c, s := ctx.GetVal()
	cmdHelp, _ := c.GetParam(m.Content, "command")

	if cmdHelp != "" {
		foundCmd, err := b.CmdHandler.FindCommand(cmdHelp, false)
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s", err))
			return
		}

		embed := &dgo.MessageEmbed{
			Description: foundCmd.Description,
			Title:       foundCmd.HelpString(),
		}

		for _, p := range foundCmd.Parameters {
			field := &dgo.MessageEmbedField{
				Name:  p.Name,
				Value: p.Description,
			}

			embed.Fields = append(embed.Fields, field)
		}

		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	} else {
		embed := &dgo.MessageEmbed{
			Description: "HEEEEELPP!",
			Title:       "This is a help message!",
		}

		for _, cmd := range b.CmdHandler.Commands {
			if cmd.Description != "" && cmd.Signature != "" {
				field := &dgo.MessageEmbedField{
					Name:   cmd.HelpString(),
					Value:  cmd.Description,
					Inline: false,
				}
				embed.Fields = append(embed.Fields, field)
			}
		}

		s.ChannelMessageSendEmbed(m.ChannelID, embed)
	}
}

// HelpCommand helps the user
// TODO: Figure out how this can go here
// TODO: Okay this can go here now since we have a full context of this current message and command usage.
// func HelpCommand(ctx *discord.MessageContext) {
// 	_, m, c, s := ctx.GetVal()

// }

// // BanCommand allows an administrator to ban a user.
// func BanCommand(ctx *discord.MessageContext) {
// 	_, m, c, s := ctx.GetVal()
// 	b.ChannelMessageSend(m.ChannelID, "test")
// }

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
