package commands

import (
	dc "github.com/g33kidd/n00b/discord"
)

// RegisterTwitchCommands registers the twitch command group
func RegisterTwitchCommands(bot *dc.Bot) {
	twitchEdit := dc.NewCommand("twitchedit", "Edits twitch channel info!", TwitchChannelEditCommand)
	twitchEdit.AddParameter("game", "The new game for twitch channel.", true)
	twitchEdit.AddParameter("title", "The new title for twitch channel.", true)

	bot.CmdHandler.AddCommand(twitchEdit)
}

// RegisterFunCommands registers the fun command group
func RegisterFunCommands(bot *dc.Bot) {

}

// RegisterImageCommands registers the fun command group
func RegisterImageCommands(bot *dc.Bot) {

}

// RegisterUtilityCommands registers the fun command group
func RegisterUtilityCommands(bot *dc.Bot) {

}

// RegisterTestingCommands for testing stuff
func RegisterTestingCommands(bot *dc.Bot) {
	macro := dc.NewCommand("macro", "Defines a macro. WIP", MacroCommand)
	macro.AddParameter("name", "Name of macro.", true)
	macro.AddParameter("cmd1", "cmd1", false)
	macro.AddParameter("cmd2", "cmd2", false)
	macro.AddParameter("cmd3", "cmd3", false)

	bot.CmdHandler.AddCommand(macro)
}

// RegisterRandomCommands registers the fun command group
func RegisterRandomCommands(bot *dc.Bot) {
	ping := dc.NewCommand("ping", "Ping!", PingPongCommand)
	pong := dc.NewCommand("pong", "Pong!", PingPongCommand)

	randomCat := dc.NewCommand("cat", "Random cat anyone?", RandomCatCommand)

	api := dc.NewCommand("api", "Makes a GET request to a JSON API and shows the content.", ApiCommand)
	api.AddParameter("url", "The URL to make a request to.", true)

	bot.CmdHandler.AddCommand(api)
	bot.CmdHandler.AddCommand(randomCat)
	bot.CmdHandler.AddCommand(ping)
	bot.CmdHandler.AddCommand(pong)
}
