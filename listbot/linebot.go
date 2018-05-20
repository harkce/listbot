package listbot

import (
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var Client *linebot.Client

func InitBot() error {
	channelSecret := os.Getenv("CHANNEL_SECRET")
	channelToken := os.Getenv("CHANNEL_TOKEN") + "="

	var err error
	if Client, err = linebot.New(channelSecret, channelToken); err != nil {
		return err
	}
	return nil
}
