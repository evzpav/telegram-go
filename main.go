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

	t := telegram.NewClient(os.Getenv("TELEGRAM_BOT_TOKEN"), os.Getenv("TELEGRAM_GROUP_ID"))
	sentMessage, err := t.SendMessage(message)
	if err != nil {
		log.Println(err)
	}
	log.Println(sentMessage)
}
