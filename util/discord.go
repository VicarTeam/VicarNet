package util

import (
	"os"

	"github.com/bwmarrin/discordgo"
)

var bot *discordgo.Session

func InitBot() {
	dcbot, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		panic(err)
	}

	bot = dcbot
}

func SendCodeThroughDiscord(discordId string, code string, exp string) {
	ch, err := bot.UserChannelCreate(discordId)
	if err != nil {
		panic(err)
	}

	_, err = bot.ChannelMessageSend(ch.ID, "Your verification code is: ```"+code+"``` (expires in "+exp+")")
	if err != nil {
		panic(err)
	}
}
