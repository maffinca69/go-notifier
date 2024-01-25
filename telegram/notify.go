package telegram

import (
	"fmt"
	"os"
	"release-notifier/github"
	"release-notifier/infrastructure"
)

func Notify(release *github.Release, name string) {
	var text = infrastructure.GetConfig().TelegramConfig.NewVersionMessage

	var request = Request{
		ChatId:   infrastructure.GetConfig().TelegramConfig.ChatID,
		Message:  fmt.Sprintf(text, name, release.TagName, release.Url),
		BotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
	}

	SendMessage(request)
}
