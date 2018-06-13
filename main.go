package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	cmds "github.com/g33kidd/n00b/commands"
	"github.com/g33kidd/n00b/services"
	"github.com/gorilla/websocket"

	"github.com/g33kidd/n00b/discord"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/joho/godotenv"
)

// TODO: get rid of this and just use the prefix set by the guild.
// By default a guild prefix should be ! or $.
const prefix = "$"

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
}

type wsClient struct {
	send     chan []byte
	sendJSON chan interface{}
	ws       *websocket.Conn
}

type wsHub struct {
	clients    map[*wsClient]bool
	broadcast  chan string
	register   chan *wsClient
	unregister chan *wsClient
}

var hub *wsHub

func main() {

	// Load our .env file in development.
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Error loading .env file.", err)
		return
	}

	bot := discord.NewBot(os.Getenv("BOT_TOKEN"), prefix)
	if bot == nil {
		log.Fatalln("Error setting up discord Bot.")
		return
	}

	cmds.RegisterRandomCommands(bot)
	cmds.RegisterTwitchCommands(bot)
	cmds.RegisterFunCommands(bot)
	cmds.RegisterImageCommands(bot)
	cmds.RegisterUtilityCommands(bot)
	cmds.RegisterTestingCommands(bot)

	hub = &wsHub{
		register:   make(chan *wsClient),
		unregister: make(chan *wsClient),
		broadcast:  make(chan string),
		clients:    make(map[*wsClient]bool),
	}

	wsh := &wsHandler{bot}

	http.HandleFunc("/", serveDashboard)
	http.HandleFunc("/ws", wsh.handle)

	go hub.run(bot)
	go http.ListenAndServe(":8000", nil)
	go bot.Connect()
	go services.TwitchLiveAlerts(bot)

	// wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// cleanly close down the discord connection
	bot.Disconnect()
}

func serveDashboard(w http.ResponseWriter, r *http.Request) {
	f, err := ioutil.ReadFile("./client/public/index.html")
	if err != nil {
		fmt.Println("could not read file.", err)
		return
	}

	fmt.Fprintf(w, "%s", f)
}

func (wsh *wsHandler) handle(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	client := &wsClient{
		ws:       ws,
		send:     make(chan []byte),
		sendJSON: make(chan interface{}),
	}
	defer ws.Close()

	hub.register <- client
	go client.readListener()
	client.writeListener()
}

func (c *wsClient) readListener() {
	defer func() {
		hub.unregister <- c
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

		hub.broadcast <- string(msg)
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
