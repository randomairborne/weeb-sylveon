package love

import (
	"github.com/bwmarrin/discordgo"
	"github.com/randomairborne/eevee/api"
	"regexp"
)

type Module struct {
	api.Module
}

var userRegex, _ = regexp.Compile("/<@!?(\\d{17,19})>/g")

type nekosApiResponseJson struct {
	URL string `json:"url"`
}

func (*Module) Load(_ *discordgo.Session) {
	api.RegisterCommand("hug", RunHugCommand)
	api.RegisterCommand("kiss", RunKissCommand)
	api.RegisterCommand("pat", RunPatsCommand)
	api.RegisterCommand("slap", RunSlapCommand)
	api.RegisterIntentNeed(discordgo.IntentsGuildMessages, discordgo.IntentsDirectMessages)
}
