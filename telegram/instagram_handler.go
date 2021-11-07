package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"instabot/instagram"
	"strconv"
	"strings"
)

type InstagramHandler struct {
	manager instagram.InstaManager
}

var (
	storiesAction = "storiesAction"
	postsAction   = "postsAction"
	backAction    = "back"
)

func NewInstagramHandler() InstagramHandler {
	return InstagramHandler{
		manager: instagram.NewInstagramManager(),
	}
}

func (i InstagramHandler) HandleUpdate(update tgbotapi.Update) bool {
	if update.Message != nil {
		return i.handleMessage(*update.Message)
	} else if update.CallbackQuery != nil {
		return i.handleCallbackQuery(*update.CallbackQuery)
	}
	return false
}

func (i InstagramHandler) handleMessage(message tgbotapi.Message) bool {
	username := message.Text
	return i.processAccount(username, message.Chat.ID)
}

func (i InstagramHandler) processAccount(username string, chatId int64) bool {
	info := i.manager.GetUserInfo(username)
	if info == nil {
		log.Errorf("Can't get info for username %s", username)
		return false
	}
	msg := tgbotapi.NewMessage(chatId, "Here is some information about user:"+
		"\n\n*Username:* "+escapeText(info.Username)+
		"\n*Biography:*\n"+escapeText(info.Biography)+
		"\n*Type:* "+info.Type+
		"\n*Followers:* "+info.Followers)
	msg.ParseMode = tgbotapi.ModeMarkdown

	storiesCallback := storiesAction + ":" + username + ":0"
	//postsCallback := postsAction + ":" + username
	markup := tgbotapi.InlineKeyboardMarkup{
		InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
			{
				{
					Text:         "Stories",
					CallbackData: &storiesCallback,
				},
				//{
				//	Text:         "Posts",
				//	CallbackData: &postsCallback,
				//},
			},
		},
	}
	msg.ReplyMarkup = markup
	return SendMessage(msg)
}

func (i InstagramHandler) handleCallbackQuery(callbackQuery tgbotapi.CallbackQuery) bool {
	callbackData := strings.Split(callbackQuery.Data, ":")
	action := callbackData[0]
	username := callbackData[1]
	switch action {
	case storiesAction:
		storyIndex, _ := strconv.Atoi(callbackData[2])
		stories := i.manager.GetStories(username)
		if len(stories) == 0 || storyIndex == len(stories) {
			msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, "There are no new stories")
			return SendMessage(msg)
		}
		story := stories[storyIndex]
		urls, _ := story.GetMediaUrls()
		text := fmt.Sprintf("<a href='%s'>Story %d</a>", urls[0], storyIndex+1)
		msg := tgbotapi.NewMessage(callbackQuery.Message.Chat.ID, text)
		msg.ParseMode = tgbotapi.ModeHTML

		storiesCallback := storiesAction + ":" + username + ":" + strconv.FormatInt(int64(storyIndex+1), 10)
		backToInfoCallback := backAction + ":" + username
		if storyIndex < len(stories)-1 {
			markup := tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
					{
						{
							Text:         "Next story",
							CallbackData: &storiesCallback,
						},
						{
							Text:         "Back to info",
							CallbackData: &backToInfoCallback,
						},
					},
				},
			}
			msg.ReplyMarkup = markup
		}
		if storyIndex == len(stories)-1 {
			markup := tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{
					{
						{
							Text:         "Back to info",
							CallbackData: &backToInfoCallback,
						},
					},
				},
			}
			msg.ReplyMarkup = markup
		}
		SendMessage(msg)
	case backAction:
		i.processAccount(username, callbackQuery.Message.Chat.ID)
	}
	return true
}

func escapeText(text string) string {
	return "`" + text + "`"
}
