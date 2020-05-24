package main

import (
	"log"
	"os"

	"github.com/evzpav/telegram-go/telegram"
)

func main() {
	message := "<b> Bold text </b>"
	message += "\n" // new line
	message += "<i>Italic text</i> \n"
	message += "<code> This is code text </code> \n"
	message += "@BotFather \n" //use existing Telegram username

	t := telegram.New(os.Getenv("TELEGRAM_TOKEN"), os.Getenv("TELEGRAM_GROUP_ID"))
	telegramResponse, err := t.SendMessage(message)
	if err != nil {
		log.Printf("failed to send telegram message: %v\n", err)
	}
	log.Printf("Response: %+v\n", telegramResponse)
}
