package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const ApiTelegramUrl = "https://api.telegram.org/bot%s/sendMessage"

type TelegramResponse struct {
	Success bool `json:"ok"`
}

type TelegramRequest struct {
	key   string
	value string
}

func SendMessage(chatId int64, message string, botToken string) {
	var apiUrl = fmt.Sprintf(ApiTelegramUrl, botToken)
	var params = []TelegramRequest{
		{
			key:   "chat_id",
			value: strconv.FormatInt(chatId, 10),
		},
		{
			key:   "text",
			value: message,
		},
		{
			key:   "parse_mode",
			value: "Markdown",
		},
	}
	var response = execute(apiUrl, params)
	if !response.Success {
		panic("Telegram Request failed")
	}
}

func execute(httpUrl string, params []TelegramRequest) TelegramResponse {
	client := &http.Client{}
	req, err := http.NewRequest("GET", httpUrl, nil)
	if err != nil {
		panic("Error creating HTTP request")
	}

	req.Header.Add("Content-Type", "application/json")
	q := req.URL.Query()
	for _, param := range params {
		q.Add(param.key, param.value)
	}

	req.URL.RawQuery = q.Encode()

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		panic("Error sending HTTP request")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic("Error reading HTTP response body")
	}

	var payload TelegramResponse
	err = json.Unmarshal(body, &payload)

	return payload
}
