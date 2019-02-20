package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type tgPayload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

var telegramURL = "https://api.telegram.org"

type TelegramClient struct {
	HttpClient *http.Client
	BotToken   string
	GroupID    string
}

func NewTelegramClient(telegramBotToken, telegramGroupID string) *TelegramClient {
	netClient := &http.Client{
		Timeout: time.Second * time.Duration(10),
	}

	return &TelegramClient{
		HttpClient: netClient,
		BotToken:   telegramBotToken,
		GroupID:    telegramGroupID,
	}
}

func (t *TelegramClient) changeHttpClient(newHttpClient *http.Client) *TelegramClient {
	t.HttpClient = newHttpClient
	return t
}

//SendMessage uses sendMessage method from Telegram API
func (t *TelegramClient) SendMessage(text string) {
	var payload tgPayload
	payload.ChatID = t.GroupID
	payload.Text = text
	payload.ParseMode = "HTML"
	bs, err := json.Marshal(&payload)
	url := fmt.Sprintf("%s/bot%s/sendMessage", telegramURL, t.BotToken)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.HttpClient.Do(req)
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
