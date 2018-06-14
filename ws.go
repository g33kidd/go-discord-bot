package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/g33kidd/n00b/discord"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 1024 * 1024
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  maxMessageSize,
	WriteBufferSize: maxMessageSize,
}

// NOTE: all this ws stuff is temp, redo into its own package or something
type wsHandler struct {
	bot *discord.Bot
	hub *wsHub
}

type wsClient struct {
	send     chan []byte
	sendJSON chan interface{}
	ws       *websocket.Conn
}

// type wsHub struct {
// 	clients    map[*wsClient]bool
// 	broadcast  chan string
// 	register   chan *wsClient
// 	unregister chan *wsClient
// }

type wsChannel struct {
	Name       string
	Clients    map[*wsClient]bool
	register   chan *wsClient
	unregister chan *wsClient
}

type wsHub struct {
	Bot        *discord.Bot
	clients    map[*wsClient]bool
	channels   map[*wsChannel]bool
	broadcast  chan string
	register   chan *wsClient
	unregister chan *wsClient
}

// So this should handle channel connections. /ws/:chan
func (h *wsHub) handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		/// probably not the right error code
		http.Error(w, "This is not allowed.", 401)
	}

	fmt.Println(r.RequestURI)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	defer ws.Close()

	client := &wsClient{
		ws:       ws,
		send:     make(chan []byte),
		sendJSON: make(chan interface{}),
	}

	h.register <- client
	go client.readListener(h)
	client.writeListener()
}

func (c *wsClient) readListener(h *wsHub) {
	defer func() {
		h.unregister <- c
		c.ws.Close()
	}()

	c.ws.SetReadLimit(maxMessageSize)
	c.ws.SetReadDeadline(time.Now().Add(pongWait))
	c.ws.SetPongHandler(func(string) error {
		c.ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		// TODO: c.ws.ReadJSON() for when the client needs to send JSON to here...
		_, msg, err := c.ws.ReadMessage()
		if err != nil {
			break
		}

		fmt.Println("received message!")
		fmt.Printf("got message: %s\n", string(msg))

		h.broadcast <- string(msg)
	}
}

func (c *wsClient) writeListener() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.write(websocket.TextMessage, message); err != nil {
				return
			}
		case data, ok := <-c.sendJSON:
			if !ok {
				c.write(websocket.CloseMessage, []byte{})
				return
			}

			if err := c.writeJSON(data); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.write(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (h *wsHub) run(b *discord.Bot) {
	for {
		select {
		case c := <-h.register:
			fmt.Println("register")
			h.clients[c] = true
			c.send <- []byte("hello")
			break
		case c := <-h.unregister:
			fmt.Println("unregister")
			if ok := h.clients[c]; ok {
				delete(h.clients, c)
				close(c.send)
			}
		case m := <-b.MessageCreate:
			fmt.Printf("messageCreate received! %s\n", m.Author.Username)
			for c := range h.clients {
				c.sendJSON <- m.Author
			}
			break
		}
	}
}

func (c *wsClient) write(mt int, message []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteMessage(mt, message)
}

func (c *wsClient) writeJSON(data interface{}) error {
	c.ws.SetWriteDeadline(time.Now().Add(writeWait))
	return c.ws.WriteJSON(data)
}
