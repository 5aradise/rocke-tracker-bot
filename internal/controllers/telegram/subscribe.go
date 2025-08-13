package telegram

import "gopkg.in/telebot.v4"

func (h *Handler) subscribe(c telebot.Context) error {
	return c.Send("Sub!")
}

func (h *Handler) unsubscribe(c telebot.Context) error {
	return c.Send("Unsub!")
}
