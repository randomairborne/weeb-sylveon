package api

import (
	"github.com/bwmarrin/discordgo"
	"github.com/randomairborne/eevee/api/logger"
	"time"
)

func GetGuild(ds *discordgo.Session, guildId string) *discordgo.Guild {
	g, err := ds.State.Guild(guildId)
	if err != nil || g.Name == "" {
		// Try fetching via REST API
		g, err = ds.Guild(guildId)
		if err != nil {
			logger.Err().Printf("unable to fetch Guild for Message, %s", err)
		} else {
			// Attempt to add this channel into our State
			err = ds.State.GuildAdd(g)
			if err != nil {
				logger.Err().Printf("error updating Guild with Channel, %s", err)
			}
		}
	}

	return g
}

func GetChannel(ds *discordgo.Session, channelId string) *discordgo.Channel {
	c, err := ds.State.Channel(channelId)
	if err != nil {
		// Try fetching via REST API
		c, err = ds.Channel(channelId)
		if err != nil {
			logger.Err().Printf("unable to fetch Channel for Message, %s", err)
		} else {
			// Attempt to add this channel into our State
			err = ds.State.ChannelAdd(c)
			if err != nil {
				logger.Err().Printf("error updating State with Channel, %s", err)
			}
		}
	}

	return c
}
func SendWithSelfDelete(ds *discordgo.Session, channelId string, duration int, message string) error {
	m, err := ds.ChannelMessageSend(channelId, message)
	if err != nil {
		return err
	}

	go func(ch, id string, session *discordgo.Session) {
		<-time.After(time.Duration(duration) * time.Second)
		_ = ds.ChannelMessageDelete(channelId, m.ID)
	}(channelId, m.ID, ds)
	return nil
}
