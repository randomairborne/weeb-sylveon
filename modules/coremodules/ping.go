package coremodules

import (
	"github.com/bwmarrin/discordgo"
	"strconv"
)

func RunPingCommand(ds *discordgo.Session, mc *discordgo.MessageCreate, _ string, _ []string) {
	_, err := ds.ChannelMessageSendReply(
		mc.ChannelID,
		strconv.FormatInt(ds.HeartbeatLatency().Milliseconds(), 10)+"ms",
		mc.Reference())
	if err != nil {
		println(err.Error())
	}
}
