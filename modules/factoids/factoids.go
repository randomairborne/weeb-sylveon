package factoids

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/jinzhu/gorm"
	"github.com/lordralex/absol/api"
	"github.com/lordralex/absol/api/database"
	"github.com/lordralex/absol/api/logger"
	"github.com/spf13/viper"
	"strings"
)

type Module struct {
	api.Module
}

func (*Module) Load(ds *discordgo.Session) {
	//api.RegisterCommand("", RunCommand)
	api.RegisterCommand("f", RunCommand)
	api.RegisterCommand("factoid", RunCommand)
}

func RunCommand(ds *discordgo.Session, mc *discordgo.MessageCreate, cmd string, args []string) {
	if len(args) == 0 {
		return
	}

	factoids := make([]string, 0)
	if cmd == "" {
		factoids = []string{cmd}
	}

	if len(mc.MentionRoles) + len(mc.MentionChannels) > 0 {
		_, _ = ds.ChannelMessageSend(mc.ChannelID, "Cannot mention to roles or channels")
		return
	}

	for _, v := range args {
		skip := false
		for _, m := range mc.Mentions {
			if "<@!" + m.ID + ">" == v || "<@" + m.ID + ">" ==v {
				skip = true
				break
			}
		}
		if !skip {
			factoids = append(factoids, v)
		}
	}

	max := viper.GetInt("factoids.max")
	if max == 0 {
		max = 5
	}
	if len(factoids) > max {
		_, _ = ds.ChannelMessageSend(mc.ChannelID, fmt.Sprintf("Cannot send more than %d factoids at once", max))
		return
	}

	db, err := database.Get()
	if err != nil {
		_, _ = ds.ChannelMessageSend(mc.ChannelID, "Failed to connect to database")
		logger.Err().Printf("Failed to connect to database\n%s", err)
		return
	}

	var data []factoid
	err = db.Where("name IN (?)", factoids).Find(&data).Error

	if gorm.IsRecordNotFoundError(err) || (err == nil && len(data) == 0) {
		_, err = ds.ChannelMessageSend(mc.ChannelID, "No factoid with the given name was found")
		return
	} else if err != nil {
		logger.Err().Printf("Failed to pull data from database\n%s", err)
		return
	}

	if len(factoids) != len(data) {
		//we have a missing one...
		missing := make([]string, 0)
		for _, v := range factoids {
			good := false
			for _, k := range data {
				if v == k.Name {
					good = true
					break
				}
			}
			if !good {
				missing = append(missing, v)
			}
		}

		_, err = ds.ChannelMessageSend(mc.ChannelID, "No factoid with the given name(s) was found: "+strings.Join(missing, ", "))
		return
	}

	msg := ""
	for i, v := range factoids {
		for _, o := range data {
			if o.Name == v {
				msg += o.Content
				if i+1 != len(factoids) {
					msg += "\n\n"
				}
			}
		}
	}

	msg = strings.Replace(msg, "[b]", "**", -1)
	msg = strings.Replace(msg, "[/b]", "**", -1)
	msg = strings.Replace(msg, "[u]", "__", -1)
	msg = strings.Replace(msg, "[/u]", "__", -1)
	msg = strings.Replace(msg, "[i]", "*", -1)
	msg = strings.Replace(msg, "[/i]", "*", -1)
	msg = strings.Replace(msg, ";;", "\n", -1)

	if strings.Contains(msg, "https://") || strings.Contains(msg, "http://") {
		msgsplit := strings.Split(msg, " ")
		for k, v := range msgsplit {
			if strings.HasPrefix(v, "https://") || strings.HasPrefix(v, "http://") {
				msgsplit[k] = "<" + v + ">"
			}
		}
		msg = strings.Join(msgsplit, " ")
	}

	header := ""
	if len(mc.Mentions) > 0 {
		//if we have an @, we'll add it to the message
		for _, v := range mc.Mentions {
			header += "<@" + v.ID + "> "
		}
		header += "Please refer to the below information."
	}

	embed := &discordgo.MessageEmbed{
		Description: msg,
		Footer: &discordgo.MessageEmbedFooter{
			Text: "I am a bot, I will not respond to you\nIssued by " + mc.Author.Username + "#" + mc.Author.Discriminator,
		},
	}

	send := &discordgo.MessageSend{
		Content: header,
		Embed:   embed,
	}

        if viper.GetBool("factoid.delete") {
		_ = ds.ChannelMessageDelete(mc.ChannelID, mc.ID)
	}

	_, err = ds.ChannelMessageSendComplex(mc.ChannelID, send)
	if err != nil {
		logger.Err().Printf("Failed to send message\n%s", err)
	}
	
}

type factoid struct {
	Name    string `gorm:"name"`
	Content string `gorm:"content"`
}
