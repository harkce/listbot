package listbot

import (
	"log"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var Client *linebot.Client

func InitBot() error {
	channelSecret := os.Getenv("CHANNEL_SECRET")
	channelToken := os.Getenv("CHANNEL_TOKEN")

	var err error
	if Client, err = linebot.New(channelSecret, channelToken); err != nil {
		log.Println(err)
		log.Println(channelSecret)
		log.Println(channelToken)
		return err
	}
	return nil
}
