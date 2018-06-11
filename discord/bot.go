package discord

import (
	"fmt"
	"log"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/fatih/color"
)

// NewBot creates a new bot!
func NewBot(token string, prefix string) *Bot {
	if prefix == "" {
		prefix = "!"
	}

	session, err := dgo.New("Bot " + token)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	cmdHandler := NewCommandHandler(prefix)

	b := &Bot{
		Session:    session,
		Prefix:     prefix,
		CmdHandler: cmdHandler,
		closed:     true,
	}

	msgCreateFunc := func(s *dgo.Session, m *dgo.MessageCreate) {
		b.messageCreate(s, m)
	}

	guildMemberAddFunc := func(s *dgo.Session, m *dgo.GuildMemberAdd) {
		b.guildMemberAdd(s, m)
	}

	ready := func(s *dgo.Session, m *dgo.Ready) {
		b.ready(s, m)
	}

	// Setup the discordgo event handlers!
	session.AddHandlerOnce(ready)
	session.AddHandler(msgCreateFunc)
	session.AddHandler(guildMemberAddFunc)

	return b
}

// Connect opens the connection to the discord gateway.
func (b *Bot) Connect() {
	b.closed = false
	if err := b.Session.Open(); err != nil {
		b.closed = true
		log.Fatalln("Error opening connection to Discord", err)
		return
	}
}

// Disconnect disconnects, yo!
func (b *Bot) Disconnect() {
	b.Session.Close()
}

func (b *Bot) messageCreate(s *dgo.Session, m *dgo.MessageCreate) {
	logMessageCreate(s, m)

	// NOTE: This is for the websocket stuff ->
	// for c := range b.Hub.clients {
	// 	c.sendJSON <- m.Author
	// }

	/// Do not respond to self, or any other bot messages.
	content := m.Content

	if s.State.User.ID == m.Author.ID || m.Author.Bot {
		return
	}

	cmd, err := b.CmdHandler.FindCommand(content, true)
	if err == nil {
		cmd.Handler(s, m, cmd)
		return
	}

	cb, err := b.CmdHandler.MaybeHandleCodeBlock(s, m)
	if err == nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```\n%s\n```\nis the result", cb))
		return
	}

	trig, err := b.CmdHandler.MaybeHandleMessageTrigger(s, m)
	if err == nil {
		s.ChannelMessageSend(m.ChannelID, trig)
		return
	}
}

// TODO: what to do with this one?
func (b *Bot) ready(s *dgo.Session, r *dgo.Ready) {
}

func (b *Bot) guildMemberAdd(s *dgo.Session, m *dgo.GuildMemberAdd) {
	fmt.Println("somebody joined the guild!")
}

func logMessageCreate(s *dgo.Session, m *dgo.MessageCreate) {
	ch, err := s.State.Channel(m.ChannelID)
	if err != nil {
		fmt.Println("didnt find the channel!")
		return
	}

	g, err := s.State.Guild(ch.GuildID)
	if err != nil {
		fmt.Println("didnt find the guild!")
		return
	}

	fmt.Printf(
		"[%s #%s] @%s : %s\n",
		color.GreenString(g.Name),
		color.BlueString(ch.Name),
		color.MagentaString(m.Author.Username),
		m.ContentWithMentionsReplaced(),
	)
}

// package main

// import (
// 	"fmt"
// 	"io/ioutil"
// 	"log"
// 	"net/http"
// 	"os"
// 	"time"

// 	dgo "github.com/bwmarrin/discordgo"
// )

// // TODO: Cleanup the websocket stuff!
// // TODO: command aliases!
// // TODO: Setup permission based commands.
// // TODO: Waiting for response system?
// // TODO: Look into a different way of setting up commands!
// // TODO: Better system for triggers. More variety in responses!

// // Setup initializes the session
// func (b *Bot) Setup() {
// 	// Get the bot token and setup the Session
// 	botToken := os.Getenv("BOT_TOKEN")
// 	dg, err := dgo.New("Bot " + botToken)
// 	if err != nil {
// 		log.Fatalln("Error setting up discord.")
// 		log.Fatalln(err)
// 		return
// 	}

// 	b.Hub = &hub{
// 		broadcast:  make(chan string),
// 		register:   make(chan *wsClient),
// 		unregister: make(chan *wsClient),
// 		clients:    make(map[*wsClient]bool),
// 		content:    "",
// 	}

// 	b.CmdHandler = NewCommandHandler("$")

// 	msgCreateFunc := func(s *dgo.Session, m *dgo.MessageCreate) {
// 		b.messageCreate(s, m)
// 	}

// 	guildMemberAddFunc := func(s *dgo.Session, m *dgo.GuildMemberAdd) {
// 		b.guildMemberAdd(s, m)
// 	}

