package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const ApiTelegramUrl = "https://api.telegram.org/bot%s/sendMessage"

type Response struct {
	Success bool `json:"ok"`
}

type QueryParam struct {
	key   string
	value string
}

type Request struct {
	ChatId   int64
	Message  string
	BotToken string
}

func SendMessage(request Request) {
	var apiUrl = fmt.Sprintf(ApiTelegramUrl, request.BotToken)
	var params = []QueryParam{
		{
			key:   "chat_id",
			value: strconv.FormatInt(request.ChatId, 10),
		},
		{
			key:   "text",
			value: request.Message,
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

func execute(httpUrl string, params []QueryParam) Response {
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

	var payload Response
	err = json.Unmarshal(body, &payload)

	return payload
}
