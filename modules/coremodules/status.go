package coremodules

import (
	"github.com/bwmarrin/discordgo"
	"github.com/randomairborne/eevee/api"
	"github.com/randomairborne/eevee/api/logger"
	"io/fs"
	"io/ioutil"
	"strings"
)

func RunStatusCommand(ds *discordgo.Session, mc *discordgo.MessageCreate, _ string, args []string) {
	if mc.Author.ID != "861733561463603240" {
		err := api.SendWithSelfDelete(ds, mc.ChannelID, 10, "You don't have permission to run that command!")
		if err != nil {
			logger.Err().Println(err.Error())
		}
		err = ds.UpdateGameStatus(0, strings.Join(args, " "))
		if err != nil {
			err := api.SendWithSelfDelete(ds, mc.ChannelID, 15, "Failed to update status: "+err.Error())
			if err != nil {
				return
			}
		}
		err = ioutil.WriteFile("status.text", []byte(strings.Join(args, " ")), fs.FileMode(0777))
		if err != nil {
			err := api.SendWithSelfDelete(ds, mc.ChannelID, 15, "Failed to write to file: "+err.Error())
			if err != nil {
				return
			}
		}
		_, err = ds.ChannelMessageSendReply(mc.ChannelID, "Changed status to "+strings.Join(args, " "), mc.Reference())
		if err != nil {
			return
		}
	}
}
