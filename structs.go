package main

import dgo "github.com/bwmarrin/discordgo"

// Command : this is a command...
// TODO: Setup permission based commands.
type Command struct {
	Signature   string
	Description string
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

// CommandHandler : Holds a list of commands and specifies the Bot prefix.
type CommandHandler struct {
	Commands []*Command
	Prefix   string
}

type commandHandlerFunc func(s *dgo.Session, m *dgo.MessageCreate, c *Command)

// type serverStats struct {
// 	SentMessages int
// 	UsersJoined  int
// }
