package discord

import dgo "github.com/bwmarrin/discordgo"

// Bot is a discord bot, yo!
type Bot struct {
	Session    *dgo.Session
	Prefix     string
	CmdHandler *CommandHandler
	closed     bool
}

// CommandHandler : Holds a list of commands and specifies the Bot prefix.
// TODO: use maps for some stuff...
type CommandHandler struct {
	Commands      []*Command
	Prefix        string
	Triggers      []*Trigger
	Conversations []*Conversation
}

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

// package main

// import (
// 	dgo "github.com/bwmarrin/discordgo"
// 	"github.com/gorilla/websocket"
// )

// // Command : this is a command...
// type Command struct {
// 	Signature   string
// 	Description string
// 	Aliases     []string
// 	Parameters  []*CommandParameter
// 	Handler     commandHandlerFunc
// }

// // CommandParameter : Specifies the name, position and what this parameter does.
// type CommandParameter struct {
// 	Name        string
// 	Description string
// 	Position    int
// 	Required    bool
// }

// // Bot is the bot, man!
// type Bot struct {
// 	Session       *dgo.Session
// 	CmdHandler    *CommandHandler
// 	Conversations []Conversation
// 	Hub           *hub
// 	closed        bool
// }

// // CommandHandler : Holds a list of commands and specifies the Bot prefix.
// // TODO: use maps for some stuff...
// type CommandHandler struct {
// 	Commands      []*Command
// 	Prefix        string
// 	Triggers      []*Trigger
// 	Conversations []*Conversation
// }

// // Trigger takes a pattern and returns a response
// type Trigger struct {
// 	Pattern  string
// 	Response string
// }

// // Conversation between the bot and a user..
// type Conversation struct {
// 	UserID string
// }

// // TODO: phoenix-like channels?
// // type wsChannel struct {
// // 	Name string
// // 	Clients map[*wsClient]bool
// // 	Messages []*Message
// // }

// // type Message struct {
// // 	timeSent
// // 	payload string
// // }

// type catImage struct {
// 	URL string `xml:"data>images>image>url"`
// }

// type wsClient struct {
// 	send     chan []byte
// 	sendJSON chan interface{}
// 	ws       *websocket.Conn
// }

// type hub struct {
// 	// channels   map[*wsChannel]bool
// 	clients    map[*wsClient]bool
// 	broadcast  chan string
// 	register   chan *wsClient
// 	unregister chan *wsClient
// 	content    string
// }

// type commandHandlerFunc func(s *dgo.Session, m *dgo.MessageCreate, c *Command)

// // type serverStats struct {
// // 	SentMessages int
// // 	UsersJoined  int
// // }

// // type Bot struct {
// // 	Session *dgo.Session
// // 	Hub *BotHub
// // }

// // type BotHub struct {}

// // type Channel struct {}
// // type Client struct {}

// // type CommandHandler struct {}
// // type Command struct {}
// // type CommandParam struct {}
