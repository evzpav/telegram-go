package telegram_test

import (
	"net/http"
	"testing"

	httpclient "github.com/evzpav/telegram-go/http_client"
	"github.com/evzpav/telegram-go/telegram"
	"github.com/stretchr/testify/assert"
)

func TestSendMessage(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		httpClientMock := &httpclient.Mock{
			ResponseStatus: http.StatusOK,
			ResponseBody: `{
				"ok": true,
				"result": {
				  "message_id": 64,
				  "from": {
					"id": 12345678,
					"is_bot": true,
					"first_name": "my-bot",
					"username": "mybot"
				  },
				  "chat": {
					"id": -10101010,
					"title": "My bot test",
					"type": "group",
					"all_members_are_administrators": true
				  },
				  "date": 1590338718,
				  "text": " Bold text \nItalic text \n This is code text  \n@BotFather",
				  "entities": [
					{ "offset": 0, "length": 11, "type": "bold" },
					{ "offset": 12, "length": 11, "type": "italic" },
					{ "offset": 25, "length": 19, "type": "code" },
					{ "offset": 46, "length": 10, "type": "mention" }
				  ]
				}
			  }
			  `,
		}

		tg := telegram.NewWithArguments("https://api.telegram.org", "token123", "groupID", httpClientMock)

		message := "<b> Bold text </b>"
		message += "\n" // new line
		message += "<i>Italic text</i> \n"
		message += "<code> This is code text </code> \n"
		message += "@BotFather \n" //use existing Telegram username

		resp, err := tg.SendMessage(message)
		assert.Nil(t, err)
		assert.Equal(t, true, resp.OK)
		assert.Equal(t, "https://api.telegram.org/bottoken123/sendMessage", httpClientMock.RequestURL)
		assert.Equal(t, " Bold text \nItalic text \n This is code text  \n@BotFather", resp.Result.Text)

	})

	t.Run("Not authorized - invalid credentials", func(t *testing.T) {
		httpClientMock := &httpclient.Mock{
			ResponseStatus: http.StatusOK,
			ResponseBody:   `{"ok":false,"error_code":401,"description":"Unauthorized"}`,
		}

		tg := telegram.NewWithArguments("https://api.telegram.org", "token123", "groupID", httpClientMock)

		message := "<b> Bold text </b>"

		resp, err := tg.SendMessage(message)
		assert.EqualError(t, err, "failed to send message. error code: [401]; description: [Unauthorized]")
		assert.Equal(t, false, resp.OK)

	})

	t.Run("Unknown error", func(t *testing.T) {
		httpClientMock := &httpclient.Mock{
			ResponseStatus: http.StatusInternalServerError,
		}

		tg := telegram.NewWithArguments("https://api.telegram.org", "token123", "groupID", httpClientMock)

		message := "<b> Bold text </b>"

		resp, err := tg.SendMessage(message)
		assert.EqualError(t, err, "request status code [500] [https://api.telegram.org/bottoken123/sendMessage]: ")
		assert.Equal(t, false, resp.OK)

	})

}
