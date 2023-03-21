package gotelegrambot

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func callAPI(t string, method string, body io.Reader) (resp *http.Response, err error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", t, method)
	return http.Post(url, "application/json", body)

}

func GetMe(t string) (resp *http.Response, err error) {
	return http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/%s", t, "getMe"))
}

func SendMessage(t string, d SendMessagePayload) (resp *http.Response, err error) {

	b, err := json.Marshal(d)

	if err != nil {
		return nil, err
	}

	r, err := callAPI(t, "sendMessage", bytes.NewBuffer(b))
	defer r.Body.Close()

	return r, err
}

func AnswerCallbackQuery(t string, d AnswerCallbackQueryPayload) (resp *http.Response, err error) {
	b, err := json.Marshal(d)

	if err != nil {
		return nil, err
	}

	r, err := callAPI(t, "answerCallbackQuery", bytes.NewBuffer(b))
	defer r.Body.Close()

	return r, err
}

func EditMessageText(t string, d EditMessageTextPayload) (resp *http.Response, err error) {
	b, err := json.Marshal(d)

	if err != nil {
		return nil, err
	}

	r, err := callAPI(t, "editMessageText", bytes.NewBuffer(b))
	defer r.Body.Close()

	return r, err
}
