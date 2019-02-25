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
		Session:       session,
		Prefix:        prefix,
		CmdHandler:    cmdHandler,
		closed:        true,
		MessageCreate: make(chan *dgo.MessageCreate),
	}

	msgCreateFunc := func(s *dgo.Session, m *dgo.MessageCreate) {
		b.MessageCreate <- m
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

	context := &MessageContext{
		Bot:     b,
		Session: s,
		Message: m,
	}

	log.Println(context)

	content := m.Content

	/// Do not respond to self, or any other bot messages.
	if s.State.User.ID == m.Author.ID || m.Author.Bot {
		return
	}

	cmd, err := b.CmdHandler.FindCommand(content, true)
	if err == nil {
		// TODO: instead of H(s, m, cmd) do H(bot, m, cmd) because bot already has a session.
		context.Command = cmd
		// TODO eventually pass in the context here..
		cmd.Handler(s, m, cmd)
		return
	}

	cb, err := b.CmdHandler.MaybeHandleCodeBlock(s, m)
	if err == nil {
		s.ChannelMessageSend(m.ChannelID, cb)
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
