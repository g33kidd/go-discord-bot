package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	dgo "github.com/bwmarrin/discordgo"
)

func setupCommands(ch *CommandHandler) {
	testCommand := &Command{
		Signature:   "test",
		Description: "Does a thing!",
		Handler:     testCommandHandler,
	}

	concurrencyTest := &Command{
		Signature:   "c",
		Description: "Testing some goroutine stuff!",
		Handler:     concurrencyTestHandler,
	}

	twitchCommand := &Command{
		Signature:   "twitch-edit",
		Description: "Does another thing!",
		Handler:     twitchChannelEditCommand,
	}

	pingCommand := &Command{
		Signature:   "ping",
		Description: "Ping pong!",
		Handler:     pingPongHandler,
	}

	pongCommand := &Command{
		Signature:   "pong",
		Description: "Ping pong!",
		Handler:     pingPongHandler,
	}

	randomCatCommand := &Command{
		Signature:   "cat",
		Description: "Gives a random cat image!",
		Handler:     catCommandHandler,
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

	ch.AddCommand(testCommand)
	ch.AddCommand(twitchCommand)
	ch.AddCommand(concurrencyTest)
	ch.AddCommand(randomCatCommand)
	ch.AddCommand(pingCommand)
	ch.AddCommand(pongCommand)
}

func testCommandHandler(s *dgo.Session, m *dgo.MessageCreate, c *Command) {
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

func getChannel() *twitchChannel {
	client := &http.Client{}
	// TODO: pass in the channel id after getting the channel id!
	req, err := http.NewRequest("GET", "https://api.twitch.tv/kraken/channel", nil)
	req.Header.Add("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Add("Authorization", "OAuth "+os.Getenv("TWITCH_OAUTH"))
	req.Header.Add("Client-ID", os.Getenv("TWITCH_CLIENT_ID"))
	resp, err := client.Do(req)
	if err != nil {
		// TODO: dont panic!
		panic(err)
	}
	defer resp.Body.Close()
	ch := &twitchChannel{}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, ch)
	return ch
}

func updateChannel(channelID string, game string, status string) error {
	update := &channelUpdate{
		Channel: channel{
			Status: status,
			Game:   game,
		},
	}

	jb, _ := json.Marshal(update)

	client := &http.Client{}
	req, err := http.NewRequest("PUT", "https://api.twitch.tv/kraken/channels/"+channelID, bytes.NewBuffer(jb))
	req.Header.Add("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Add("Authorization", "OAuth "+os.Getenv("TWITCH_OAUTH"))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Client-ID", os.Getenv("TWITCH_CLIENT_ID"))

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func twitchChannelEditCommand(s *dgo.Session, m *dgo.MessageCreate, c *Command) {

	ch := getChannel()
	fmt.Println(ch.Status)
	fmt.Println(ch.Game)

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

	if update := updateChannel(ch.ID, game, title); update != nil {
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
func concurrencyTestHandler(s *dgo.Session, m *dgo.MessageCreate, c *Command) {
	for i := 0; i < 5; i++ {
		c := timer(3 * time.Second)
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%d", i))
		<-c
	}
	s.ChannelMessageSend(m.ChannelID, "ch finished")
}

func catCommandHandler(s *dgo.Session, m *dgo.MessageCreate, c *Command) {
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

func pingPongHandler(s *dgo.Session, m *dgo.MessageCreate, c *Command) {
	if c.Signature == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}

// TODO: Access to commandHandler in here? Yes/no?
// func conversationHandler(s *dgo.Session, m *dgo.MessageCreate, c *Command) {
// 		ch.AddConversation()
// }
