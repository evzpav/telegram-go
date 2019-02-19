package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var netClient = &http.Client{
	Timeout: time.Second * time.Duration(10),
}

type tgPayload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

var telegramURL = "https://api.telegram.org"
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

//SendMessage uses sendMessage method from Telegram API
func SendMessage(text string) {
	var payload tgPayload
	payload.ChatID = telegramGroupID
	payload.Text = text
	payload.ParseMode = "HTML"
	bs, err := json.Marshal(&payload)
	url := fmt.Sprintf("%s/bot%s/sendMessage", telegramURL, telegramBotToken)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	req.Header.Set("Content-Type", "application/json")
	resp, err := netClient.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Telegram message sent:", string(body))

}
