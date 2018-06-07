package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/websocket"

	"github.com/joho/godotenv"
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

const prefix = "$"

func main() {

	// Load our .env file in development.
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file.")
		log.Fatalln(err)
		return
	}

	discordBot := &Bot{}
	discordBot.Setup()

	http.HandleFunc("/", discordBot.serveDashboard)
	http.HandleFunc("/ws", discordBot.serveWs)

	// go discordBot.runWs()
	go discordBot.Connect()
	go http.ListenAndServe(":8000", nil)
	go discordBot.ReportStatus()

	// wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// cleanly close down the discord connection
	discordBot.Close()
}
