package telegram

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"release-notifier/github"
	"release-notifier/infrastructure"
)

func Notify(release *github.Release, name string) {
	var text = infrastructure.Config().TelegramConfig.NewVersionMessage

	var messageText = fmt.Sprintf(text, name, release.TagName, release.Url)
	msg := tgbotapi.NewMessage(infrastructure.Config().TelegramConfig.ChatID, messageText)
	msg.ParseMode = "Markdown"

	_, err := BotAPI().Send(msg)
	if err != nil {
		_ = errors.New(err.Error())
	}
}
