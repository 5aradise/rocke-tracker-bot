package telegram

import (
	"bot/config"
	"context"

	"gopkg.in/telebot.v4"
)

func (h *Handler) start(c telebot.Context) error {
	user := c.Sender()
	userLang := userLanguage(user)

	_, serr := h.users.CreateUser(context.TODO(), user.ID)
	if !serr.IsZero() {
		switch serr.Code {
		default:
			return c.Send(unexpectedErrorMsg(userLang, serr.Error()))
		case config.CodeUserWithTgIDExist:
		}
	}

	return c.Send(greetingsMsg(userLang, user.Username))
}
