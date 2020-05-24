package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const telegramURL string = "https://api.telegram.org"

type telegramPayload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

//Client struct
type Client struct {
	httpClient *http.Client
	botToken   string
	groupID    string
	baseURL    string
}

//New creates telegram Client
func New(botToken, telegramGroupID string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 60,
			Transport: &http.Transport{
				Dial:                (&net.Dialer{Timeout: 60}).Dial,
				TLSHandshakeTimeout: 60,
			},
		},
		baseURL:  telegramURL,
		botToken: botToken,
		groupID:  telegramGroupID,
	}
}

//WithHTTPClient set new http Client if needed
func (c *Client) WithHTTPClient(newHTTPClient *http.Client) {
	c.httpClient = newHTTPClient
}

//WithTimeout set new timeout in seconds
func (c *Client) WithTimeout(newTimeout int) {
	c.httpClient.Timeout = time.Second * time.Duration(newTimeout)
}

//WithURL set new URL
func (c *Client) WithURL(url string) {
	c.baseURL = url
}

func (c *Client) baseRequest(method, path string, body interface{}) ([]byte, error) {
	var bodyReader io.Reader

	if body != nil {
		bodyBs, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("error to marshal body %v", err)
		}
		bodyReader = bytes.NewBuffer(bodyBs)
	}

	url := c.baseURL + path

	req, err := http.NewRequest(method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("error in create request to [%s]: %v", url, err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error to complete request to [%s]: %v", url, err)
	}

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error of [%s] reading body: %v", url, err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request status code [%d] [%s]: %v", resp.StatusCode, url, string(bs))
	}

	return bs, nil
}

//SendMessage uses sendMessage method from Telegram API
func (c *Client) SendMessage(text string) (string, error) {
	payload := telegramPayload{
		ChatID:    c.groupID,
		Text:      text,
		ParseMode: "HTML",
	}
	path := fmt.Sprintf("/bot%s/sendMessage", c.botToken)

	respBs, err := c.baseRequest(http.MethodPost, path, payload)
	if err != nil {
		return "", err
	}

	return string(respBs), nil
}
