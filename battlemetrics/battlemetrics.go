package battlemetrics

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// This is just for fun

// MakeRequest makes a request to battlemetrics APIs
func MakeRequest(endpoint string) (data []byte, err error) {
	data = []byte{}
	client := &http.Client{}
	resp, err := client.Get("https://api.battlemetrics.com/" + endpoint)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	data = b
	return
}

// GetServer gets a server from the API.
func GetServer(serverID string) *Server {
	r, err := MakeRequest("servers/" + serverID)
	if err != nil {
		log.Fatalln(err)
		return nil
	}

	server := &Server{}
	json.Unmarshal(r, server)
	return server
}
