package telegram

import (
	"sync"

	"gopkg.in/telebot.v4"
)

type Handler struct {
	adminID id

	adminModeMu sync.RWMutex
	inAdminMode map[int64]struct{}
}

func NewHandler(adminID int64) *Handler {
	return &Handler{
		adminID:     newID(adminID),
		inAdminMode: make(map[int64]struct{}),
	}
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
			handler: h.adminMode,
		},
	}
	err := cmds.set(b)
	if err != nil {
		return err
	}

	b.Handle(telebot.OnText, h.onText)
	return nil
}
