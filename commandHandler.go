package main

import (
	"errors"
	"fmt"
	"strings"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/robertkrimen/otto"
)

// GetParam : parses the params and gets the correct value for the param with the name given.
func (c *Command) GetParam(content string, name string) (pr string, err error) {
	trimmedContent := strings.TrimPrefix(content, c.SignatureWithPrefix())

	pr = ""
	param := &CommandParameter{}
	foundParam := false
	foundMatchingParam := false
	for _, p := range c.Parameters {
		if p.Name == name {
			param = p
			foundParam = true
			break
		}
	}

	if len(c.Parameters) <= 0 {
		fmt.Println("command has no parameters")
		err = errors.New("command has no parameters set")
		return
	}

	if !foundParam {
		fmt.Printf("couldn't find the param named '%s' in the command '%s'", name, c.Signature)
		err = errors.New("could not find param")
		return
	}

	parsed := ParseParams(trimmedContent)

	if len(parsed) > len(c.Parameters) {
		err = errors.New("too many parameters passed")
		return
	}

	for pos, s := range parsed {
		if pos == param.Position {
			pr = FormatParamString(s)
			err = nil
			foundMatchingParam = true
			break
		}
	}

	if param.Required && !foundMatchingParam {
		err = errors.New("required parameter was not found")
		return
	}

	return
}

// FormatParamString removes certain characters like "" from the param response string.
func FormatParamString(pr string) string {
	f := func(r rune) bool {
		switch {
		case r == '"':
			return true
		default:
			return false
		}
	}

	return strings.TrimFunc(pr, f)
}

// SignatureWithPrefix returns a concatenated string with prefix and signature
func (c *Command) SignatureWithPrefix() string {
	return fmt.Sprintf("%s%s", prefix, c.Signature)
}

// AddParam adds a param to a command.
func (c *Command) AddParam(param *CommandParameter) {
	c.Parameters = append(c.Parameters, param)
}

// HelpString returns a string of how this Command is used.
func (c *Command) HelpString() string {
	parts := make([]string, 0)
	for _, p := range c.Parameters {
		if !p.Required {
			parts = append(parts, fmt.Sprintf("<%s:optional>", p.Name))
		} else {
			parts = append(parts, fmt.Sprintf("<%s:required>", p.Name))
		}
	}
	joined := strings.Join(parts, " ")
	return fmt.Sprintf("%s %s", c.SignatureWithPrefix(), joined)
}

// ParseParams gets a slice of arguments in a string.
// 		Goes through each rune and determines what to do with it.
//		If the rune is in Quotes "something something", it will return the whole string in quotes.
func ParseParams(content string) []string {
	inQuotes := false
	f := func(r rune) bool {
		switch {
		case r == '"':
			inQuotes = !inQuotes
			return false
		case inQuotes:
			return false
		default:
			return r == ' '
		}
	}

	return strings.FieldsFunc(content, f)
}

// NewCommandHandler creates a new CommandHandler with Prefix and sets up the default commands.
//	default commands: help
func NewCommandHandler(Prefix string) *CommandHandler {
	cmdh := &CommandHandler{
		Prefix: Prefix,
	}

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

				me, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
				if err != nil {
					panic(err)
				} else {
					fmt.Println(me)
				}
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

				me, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
				if err != nil {
					panic(err)
				} else {
					fmt.Println(me)
				}
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

	for _, cmd := range ch.Commands {
		if alreadyFound {
			break
		}

		if withPrefix && strings.HasPrefix(content, cmd.SignatureWithPrefix()) {
			c = cmd
			err = nil
			alreadyFound = true
		}

		if !withPrefix && (content == cmd.Signature || content == cmd.SignatureWithPrefix()) {
			c = cmd
			err = nil
			alreadyFound = true
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
// TODO: Replace console.log and have it send a message to the discord channel instead of fmt.Print
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
	vm := otto.New()
	vm.Set("discordLog", func(call otto.FunctionCall) otto.Value {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("`%s`", call.Argument(0).String()))
		return otto.Value{}
	})

	// Evaluate the code given in the code block and get a Value
	trimmedCodeBlock := strings.TrimFunc(m.Content, f)
	val, err := vm.Eval(trimmedCodeBlock)
	if err != nil {
		panic(err)
	}

	// Try to get that Value as a string.
	r, err = val.ToString()

	return
}
