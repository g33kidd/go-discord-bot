package main

import (
	dgo "github.com/bwmarrin/discordgo"
	"github.com/gorilla/websocket"
)

// Command : this is a command...
type Command struct {
	Signature   string
	Description string
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

// Bot is the bot, man!
type Bot struct {
	Session       *dgo.Session
	CmdHandler    *CommandHandler
	Conversations []Conversation
	Hub           *hub
	closed        bool
}

// CommandHandler : Holds a list of commands and specifies the Bot prefix.
// TODO: use maps for some stuff...
type CommandHandler struct {
	Commands      []*Command
	Prefix        string
	Triggers      []*Trigger
	Conversations []*Conversation
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

type twitchChannel struct {
	Status string `json:"status"`
	Game   string `json:"game"`
	ID     string `json:"_id"`
}

type channel struct {
	Status string `json:"status"`
	Game   string `json:"game"`
}

type channelUpdate struct {
	Channel channel `json:"channel"`
}

type catImage struct {
	URL string `xml:"data>images>image>url"`
}

type wsClient struct {
	send     chan []byte
	sendJSON chan interface{}
	ws       *websocket.Conn
}

type hub struct {
	clients    map[*wsClient]bool
	broadcast  chan string
	register   chan *wsClient
	unregister chan *wsClient
	content    string
}

type commandHandlerFunc func(s *dgo.Session, m *dgo.MessageCreate, c *Command)

// type serverStats struct {
// 	SentMessages int
// 	UsersJoined  int
// }
