package telegram

import (
	"errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"release-notifier/telegram/commands"
)

type Command interface {
	Handle(message *tgbotapi.Message, api *tgbotapi.BotAPI)
}

func ResolveCommand(command string) Command {
	switch command {
	case "list":
		return commands.List{}
	case "ping":
		return commands.Ping{}
	default:
		_ = errors.New("No implemented yet: " + command)
		return nil
	}
}