// 	ready := func(s *dgo.Session, m *dgo.Ready) {
// 		b.ready(s, m)
// 	}

// 	// Setup the discordgo event handlers!
// 	b.Session = dg
// 	b.Session.AddHandlerOnce(ready)
// 	b.Session.AddHandler(msgCreateFunc)
// 	b.Session.AddHandler(guildMemberAddFunc)
// 	setupCommands(b.CmdHandler)
// }

// func (b *Bot) serveDashboard(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "GET" {
// 		http.Error(w, "Method not supported", 405)
// 		return
// 	}

// 	f, err := ioutil.ReadFile("index.html")
// 	if err != nil {
// 		fmt.Println("Could not open the file.", err)
// 	}
// 	fmt.Fprintf(w, "%s", f)
// }

// func (b *Bot) serveWs(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "GET" {
// 		http.Error(w, "Method not supported", 405)
// 		return
// 	}

// 	ws, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println(err)
// 		return
// 	}

// 	defer ws.Close()

// 	client := &wsClient{
// 		send:     make(chan []byte, maxMessageSize),
// 		sendJSON: make(chan interface{}, maxMessageSize),
// 		ws:       ws,
// 	}

// 	b.Hub.register <- client
// 	go client.writeListener()
// 	client.readListener(b)
// }

// // Connect connects to the discord!
// func (b *Bot) Connect() {
// 	// Open the connection to discord.
// 	b.closed = false
// 	go b.Hub.run()
// 	if err := b.Session.Open(); err != nil {
// 		log.Fatalln("Error opening connection to Discord", err)
// 		return
// 	}
// }

// // Close closes the thing!
// func (b *Bot) Close() {
// 	b.closed = true
// 	b.Session.Close()
// }

// // RunHub waits for messages from the Hub channels
// func (h *hub) run() {
// 	for {
// 		select {
// 		case client := <-h.register:
// 			h.clients[client] = true
// 			client.send <- []byte(h.content)
// 			break
// 		case c := <-h.unregister:
// 			if ok := h.clients[c]; ok {
// 				delete(h.clients, c)
// 				close(c.send)
// 			}
// 		}
// 	}
// }

// // ReportStatus reports the status to the websocket clients
// // TODO: Find other data to report to wsClients
// func (b *Bot) ReportStatus() {
// 	for {
// 		if b.closed {
// 			break
// 		}

// 		if len(b.Hub.clients) <= 0 {
// 			time.Sleep(30 * time.Second)
// 			return
// 		}

// 		for c := range b.Hub.clients {
// 			c.sendJSON <- b.Session.State.Presences
// 		}

// 		time.Sleep(10 * time.Second)
// 	}
// }

// func (b *Bot) messageCreate(s *dgo.Session, m *dgo.MessageCreate) {
// 	logMessageCreate(s, m)

// 	for c := range b.Hub.clients {
// 		c.sendJSON <- m.Author
// 	}

// 	/// Do not respond to self, or any other bot messages.
// 	content := m.Content

// 	if s.State.User.ID == m.Author.ID || m.Author.Bot {
// 		return
// 	}

// 	cmd, err := b.CmdHandler.FindCommand(content, true)
// 	if err == nil {
// 		cmd.Handler(s, m, cmd)
// 		return
// 	}

// 	cb, err := b.CmdHandler.MaybeHandleCodeBlock(s, m)
// 	if err == nil {
// 		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```\n%s\n```\nis the result", cb))
// 		return
// 	}

// 	trig, err := b.CmdHandler.MaybeHandleMessageTrigger(s, m)
// 	if err == nil {
// 		s.ChannelMessageSend(m.ChannelID, trig)
// 		return
// 	}
// }

// func (b *Bot) ready(s *dgo.Session, r *dgo.Ready) {
// }

// func (b *Bot) guildMemberAdd(s *dgo.Session, m *dgo.GuildMemberAdd) {
// 	fmt.Println("somebody joined the guild!")
// }

// // var channels = []string{"g33kidd", "neoplatonist", "pixelogicdev", "naysayer88"}

// // func (b *Bot) watchTwitchChannels() {
// // 	for {
// // 		users := twitch.GetUsers(channels)
// // 		for _, u := range users {
// // 			s := twitch.GetStream(u.ID)
// // 			if s.Data == nil {
// // 				return
// // 			}

// // 			fmt.Printf("Channel: %s\nViewers: %v\nPreview: %s\n---\n", u.Name, s.Data.Viewers, s.Data.Preview.Small)
// // 			fmt.Println(s.Data)
// // 		}
// // 		time.Sleep(120 * time.Second)
// // 	}
// // }
