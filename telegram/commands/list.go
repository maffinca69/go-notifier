package commands

import (
	"errors"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"release-notifier/cache"
	"release-notifier/infrastructure"
)

type List struct {
}

func (l List) Handle(message *tgbotapi.Message, api *tgbotapi.BotAPI) {
	var messageText = getFormattedRepositories()
	msg := tgbotapi.NewMessage(message.Chat.ID, messageText)
	msg.ParseMode = "Markdown"
	msg.ReplyToMessageID = message.MessageID

	_, err := api.Send(msg)
	if err != nil {
		_ = errors.New(err.Error())
	}
}

func getFormattedRepositories() string {
	var repositories = "🗒️*Repositories*\n\n"
	for _, repository := range infrastructure.Config().Repository {
		var name = repository.Name
		var url = repository.Url

		var currentVersion = cache.GetCurrentVersion(name)
		if len(currentVersion) == 0 {
			currentVersion = "Unknown"
		}

		repositories += fmt.Sprintf("🔥 *%s:*\n*Url:* %s\n*Version:* %s\n\n", name, url, currentVersion)
	}

	return repositories
}
