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

	httpclient "github.com/evzpav/telegram-go/http_client"
)

const telegramURL string = "https://api.telegram.org"
const defaultTimeout = 60

type telegramPayload struct {
	ChatID    string `json:"chat_id"`
	Text      string `json:"text"`
	ParseMode string `json:"parse_mode"`
}

//From ...
type From struct {
	ID        int64  `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
}

//Chat ...
type Chat struct {
	ID                          int64  `json:"id"`
	Title                       string `json:"title"`
	Type                        string `json:"type"`
	AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
}

//Result ...
type Result struct {
	From     From                     `json:"from"`
	Chat     Chat                     `json:"chat"`
	Date     int64                    `json:"date"`
	Text     string                   `json:"text"`
	Entities []map[string]interface{} `json:"entities"`
}

//Response is the struct of the status code 200 response
type Response struct {
	OK          bool   `json:"ok"`
	Result      Result `json:"result,omitempty"`
	ErrorCode   int    `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
}

//Client struct
type Client struct {
	httpClient httpclient.HTTPClient
	botToken   string
	groupID    string
	baseURL    string
	timeout    int
}

//New creates telegram Client
func New(botToken, telegramGroupID string) *Client {
	httpClient := &http.Client{
		Timeout: defaultTimeout * time.Second,
		Transport: &http.Transport{
			Dial:                (&net.Dialer{Timeout: defaultTimeout * time.Second}).Dial,
			TLSHandshakeTimeout: defaultTimeout * time.Second,
		},
	}
	return NewWithArguments(telegramURL, botToken, telegramGroupID, httpClient)
}

//NewWithArguments creates telegram Client with arguments
func NewWithArguments(url, botToken, groupID string, httpClient httpclient.HTTPClient) *Client {
	return &Client{
		httpClient: httpClient,
		baseURL:    url,
		botToken:   botToken,
		groupID:    groupID,
	}
}

//WithHTTPClient set new http Client if needed
func (c *Client) WithHTTPClient(newHTTPClient httpclient.HTTPClient) {
	c.httpClient = newHTTPClient
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
func (c *Client) SendMessage(text string) (Response, error) {
	payload := telegramPayload{
		ChatID:    c.groupID,
		Text:      text,
		ParseMode: "HTML",
	}
	path := fmt.Sprintf("/bot%s/sendMessage", c.botToken)

	respBs, err := c.baseRequest(http.MethodPost, path, payload)
	if err != nil {
		return Response{}, err
	}

	var response Response
	if err = json.Unmarshal(respBs, &response); err != nil {
		return Response{}, err
	}

	if !response.OK {
		return Response{}, fmt.Errorf("failed to send message. error code: [%d]; description: [%s]", response.ErrorCode, response.Description)
	}

	return response, nil
}
