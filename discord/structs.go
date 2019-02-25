package discord

import dgo "github.com/bwmarrin/discordgo"

// Bot is a discord bot, yo!
type Bot struct {
	Session       *dgo.Session
	Prefix        string
	CmdHandler    *CommandHandler
	closed        bool
	MessageCreate chan *dgo.MessageCreate
}

// CommandHandler : Holds a list of commands and specifies the Bot prefix.
// TODO: use maps for some stuff...
type CommandHandler struct {
	Commands      []*Command
	Prefix        string
	Triggers      []*Trigger
	Conversations []*Conversation
}

// MessageContext is usually just for when Command handlers are being run, but this could be used for different things.
type MessageContext struct {
	Session *dgo.Session
	Bot     *Bot
	Message *dgo.MessageCreate
	Command *Command
}

// Command : this is a command...
type Command struct {
	Signature   string
	Description string
	Prefix      string
	Aliases     []string
	Parameters  []*CommandParameter
	Handler     commandHandlerFunc
}

// CommandParameter : Specifies the name, position and what this parameter does.
type CommandParameter struct {
	Name        string
	Description string
	Position    int
	Required    bool
}

// Trigger takes a pattern and returns a response
type Trigger struct {
	Pattern  string
	Response string
}

// Conversation between the bot and a user..
type Conversation struct {
	UserID string
}

type commandHandlerFunc func(s *dgo.Session, m *dgo.MessageCreate, c *Command)
