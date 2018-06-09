package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/g33kidd/n00b/discord"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/joho/godotenv"
)

const prefix = "$"

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

	go bot.Connect()

	// wait here until CTRL-C or other term signal is received.
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// cleanly close down the discord connection
	bot.Disconnect()
}
