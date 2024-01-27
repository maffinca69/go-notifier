package commands

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Ping struct {
}

func (l Ping) Handle(message *tgbotapi.Message, api *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(message.Chat.ID, "Pong")
	msg.ReplyToMessageID = message.MessageID

	_, err := api.Send(msg)
	if err != nil {
		_ = errors.New(err.Error())
	}
}
