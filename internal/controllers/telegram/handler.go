package telegram

import (
	"bot/internal/utils/lang"
	"sync"

	"gopkg.in/telebot.v4"
)

type Handler struct {
	adminID id

	adminModeMu sync.RWMutex
	inAdminMode map[int64]struct{}
}

func New(adminID int64) *Handler {
	return &Handler{
		adminID:     newID(adminID),
		inAdminMode: make(map[int64]struct{}),
	}
}

func (h *Handler) Use(b *telebot.Bot) error {
	cmds := commands{
		{
			cmd: "start",
			description: lang.NewString(
				"Greetings",
				"Привітання",
			),
			handler: h.start,
		},
		{
			cmd: "subscribe",
			description: lang.NewString(
				"Subscribe to tournament notifications",
				"Підписка на повідомлення про турніри",
			),
			handler: h.subscribe,
		},
		{
			cmd: "unsubscribe",
			description: lang.NewString(
				"Unsubscribe from tournament notifications",
				"Відписка від повідомлень про турніри",
			),
			handler: h.unsubscribe,
		},
		{
			cmd: "admin",
			description: lang.NewString(
				"Contact the administration",
				"Зв'язатися з адміністрацією",
			),
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
