package telegram

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"os"
	"strconv"
)

var apiInstance *tgbotapi.BotAPI

func BotAPI() *tgbotapi.BotAPI {
	if apiInstance == nil {
		apiInstance = setupBotAPI()
	}

	return apiInstance
}

func setupBotAPI() *tgbotapi.BotAPI {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_BOT_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug, _ = strconv.ParseBool(os.Getenv("TELEGRAM_BOT_DEBUG"))

	return bot
}
