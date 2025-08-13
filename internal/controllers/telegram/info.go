package telegram

import "gopkg.in/telebot.v4"

func (h *Handler) start(c telebot.Context) error {
	user := c.Sender()
	return c.Send(greetingsMsg(user.LanguageCode, user.Username))
}
