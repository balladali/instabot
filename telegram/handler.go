package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

type Handler interface {
	HandleUpdate(update tgbotapi.Update) bool
}

type StartHandler struct{}

func (s StartHandler) HandleUpdate(update tgbotapi.Update) bool {
	return update.Message != nil && s.handleMessage(*update.Message)
}

func (s StartHandler) handleMessage(message tgbotapi.Message) bool {
	switch message.Text {
	case "/start":
		msg := tgbotapi.NewMessage(message.Chat.ID, "Hi! I can get you postsAction, storiesAction, reels etc of Instagram user "+
			":) \nJust send me Instagram username")
		_, err := Bot.Send(msg)
		if err != nil {
			log.Errorf("Can't send message to the chat: %v", err)
			return false
		}
		return true
	}
	return false
}
