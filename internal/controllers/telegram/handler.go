package telegram

import (
	rocketleague "bot/internal/models/rocket-league"
	subservice "bot/internal/services/subscriptions"
	userservice "bot/internal/services/users"
	"bot/pkg/cache"
	"bot/pkg/lang"
	"sync"
	"time"

	"gopkg.in/telebot.v4"
)

const recentlyUnsubedClearDelay = 10 * time.Minute

type storage[K comparable, V any] interface {
	LoadOrStore(K, V) (actual V, loaded bool)
	Delete(K)

	Stop()
}

type Handler struct {
	users *userservice.Service
	subs  *subservice.Service

	selectedPlayersMu sync.Mutex
	selectedPlayers   map[int64]rocketleague.Players

	// prevent rapid unsubscribing becouse client can do it easily
	recentlyUnsubed storage[subscription, struct{}]

	adminID id

	adminModeMu sync.RWMutex
	inAdminMode map[int64]struct{}
}

func New(userServ *userservice.Service, subServ *subservice.Service, adminID int64,
) *Handler {
	return &Handler{
		users: userServ,
		subs:  subServ,

		selectedPlayers: make(map[int64]rocketleague.Players),

		recentlyUnsubed: cache.New[subscription, struct{}](recentlyUnsubedClearDelay),

		adminID: newID(adminID),

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
			cmd: "subscriptions",
			description: lang.NewString(
				"Current tournament subscriptions",
				"Поточні підписки на турніри",
			),
			handler: h.subscriptions,
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

	// static reply buttons (subscription)
	b.Handle(&subPlayers2x2Btn, h.onPlayersBtn(rocketleague.P2x2))
	b.Handle(&subPlayers3x3Btn, h.onPlayersBtn(rocketleague.P3x3))
	// V V V
	b.Handle(&subModeSoccerBtn, h.onModeBtn(rocketleague.Soccer))
	b.Handle(&subModePentathlonBtn, h.onModeBtn(rocketleague.Pentathlon))

	// dynamic inline buttons (for now only unsubscription)
	b.Handle(telebot.OnCallback, h.onSelectedUnsubBtn)

	// for admin mode
	b.Handle(telebot.OnText, h.onText)

	return nil
}

func (h *Handler) Shutdown() {
	h.recentlyUnsubed.Stop()
}
