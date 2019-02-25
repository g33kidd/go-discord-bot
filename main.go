package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	cmds "github.com/g33kidd/n00b/commands"

	"github.com/g33kidd/n00b/discord"
	_ "github.com/jinzhu/gorm/dialects/sqlite"

	"github.com/joho/godotenv"
)

// TODO: get rid of this and just use the prefix set by the guild.
// By default a guild prefix should be ! or $.
const prefix = "$"

func main() {

	// Load our .env file in development.
	if err := godotenv.Load(); err != nil {
		log.Fatalln("Error loading .env file.", err)
		return
	}

	// Creates a new bot instance of DiscordGo
	bot := discord.NewBot(os.Getenv("BOT_TOKEN"), prefix)
	if bot == nil {
		log.Fatalln("Error setting up discord Bot.")
		return
	}

	// Register all the different commands for the bot...
	// TODO is there a better way to package all these commands without doing something like this?
	cmds.RegisterRandomCommands(bot)
	cmds.RegisterTwitchCommands(bot)
	cmds.RegisterFunCommands(bot)
	cmds.RegisterImageCommands(bot)
	cmds.RegisterUtilityCommands(bot)
	cmds.RegisterTestingCommands(bot)
	cmds.RegisterSpacexCommands(bot)

	// Dashboard stuff for handling real-time updates and the admin dashboard.
	hub := &wsHub{
		Bot:        bot,
		register:   make(chan *wsClient),
		unregister: make(chan *wsClient),
		broadcast:  make(chan string),
		clients:    make(map[*wsClient]bool),
	}

	http.HandleFunc("/", serveDashboard)
	http.HandleFunc("/ws", hub.handle)

	// All the goroutines, this stuff is running at the same time btw..
	go hub.run(bot)
	go http.ListenAndServe(":8000", nil)
	go bot.Connect()

	// TODO this shit can be turned off after starting the bot once..
	// Need a database before keeping this all the time.
	// go services.TwitchLiveAlerts(bot)

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
