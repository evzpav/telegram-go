package main

import (
	"github.com/evzpav/telegram-go/telegram"
	"github.com/joho/godotenv"
	"log"
	"os"
)

var telegramBotToken string
var telegramGroupID string

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	telegramBotToken = os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramGroupID = os.Getenv("TELEGRAM_GROUP_ID")
}

func main() {

	message := "<b> Bold text </b>"
	message += "\n" // new line
	message += "<i>Italic text</i> \n"
	message += "<code> This is code text </code> \n"
	message += "@BotFather \n" //use existing Telegram username

	t := telegram.NewTelegramClient(telegramBotToken, telegramGroupID)
	t.SendMessage(message)
}
