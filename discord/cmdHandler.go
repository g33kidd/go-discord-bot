package discord

import (
	"errors"
	"fmt"
	"strings"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/robertkrimen/otto"
)

// NewCommandHandler creates a new CommandHandler with Prefix and sets up the default commands.
//	default commands: help
// TODO: Move this away from here!
func NewCommandHandler(Prefix string) *CommandHandler {
	cmdh := &CommandHandler{
		Prefix: Prefix,
	}

	// TODO: does this really need to be in here?
	helpCommand := &Command{
		Signature:   "help",
		Description: "Help meeeeee!",
		Handler: func(s *dgo.Session, m *dgo.MessageCreate, c *Command) {
			cmdHelp, err := c.GetParam(m.Content, "command")
			if err != nil {
				fmt.Println("Something happened")
				fmt.Println(err)
			}

			if cmdHelp != "" {
				foundCmd, err := cmdh.FindCommand(cmdHelp, false)
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

				for _, cmd := range cmdh.Commands {
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
		},
	}

	helpCommand.AddParam(&CommandParameter{
		Name:        "command",
		Description: "Gets the help for this command only!",
	})

	cmdh.AddCommand(helpCommand)
	return cmdh
}

// AddCommand adds a command to the command handler.
func (ch *CommandHandler) AddCommand(cmd *Command) {
	ch.Commands = append(ch.Commands, cmd)
}

// FindCommand finds a command (if it exists) within the commandHandler.
func (ch *CommandHandler) FindCommand(content string, withPrefix bool) (c *Command, err error) {
	alreadyFound := false

	contentSlice := strings.Split(content, " ")
	content = contentSlice[0]

	for _, cmd := range ch.Commands {
		if withPrefix && content == cmd.SignatureWithPrefix() {
			c = cmd
			err = nil
			alreadyFound = true
			break
		}

		if !withPrefix && content == cmd.Signature {
			c = cmd
			err = nil
			alreadyFound = true
			break
		}
	}

	if !alreadyFound {
		c = nil
		err = errors.New("could not find the command")
	}

	return
}

// MaybeHandleCodeBlock takes a string and figures out if it's in a codeblock format.
// TODO: figure out what language is being passed in. ```js <- js should be the language.
// TODO: Have a separate Handler for each language that is supported.
func (ch *CommandHandler) MaybeHandleCodeBlock(s *dgo.Session, m *dgo.MessageCreate) (r string, err error) {
	codeBlockStart := strings.HasPrefix(m.Content, "```")
	codeBlockEnd := strings.HasSuffix(m.Content, "```")

	// Remove the backticks from the message. It would mess up the eval if we didn't!
	f := func(r rune) bool {
		switch {
		case r == '`':
			return true
		default:
			return false
		}
	}

	if !codeBlockStart && !codeBlockEnd {
		r = ""
		err = errors.New("did not find a code block")
		return
	}

	// Setup the Otto VM and add a function called `discordLog` that can be used in the JS code.
	// TODO: this https://godoc.org/github.com/robertkrimen/otto#hdr-Halting_Problem
	vm := otto.New()
	vm.Set("botToken", "Nice try! No bot tokens here!")
	vm.Set("discordLog", func(call otto.FunctionCall) otto.Value {
		s.ChannelMessageSend(m.ChannelID, call.Argument(0).String())
		return otto.Value{}
	})

	// Evaluate the code given in the code block and get a Value
	trimmedCodeBlock := strings.TrimFunc(m.Content, f)
	val, err := vm.Eval(trimmedCodeBlock)
	if err != nil {
		if val.IsUndefined() {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Something went wrong when evaluating your code! Here's what happened:\n```%s```", err))
		} else {
			s.ChannelMessageSend(m.ChannelID, "I don't know what happened. Sorry!")
		}

		return
	}

	// Try to get that Value as a string.
	r, err = val.ToString()

	return
}

func (ch *CommandHandler) startConversation(id string) {
	conv := &Conversation{
		UserID: id,
	}
	ch.Conversations = append(ch.Conversations, conv)
}

func (ch *CommandHandler) stopConversation(id string) {
	// remove conversation from the list...
}

func (ch *CommandHandler) hasConversation(id string) bool {
	found := false
	for _, c := range ch.Conversations {
		if c.UserID == id {
			found = true
			break
		}
	}
	return found
}

// AddMessageTrigger adds a new Trigger to the Triggers slice on CommandHandler
func (ch *CommandHandler) AddMessageTrigger(s string, r string) {
	ch.Triggers = append(ch.Triggers, &Trigger{Pattern: s, Response: r})
}

// MaybeHandleMessageTrigger looks for a pattern in any message and responds
func (ch *CommandHandler) MaybeHandleMessageTrigger(s *dgo.Session, m *dgo.MessageCreate) (r string, err error) {
	content := m.Content
	foundTrigger := false
	err = nil

	for _, t := range ch.Triggers {
		if strings.Contains(content, t.Pattern) {
			foundTrigger = true
			r = t.Response
			break
		}
	}

	if !foundTrigger {
		err = errors.New("didn't find a message trigger")
	}

	return
}
