package telegram

import (
	rocketleague "bot/internal/models/rocket-league"
	subservice "bot/internal/services/subscriptions"
	userService "bot/internal/services/users"
	"bot/pkg/lang"
	"sync"

	"gopkg.in/telebot.v4"
)

type Handler struct {
	users *userService.Service
	subs  *subservice.Service

	selectedPlayersMu sync.Mutex
	selectedPlayers   map[int64]rocketleague.Players

	adminID id

	adminModeMu sync.RWMutex
	inAdminMode map[int64]struct{}
}

func New(userServ *userService.Service, subServ *subservice.Service, adminID int64,
) *Handler {
	return &Handler{
		users: userServ,
		subs:  subServ,

		selectedPlayers: make(map[int64]rocketleague.Players),

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

	// buttons (subscription)
	b.Handle(&players2x2Btn, h.onPlayersBtn(rocketleague.P2x2))
	b.Handle(&players3x3Btn, h.onPlayersBtn(rocketleague.P3x3))
	b.Handle(&modeSoccerBtn, h.onModeBtn(rocketleague.Soccer))
	b.Handle(&modePentathlonBtn, h.onModeBtn(rocketleague.Pentathlon))

	// for admin mode
	b.Handle(telebot.OnText, h.onText)

	return nil
}
