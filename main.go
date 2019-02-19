package main

import (
	"github.com/evzpav/telegram-go/telegram"
)

func main() {

	message := "<b> Bold text </b>"
	message += "\n" // new line
	message += "<i>Italic text</i> \n"
	message += "<code> This is code text </code> \n"
	message += "@BotFather \n" //use existing Telegram username

	telegram.SendMessage(message)
}
