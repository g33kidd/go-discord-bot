package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	dgo "github.com/bwmarrin/discordgo"
)

// TODO: command aliases!
// TODO: Setup permission based commands.
// TODO: Waiting for response system?
// 	User: $derp
//	Bot: What derp? <waits for another message from User>
// 	User: Do the thing!
//	Bot: runs a response function or checks to see if this user is in the queue and response accordingly.
// TODO: Look into a different way of setting up commands!
// 	twitchCommand := Command.New("twitch-edit", twitchChannelEditCommand)
// 	twitchCommand.AddParam("title", "Sets the title for the twitch stream", 1, true)
// TODO: Better system for triggers. More variety in responses!

// Setup initializes the session
func (b *Bot) Setup() {
	// Get the bot token and setup the Session
	botToken := os.Getenv("BOT_TOKEN")
	dg, err := dgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalln("Error setting up discord.")
		log.Fatalln(err)
		return
	}

	b.Hub = &hub{
		broadcast:  make(chan string),
		register:   make(chan *wsClient),
		unregister: make(chan *wsClient),
		clients:    make(map[*wsClient]bool),
		content:    "",
	}

	b.CmdHandler = NewCommandHandler("$")

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
	b.Session = dg
	b.Session.AddHandlerOnce(ready)
	b.Session.AddHandler(msgCreateFunc)
	b.Session.AddHandler(guildMemberAddFunc)
	setupCommands(b.CmdHandler)
}

func (b *Bot) serveDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not supported", 405)
		return
	}

	f, err := ioutil.ReadFile("index.html")
	if err != nil {
		fmt.Println("Could not open the file.", err)
	}
	fmt.Fprintf(w, "%s", f)
}

func (b *Bot) serveWs(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not supported", 405)
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	defer ws.Close()

	client := &wsClient{
		send:     make(chan []byte, maxMessageSize),
		sendJSON: make(chan interface{}, maxMessageSize),
		ws:       ws,
	}

	b.Hub.register <- client
	go client.writeListener()
	client.readListener(b)
}

// Connect connects to the discord!
func (b *Bot) Connect() {
	// Open the connection to discord.
	b.closed = false
	go b.Hub.run()
	if err := b.Session.Open(); err != nil {
		log.Fatalln("Error opening connection to Discord", err)
		return
	}
}

// Close closes the thing!
func (b *Bot) Close() {
	b.closed = true
	b.Session.Close()
}

// RunHub waits for messages from the Hub channels
func (h *hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			client.send <- []byte(h.content)
			break
		case c := <-h.unregister:
			if ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
		}
	}
}

// ReportStatus reports the status to the websocket clients
// TODO: Find other data to report to wsClients
func (b *Bot) ReportStatus() {
	for {
		if b.closed {
			break
		}

		if len(b.Hub.clients) <= 0 {
			time.Sleep(30 * time.Second)
			return
		}

		for c := range b.Hub.clients {
			c.sendJSON <- b.Session.State.Presences
		}

		time.Sleep(10 * time.Second)
	}
}

func (b *Bot) messageCreate(s *dgo.Session, m *dgo.MessageCreate) {
	logMessageCreate(s, m)

	for c := range b.Hub.clients {
		c.sendJSON <- m.Author
	}

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

func (b *Bot) ready(s *dgo.Session, r *dgo.Ready) {
}

func (b *Bot) guildMemberAdd(s *dgo.Session, m *dgo.GuildMemberAdd) {
	fmt.Println("somebody joined the guild!")
}
