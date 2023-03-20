package gotelegrambot

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func NewBot(token string) (*Bot, error) {
	var getMeResponse GetMeResponse
	b := &Bot{
		Token:    token,
		Handlers: make(map[string]func(*Bot, *Update, []string)),
	}

	r, err := b.GetMe()

	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&getMeResponse)

	if err != nil {
		return nil, err
	}

	if getMeResponse.User.Username == "" {
		err := errors.New("error: getMe response has no username")
		return nil, err
	}

	b.Username = getMeResponse.User.Username

	return b, err
}

func (bot *Bot) callAPI(method string, body io.Reader) (resp *http.Response, err error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", bot.Token, method)
	return http.Post(url, "application/json", body)

}

func (bot *Bot) GetMe() (resp *http.Response, err error) {
	return http.Get(fmt.Sprintf("https://api.telegram.org/bot%s/%s", bot.Token, "getMe"))
}

func (bot *Bot) SendMessage(d SendMessagePayload) (resp *http.Response, err error) {

	b, err := json.Marshal(d)

	if err != nil {
		return nil, err
	}

	r, err := bot.callAPI("sendMessage", bytes.NewBuffer(b))
	defer r.Body.Close()

	return r, err
}

func (bot *Bot) AnswerCallbackQuery(d AnswerCallbackQueryPayload) (resp *http.Response, err error) {
	b, err := json.Marshal(d)

	if err != nil {
		return nil, err
	}

	r, err := bot.callAPI("answerCallbackQuery", bytes.NewBuffer(b))
	defer r.Body.Close()

	return r, err
}

func (bot *Bot) EditMessageText(d EditMessageTextPayload) (resp *http.Response, err error) {
	b, err := json.Marshal(d)

	if err != nil {
		return nil, err
	}

	r, err := bot.callAPI("editMessageText", bytes.NewBuffer(b))
	defer r.Body.Close()

	return r, err
}

func (bot *Bot) SetUpdateHandler(s string, ud UpdateHandler) {
	bot.Handlers[s] = ud
	bot.Handlers[fmt.Sprintf("%s@%s", s, bot.Username)] = ud
}

func (bot *Bot) HandleUpdate(u *Update) {
	s := u.Message.Text

	trimmed := strings.Trim(s, " ")
	tokens := strings.Split(trimmed, " ")
	funcName := tokens[0]
	args := tokens[1:]

	if f, ok := bot.Handlers[funcName]; ok {
		f(bot, u, args)
	}

}
