package telegram

import (
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"instabot/config"
)

var Bot *tgbotapi.BotAPI
var handlers = []Handler{StartHandler{}, NewInstagramHandler()}

func StartBot() {
	bot, err := tgbotapi.NewBotAPI(config.Cfg.Bot.Token)
	if err != nil {
		log.Panic(err)
	}

	Bot = bot
	Bot.Debug = config.Cfg.Bot.Debug

	log.Printf("Authorized on account %s", Bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := Bot.GetUpdatesChan(updateConfig)

	for update := range updates {
		for _, handler := range handlers {
			if handler.HandleUpdate(update) {
				break
			}
		}
	}
}

func SendMessage(msg tgbotapi.Chattable) bool {
	_, err := Bot.Send(msg)
	if err != nil {
		log.Errorf("Can't send message to the chat: %v", err)
		return false
	}
	return true
}
