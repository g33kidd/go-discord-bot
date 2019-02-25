package discord

import (
	"bufio"
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
	helpCommand := NewCommand("help", "Displays help messages.", func(s *dgo.Session, m *dgo.MessageCreate, c *Command) {
		cmdHelp, _ := c.GetParam(m.Content, "command")

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
	})

	helpCommand.AddParameter("command", "Gets the help message for the command given, if there is one.", false)

	cmdh.AddCommand(helpCommand)
	return cmdh
}

// AddCommand adds a command to the command handler.
func (ch *CommandHandler) AddCommand(cmd *Command) {
	cmd.Prefix = ch.Prefix
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
// TODO: Can an API be generalized for creating Maybe handlers?
// TODO: Create a handler for a specific language. Do something!!!!
func (ch *CommandHandler) MaybeHandleCodeBlock(s *dgo.Session, m *dgo.MessageCreate) (res string, err error) {
	content := m.Content
	scanner := bufio.NewScanner(strings.NewReader(content))

	tickTrimFunc := func(r rune) bool {
		switch {
		case r == '`':
			return true
		default:
			return false
		}
	}

	// check if this really is a code block.
	start := strings.HasPrefix(content, "```")
	end := strings.HasSuffix(content, "```")
	if !start && !end {
		res = ""
		err = errors.New("was not a code block")
		return
	}

	// figure out what language this is, or if any is set.
	firstLine := true
	var language string
	for scanner.Scan() {
		line := scanner.Text()

		if !firstLine {
			fmt.Println("not first line...")
			break
		}

		if strings.HasPrefix(line, "```") && strings.HasSuffix(line, "```") && len(line) > 3 {
			res = "Codeblock should be formatted properly.\nNewline after start and before end block ticks."
			return
		}

		language = strings.TrimFunc(line, tickTrimFunc)

		fmt.Printf("%s\n", line)

		firstLine = false
	}

	// remove the back ticks.
	blockText := strings.TrimFunc(content, tickTrimFunc)
	blockText = strings.Trim(blockText, language)

	// return fmt.Sprintf("language was **%s**\ntext was\n```%s\n%s\n```", language, language, blockText), nil

	// codeBlockStart := strings.HasPrefix(m.Content, "```")
	// codeBlockEnd := strings.HasSuffix(m.Content, "```")

	// // Remove the backticks from the message. It would mess up the eval if we didn't!

	// if !codeBlockStart && !codeBlockEnd {
	// 	r = ""
	// 	err = errors.New("did not find a code block")
	// 	return
	// }

	// // Setup the Otto VM and add a function called `discordLog` that can be used in the JS code.
	// // TODO: this https://godoc.org/github.com/robertkrimen/otto#hdr-Halting_Problem
	vm := otto.New()
	vm.Set("botToken", "Nice try! No bot tokens here!")
	vm.Set("discordLog", func(call otto.FunctionCall) otto.Value {
		s.ChannelMessageSend(m.ChannelID, call.Argument(0).String())
		return otto.Value{}
	})

	// // Evaluate the code given in the code block and get a Value
	// trimmedCodeBlock := strings.TrimFunc(m.Content, f)
	val, err := vm.Eval(blockText)
	if err != nil {
		if val.IsUndefined() {
			res = fmt.Sprintf("Something went wrong when evaluating your code! Here's what happened:\n```%s```", err)
		} else {
			res = "I don't know what happened. Sorry!"
		}

		return
	}

	// // Try to get that Value as a string.
	res, err = val.ToString()

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
