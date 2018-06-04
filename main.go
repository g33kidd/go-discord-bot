package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

var cmdHandler *CommandHandler

const prefix = "$"

// TODO: Waiting for response system?
// 	User: $derp
//	Bot: What derp? <waits for another message from User>
// 	User: Do the thing!
//	Bot: runs a response function or checks to see if this user is in the queue and response accordingly.

// TODO: Look into a different way of setting up commands!
// 	twitchCommand := Command.New("twitch-edit", twitchChannelEditCommand)
// 	twitchCommand.AddParam("title", "Sets the title for the twitch stream", 1, true)

func main() {

	// Test Command
	testCommand := &Command{
		Signature:   "test",
		Description: "Does a thing!",
		Handler:     testCommandHandler,
	}

	testCommand.AddParam(&CommandParameter{
		Name:        "name",
		Description: "Sets the name of test",
		Position:    0,
		Required:    true,
	})

	testCommand.AddParam(&CommandParameter{
		Name:        "something",
		Description: "Does the something",
		Position:    1,
		Required:    true,
	})

	// Edit Twitch Channel
	twitchCommand := &Command{
		Signature:   "twitch-edit",
		Description: "Does another thing!",
		Handler:     twitchChannelEditCommand,
	}

	twitchCommand.AddParam(&CommandParameter{
		Name:        "game",
		Description: "Sets the game for the twitch stream",
		Position:    0,
		Required:    true,
	})

	twitchCommand.AddParam(&CommandParameter{
		Name:        "title",
		Description: "Sets the title for the twitch stream",
		Position:    1,
		Required:    true,
	})

	// Setup the command handler and add some commands!
	cmdHandler = NewCommandHandler("$")
	cmdHandler.AddCommand(testCommand)
	cmdHandler.AddCommand(twitchCommand)

	// Load our .env file in development.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
		log.Fatalln(err)
		return
	}

	// Get the bot token and setup the Session
	botToken := os.Getenv("BOT_TOKEN")
	dg, err := dgo.New("Bot " + botToken)
	if err != nil {
		log.Fatalln("Error setting up discord.")
		log.Fatalln(err)
		return
	}

	// Setup the discordgo event handlers!
	dg.AddHandlerOnce(ready)
	dg.AddHandler(messageCreate)

	// Open the connection to discord.
	err = dg.Open()
	if err != nil {
		log.Fatalln("Error opening connection to Discord", err)
		return
	}

	// wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running.   Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// cleanly close down the discord connection
	dg.Close()

}

func ready(s *dgo.Session, r *dgo.Ready) {
	fmt.Println("Bot is ready to Go!")
}

func logMessageCreate(s *dgo.Session, m *dgo.MessageCreate) {
	fmt.Printf("%-15s%s\n", m.Author.Username)
}

// TODO: add some fun colorized logging for Discord events.
// TODO: add some fun colorized logging for Discord events.
// TODO: add some fun colorized logging for Discord events.
// TODO: add some fun colorized logging for Discord events.
// TODO: add some fun colorized logging for Discord events.
// TODO: add some fun colorized logging for Discord events.
// TODO: add some fun colorized logging for Discord events.
// TODO: add some fun colorized logging for Discord events.
func messageCreate(s *dgo.Session, m *dgo.MessageCreate) {
	/// Do not respond to self, or any other bot messages.
	if s.State.User.ID == m.Author.ID || m.Author.Bot {
		return
	}

	logMessageCreate(s, m)

	// Try to find the command the user is trying to use, if it exists in the message.
	cmd, err := cmdHandler.FindCommand(m.Content, true)
	if err != nil {
		fmt.Println("command not found!")
	} else {
		cmd.Handler(s, m, cmd)
		return
	}

	res, err := cmdHandler.MaybeHandleCodeBlock(s, m)
	if err != nil {
		fmt.Println("was not a code block!")
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```\n%s\n```", res))
		return
	}
}
