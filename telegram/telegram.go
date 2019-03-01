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
var defaultTimeout = 10

//Client struct
type Client struct {
	HTTPClient *http.Client
	BotToken   string
	GroupID    string
}

//NewClient creates telegram client
func NewClient(telegramBotToken, telegramGroupID string) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: setSecondsDuration(defaultTimeout),
		},
		BotToken: telegramBotToken,
		GroupID:  telegramGroupID,
	}
}

//ChangeHTTPClient set new http client if needed
func (t *Client) ChangeHTTPClient(newHTTPClient *http.Client) {
	t.HTTPClient = newHTTPClient
}

//ChangeTimeout set new timeout in seconds
func (t *Client) ChangeTimeout(newTimeout int) {
	t.HTTPClient.Timeout = setSecondsDuration(newTimeout)
}

//SendMessage uses sendMessage method from Telegram API
func (t *Client) SendMessage(text string) (string, error) {
	var payload tgPayload
	payload.ChatID = t.GroupID
	payload.Text = text
	payload.ParseMode = "HTML"
	bs, err := json.Marshal(&payload)
	url := fmt.Sprintf("%s/bot%s/sendMessage", telegramURL, t.BotToken)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(bs))
	req.Header.Set("Content-Type", "application/json")
	resp, err := t.HTTPClient.Do(req)
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

func setSecondsDuration(seconds int) time.Duration {
	return time.Second * time.Duration(seconds)
}
