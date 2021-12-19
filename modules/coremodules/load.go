package coremodules

import (
	"github.com/bwmarrin/discordgo"
	"github.com/randomairborne/eevee/api"
)

type Module struct {
	api.Module
}

func (*Module) Load(_ *discordgo.Session) {
	api.RegisterCommand("ping", RunPingCommand)
	api.RegisterCommand("status", RunStatusCommand)

	api.RegisterIntentNeed(discordgo.IntentsGuildMessages, discordgo.IntentsDirectMessages)
}
