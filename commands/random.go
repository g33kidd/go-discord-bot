package commands

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	dc "github.com/g33kidd/n00b/discord"
)

type catImage struct {
	URL string `xml:"data>images>image>url"`
}

// PingPongCommand does ping and pong!
func PingPongCommand(ctx *dc.MessageContext) {
	// You can use GetVal to extract some information from the Context, but you can also just
	// use the values directly from the context.
	// b, m, c, _ := ctx.GetVal()
	if ctx.Command.Signature == "pong" {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Ping!")
	} else {
		ctx.Session.ChannelMessageSend(ctx.Message.ChannelID, "Pong!")
	}
}

// RandomCatCommand gives a random cat image
func RandomCatCommand(ctx *dc.MessageContext) {
	_, m, _, s := ctx.GetVal()

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

// APICommand does ping and pong!
func APICommand(ctx *dc.MessageContext) {
	_, m, c, s := ctx.GetVal()
	client := &http.Client{}
	url, err := c.GetParam(m.Content, "url")
	if err != nil {
		// TODO: better way to generate these kinds of messages!
		s.ChannelMessageSend(m.ChannelID, "You need to give me the `url` param. For more info do `$help api`!")
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
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
