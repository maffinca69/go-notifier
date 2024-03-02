package commands

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strconv"
)

type GetMeInfo struct {
}

func (l GetMeInfo) Handle(message *tgbotapi.Message, api *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, getInfo(*message))
	msg.ReplyToMessageID = message.MessageID
	msg.ParseMode = "Markdown"

	_, err := api.Send(msg)
	if err != nil {
		_ = errors.New(err.Error())
	}
}

func getInfo(message tgbotapi.Message) string {
	var (
		text string
	)

	text += fmt.Sprintf("ID: `%s`\n\n", strconv.FormatInt(message.From.ID, 10))
	text += fmt.Sprintf("Username: @%s\n", message.From.UserName)
	text += fmt.Sprintf("Lang: %s\n", message.From.LanguageCode)
	text += fmt.Sprintf("Name: %s %s\n", message.From.FirstName, message.From.LastName)

	return text
}
