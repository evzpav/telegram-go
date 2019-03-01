package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

type tgPayload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

var telegramURL = "https://api.telegram.org"
var defaultTimeout = time.Second * time.Duration(10)

type TelegramClient struct {
	HttpClient *http.Client
	BotToken   string
	GroupID    string
}

func NewTelegramClient(telegramBotToken, telegramGroupID string) *TelegramClient {
	return &TelegramClient{
		HttpClient: &http.Client{
			Timeout: defaultTimeout,
		},
		BotToken: telegramBotToken,
		GroupID:  telegramGroupID,
	}
}

func (t *TelegramClient) ChangeHttpClient(newHttpClient *http.Client) *TelegramClient {
	t.HttpClient = newHttpClient
	return t
}

func (t *TelegramClient) ChangeTimeout(newTimeout time.Duration) *TelegramClient {
	t.HttpClient.Timeout = newTimeout
	return t
}

//SendMessage uses sendMessage method from Telegram API
func (t *TelegramClient) SendMessage(text string) (string, error) {
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
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil

}
