package telegram

import (
	"gopkg.in/telebot.v4"
)

type Handler struct {
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) Use(b *telebot.Bot) error {
	cmds := commands{
		{
			cmd: "start",
			description: langString{
				other: "Greetings",
				uaru:  "Привітання",
			},
			handler: h.start,
		},
		{
			cmd: "subscribe",
			description: langString{
				other: "Subscribe to tournament notifications",
				uaru:  "Підписка на повідомлення про турніри",
			},
			handler: h.subscribe,
		},
		{
			cmd: "unsubscribe",
			description: langString{
				other: "Unsubscribe from tournament notifications",
				uaru:  "Відписка від повідомлень про турніри",
			},
			handler: h.unsubscribe,
		},
		{
			cmd: "admin",
			description: langString{
				other: "Contact the administration",
				uaru:  "Зв'язатися з адміністрацією",
			},
			handler: h.admin,
		},
	}
	return cmds.set(b)
}

func (h Handler) start(c telebot.Context) error {
	sender := c.Sender()
	return c.Send(greetings(sender.LanguageCode, sender.Username))
}

func (h Handler) subscribe(c telebot.Context) error {
	return c.Send("Sub!")
}

func (h Handler) unsubscribe(c telebot.Context) error {
	return c.Send("Unsub!")
}

func (h Handler) admin(c telebot.Context) error {
	return c.Send("Admin!")
}
