package hugs

import (
	"encoding/json"
	"github.com/bwmarrin/discordgo"
	"github.com/randomairborne/eevee/api"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
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
	api.RegisterIntentNeed(discordgo.IntentsGuildMessages, discordgo.IntentsDirectMessages)
}

func RunHugCommand(session *discordgo.Session, message *discordgo.MessageCreate, _ string, args []string) {
	resp, err := http.Get("https://nekos.life/api/v2/img/hug")
	if err != nil {
		err := api.SendWithSelfDelete(session, message.ChannelID, 10, "Failed to query API, error: `"+err.Error()+"`")
		if err != nil {
			println(err.Error())
		}
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		println(err.Error())
	}
	var response nekosApiResponseJson
	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	embed := &discordgo.MessageEmbed{
		Description: "Aww, " + message.Author.Mention() + " hugs you!",
		Image: &discordgo.MessageEmbedImage{
			URL: response.URL,
		},
	}
	var send *discordgo.MessageSend
	if len(args) > 0 {
		// check if valid mention
		if !userRegex.MatchString(args[0]) {
			err := api.SendWithSelfDelete(session, message.ChannelID, 10, "That's not a valid mention!")
			if err != nil {
				return
			}
		}
		if err != nil {
			err := api.SendWithSelfDelete(session, message.ChannelID, 10, "Regex matching failed, `"+err.Error()+"`")
			if err != nil {
				return
			}
		}
		// This gets the users ID so it'll ping them
		user := strings.Trim(args[0], "&@!<>")
		send = &discordgo.MessageSend{
			Content: args[0],
			Embed:   embed,
			AllowedMentions: &discordgo.MessageAllowedMentions{
				Users: []string{user},
			},
		}
	} else {
		send = &discordgo.MessageSend{
			Embed:           embed,
			AllowedMentions: &discordgo.MessageAllowedMentions{},
		}
	}

	_, err = session.ChannelMessageSendComplex(message.ChannelID, send)
	if err != nil {
		println(err.Error())
	}
}
