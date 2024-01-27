package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
)

const LongPullingTimeout = 60
const LongPullingOffset = 0

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
		var commandHandler = ResolveCommand(command)
		if commandHandler != nil {
			commandHandler.Handle(message, bot)
		}
	}
}
