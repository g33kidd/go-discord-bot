package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	dgo "github.com/bwmarrin/discordgo"
	"github.com/g33kidd/n00b/discord"
	"github.com/g33kidd/n00b/twitch"
)

func setupCommands(ch *discord.CommandHandler) {
	testCommand := &discord.Command{
		Signature:   "test",
		Description: "Does a thing!",
		Handler:     testCommandHandler,
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

	ch.AddCommand(testCommand)
	ch.AddCommand(twitchCommand)
	ch.AddCommand(concurrencyTest)
	ch.AddCommand(randomCatCommand)
	ch.AddCommand(pingCommand)
	ch.AddCommand(pongCommand)
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
