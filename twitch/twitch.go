package twitch

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// MakeRequest makes a request using the proper credentials to the Twitch API v5
func MakeRequest(m string, e string, d io.Reader) (data []byte, err error) {

	client := &http.Client{}
	req, err := http.NewRequest(m, "https://api.twitch.tv/kraken/"+e, d)
	req.Header.Add("Accept", "application/vnd.twitchtv.v5+json")
	req.Header.Add("Authorization", "OAuth "+os.Getenv("TWITCH_OAUTH"))
	req.Header.Add("Client-ID", os.Getenv("TWITCH_CLIENT_ID"))
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	data, err = ioutil.ReadAll(resp.Body)
	return
}

// GetUser gets a single user based on the login
// TODO: check if there are users..
func GetUser(login string) *TwitchUser {
	u := &TwitchUsers{}
	data, err := MakeRequest("GET", "users?login="+login, nil)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, u)
	return u.Users[0]
}

// GetUsers gets a list of users
func GetUsers(logins []string) []*TwitchUser {
	loginsStr := strings.Join(logins, ",")
	u := &TwitchUsers{}
	data, err := MakeRequest("GET", "users?login="+loginsStr, nil)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, u)
	return u.Users
}

// GetStream gets the stream for a channel if it is online.
func GetStream(channelID string) *TwitchStreamData {
	s := &TwitchStream{}
	data, err := MakeRequest("GET", "streams/"+channelID, nil)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, s)
	return s.Data
}

// GetMyChannel gets the channel from OAuth token
func GetMyChannel() *TwitchChannel {
	c := &TwitchChannel{}
	data, err := MakeRequest("GET", "channel", nil)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, c)
	return c
}

// GetChannel returns a TwitchChannel given a channelID
func GetChannel(channelID string) *TwitchChannel {
	c := &TwitchChannel{}
	data, err := MakeRequest("GET", "channels/"+channelID, nil)
	if err != nil {
		panic(err)
	}
	json.Unmarshal(data, c)
	return c
}

// UpdateChannel updates the twitch channel
func UpdateChannel(channelID string, channelUpdate *TwitchChannelEditData) error {
	update := &TwitchChannelEdit{
		Channel: channelUpdate,
	}

	data, _ := json.Marshal(update)
	data, err := MakeRequest("PUT", "channels/"+channelID, bytes.NewBuffer(data))
	if err != nil {
		return err
	}

	return nil
}
