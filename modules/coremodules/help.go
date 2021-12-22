package coremodules

import (
	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

func RunHelpCommand(session *discordgo.Session, message *discordgo.MessageCreate, _ string, _ []string) {
	prefix := viper.GetString("prefix")
	if prefix == "" {
		prefix = "eevee!"
	}
	helpcontent := "`" + prefix + "help` - show this message\n`" + prefix + "hug <@member>` - hug someone\n`" + prefix + "kiss <@member>` - kiss someone\n`" + prefix + "pat <@member>` - pat someone\n`" + prefix + "slap <@member>` - slap someone\n`" + prefix + "ping` - Get bot ping`"
	helpembed := &discordgo.MessageEmbed{
		Description: helpcontent,
		Color:       0,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    "Requested by " + message.Author.Username + "#" + message.Author.Discriminator,
			IconURL: message.Author.AvatarURL(""),
		},
	}
	helpmessage := &discordgo.MessageSend{
		Embed:     helpembed,
		Reference: message.Reference(),
	}
	_, err := session.ChannelMessageSendComplex(message.ChannelID, helpmessage)
	if err != nil {
		println(err.Error())
	}
}
