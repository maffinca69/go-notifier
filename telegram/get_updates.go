package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"release-notifier/infrastructure"
	"slices"
)

const (
	LongPullingTimeout = 60
	LongPullingOffset  = 0
)

func LongPulling() {
	bot := BotAPI()

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(LongPullingOffset)
	u.Timeout = LongPullingTimeout

	updates := bot.GetUpdatesChan(u)
	for update := range updates {
		var message = update.Message
		if message == nil || !message.IsCommand() {
			continue
		}

		var command = message.Command()
		if isRestrictedUser(message.From.ID) && isRestrictedCommand(command) {
			log.Println(fmt.Sprintf("Restricted user: %s", message.From.UserName))
			continue
		}

		var commandHandler = ResolveCommand(command)
		if commandHandler != nil {
			commandHandler.Handle(message, bot)
		}
	}
}

func isRestrictedCommand(command string) bool {
	allowedCommand := infrastructure.Config().AllowedCommands

	return !slices.Contains(allowedCommand, command)
}

func isRestrictedUser(userId int64) bool {
	allowedUsers := infrastructure.Config().AllowedUsers

	return !slices.Contains(allowedUsers, userId)
}
