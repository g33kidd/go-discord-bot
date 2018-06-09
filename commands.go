package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/g33kidd/n00b/discord"
	"github.com/g33kidd/n00b/twitch"
)

// TODO: Move commands into separate files/modules with helper functions to create them easily and add them to the bot!

func setupCommands(ch *discord.CommandHandler) {
	testCommand := &discord.Command{
		Signature:   "test",
		Description: "Does a thing!",
		Handler:     testCommandHandler,
	}

	apiCommand := &discord.Command{
		Signature:   "api",
		Description: "Sends a GET request to an API endpoint. Must be a JSON response.",
		Handler:     apiCommandHandler,
	}

	concurrencyTest := &discord.Command{
		Signature:   "c",
		Description: "Testing some goroutine stuff!",
		Handler:     concurrencyTestHandler,
	}

	twitchCommand := &discord.Command{
		Signature:   "twitch-edit",
		Description: "Does another thing!",
		Handler:     twitchChannelEditCommand,
	}

	pingCommand := &discord.Command{
		Signature:   "ping",
		Description: "Ping pong!",
		Handler:     pingPongHandler,
	}

	pongCommand := &discord.Command{
		Signature:   "pong",
		Description: "Ping pong!",
		Handler:     pingPongHandler,
	}

	randomCatCommand := &discord.Command{
		Signature:   "cat",
		Description: "Gives a random cat image!",
		Handler:     catCommandHandler,
	}

	testCommand.AddParam(&discord.CommandParameter{
		Name:        "name",
		Description: "Sets the name of test",
		Position:    0,
		Required:    true,
	})

	testCommand.AddParam(&discord.CommandParameter{
		Name:        "something",
		Description: "Does the something",
		Position:    1,
		Required:    true,
	})

	twitchCommand.AddParam(&discord.CommandParameter{
		Name:        "game",
		Description: "Sets the game for the twitch stream",
		Position:    0,
		Required:    true,
	})

	twitchCommand.AddParam(&discord.CommandParameter{
		Name:        "title",
		Description: "Sets the title for the twitch stream",
		Position:    1,
		Required:    true,
	})

	apiCommand.AddParam(&discord.CommandParameter{
		Name:        "url",
		Description: "The URL to make a request to.",
		Position:    0,
		Required:    true,
	})

	macroCommand := &discord.Command{
		Signature:   "macro",
		Description: "Defines a macro that runs one or more commands in order. Shortcut for using other commands basically. Currently only supports up to 3 commands. Ability to pipe data into the next command will be added later.",
		Handler:     macroCommandHandler,
	}

	// TODO: before finishing this command, rework how you add params and create commands!
	// This is getting quite messy!
	macroCommand.AddParam(&discord.CommandParameter{
		Name:        "name",
		Description: "Name of the macro to be ran or what it should be called!",
		Position:    0,
		Required:    true,
	})

	macroCommand.AddParam(&discord.CommandParameter{
		Name:        "cmd1",
		Description: "First command to run.",
		Position:    1,
		Required:    false,
	})

	macroCommand.AddParam(&discord.CommandParameter{
		Name:        "cmd2",
		Description: "Second command to run.",
		Position:    2,
		Required:    false,
	})

	macroCommand.AddParam(&discord.CommandParameter{
		Name:        "cmd3",
		Description: "Third command to run.",
		Position:    3,
		Required:    false,
	})

	ch.AddCommand(apiCommand)
	ch.AddCommand(testCommand)
	ch.AddCommand(twitchCommand)
	ch.AddCommand(concurrencyTest)
	ch.AddCommand(randomCatCommand)
	ch.AddCommand(pingCommand)
	ch.AddCommand(pongCommand)
	ch.AddCommand(macroCommand)
}

// TODO: More work on macro system!
func macroCommandHandler(s *dgo.Session, m *dgo.MessageCreate, c *discord.Command) {
	content := m.Content

	// TODO: Re-work this so that I can return nil for name or cmd1. Just requires returning it as a pointer.
	name, err := c.GetParam(content, "name")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, "whoops... you are missing a few params. Try `$help macro` for more information!")
		return
	}

	cmd1, _ := c.GetParam(content, "cmd1")
	// cmd2, _ := c.GetParam(content, "cmd2")
	// cmd3, _ := c.GetParam(content, "cmd3")

	if cmd1 == "" && name != "" {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Okay, running this macro called `%s`.", name))
		return
	}

	fmt.Println(cmd1)

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Okay creating this macro `%s` doing `%s`", name, cmd1))
}

// TODO: add param called 'value' to get the value of whatever json was responded! Whether it's an image or whatever.
func apiCommandHandler(s *dgo.Session, m *dgo.MessageCreate, c *discord.Command) {
	client := &http.Client{}
	url, err := c.GetParam(m.Content, "url")
	if err != nil {
		// TODO: better way to generate these kinds of messages!
		s.ChannelMessageSend(m.ChannelID, "You need to give me the `url` param. For more info do `$help api`!")
	}
	resp, err := client.Get(url)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", err))
		return
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", err))
		return
	}
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, data, "", "\t")
	if err != nil {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```%s```", err))
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("```json\n%s```", prettyJSON.Bytes()))
}

func testCommandHandler(s *dgo.Session, m *dgo.MessageCreate, c *discord.Command) {
	nameP, err := c.GetParam(m.Content, "name")
	if err != nil {
		fmt.Println(err)
	}
	s.ChannelMessageSend(m.ChannelID, nameP)

	somethingP, err := c.GetParam(m.Content, "something")
	if err != nil {
		fmt.Println(err)
	}
	s.ChannelMessageSend(m.ChannelID, somethingP)
}

func twitchChannelEditCommand(s *dgo.Session, m *dgo.MessageCreate, c *discord.Command) {

	ch := twitch.GetMyChannel()

	game, err := c.GetParam(m.Content, "game")
	if err != nil {
		fmt.Println(err)
		return
	}

	title, err := c.GetParam(m.Content, "title")
	if err != nil {
		fmt.Println(err)
		return
	}

	chUpdate := &twitch.TwitchChannelEditData{
		Game:   game,
		Status: title,
	}

	if update := twitch.UpdateChannel(ch.ID, chUpdate); update != nil {
		s.ChannelMessageSend(m.ChannelID, "Something did not go as planned!")
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("I set the game to **%s** and the title to **%s**", game, title))
	}
}

func timer(d time.Duration) <-chan int {
	c := make(chan int)
	go func() {
		time.Sleep(d)
		c <- 1
	}()
	return c
}

// TODO: come back to this!
func concurrencyTestHandler(s *dgo.Session, m *dgo.MessageCreate, c *discord.Command) {
	for i := 0; i < 5; i++ {
		c := timer(3 * time.Second)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%d", i))
		<-c
	}
	s.ChannelMessageSend(m.ChannelID, "ch finished")
}

func catCommandHandler(s *dgo.Session, m *dgo.MessageCreate, c *discord.Command) {
	resp, err := http.Get("http://thecatapi.com/api/images/get?format=xml&results_per_page=1")
	if err != nil {
		fmt.Println("Didnt get cat image!")
		s.ChannelMessageSend(m.ChannelID, "Could not get the cat image you wanted! Sorry!")
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Didnt get cat image!")
		s.ChannelMessageSend(m.ChannelID, "Could not get the cat image you wanted! Sorry!")
		return
	}

	cat := &catImage{}
	xml.Unmarshal(data, cat)

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Found the image you wanted! Here! %s", cat.URL))
}

func pingPongHandler(s *dgo.Session, m *dgo.MessageCreate, c *discord.Command) {
	if c.Signature == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}
