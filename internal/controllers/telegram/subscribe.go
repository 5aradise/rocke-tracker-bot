package telegram

import (
	"bot/config"
	model "bot/internal/models"
	rocketleague "bot/internal/models/rocket-league"
	"context"

	"gopkg.in/telebot.v4"
)

var (
	playersSelector = &telebot.ReplyMarkup{}
	players2x2Btn   = playersSelector.Text("2x2")
	players3x3Btn   = playersSelector.Text("3x3")

	modeSelector      = &telebot.ReplyMarkup{OneTimeKeyboard: true}
	modeSoccerBtn     = playersSelector.Text("Soccer")
	modePentathlonBtn = playersSelector.Text("Pentathlon")
)

func init() {
	playersSelector.Reply(
		playersSelector.Row(players2x2Btn, players3x3Btn),
	)

	modeSelector.Reply(
		modeSelector.Row(modeSoccerBtn, modePentathlonBtn),
	)
}

func (h *Handler) subscribe(c telebot.Context) error {
	user := c.Sender()
	userLang := userLanguage(user)

	return c.Send(choosePlayersModeMsg.In(userLang), playersSelector)
}

func (h *Handler) onPlayersBtn(players rocketleague.Players) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		user := c.Sender()
		userLang := userLanguage(user)

		h.selectedPlayersMu.Lock()
		h.selectedPlayers[user.ID] = players
		h.selectedPlayersMu.Unlock()

		return c.Send(chooseGameModeMsg.In(userLang), modeSelector)
	}
}

func (h *Handler) onModeBtn(mode rocketleague.Mode) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		user := c.Sender()

		h.selectedPlayersMu.Lock()
		players, ok := h.selectedPlayers[user.ID]
		if !ok {
			panic("selecting a mode without selecting players")
		}
		delete(h.selectedPlayers, user.ID)
		h.selectedPlayersMu.Unlock()

		return h.createSub(c, model.Subscription{
			Players: players,
			Mode:    mode,
		})
	}
}

func (h *Handler) createSub(c telebot.Context, sub model.Subscription) error {
	user := c.Sender()
	userLang := userLanguage(user)

	_, serr := h.subs.SubscribeByTelegram(context.TODO(), user.ID, sub)
	if !serr.IsZero() {
		switch serr.Code {
		case config.CodeUserHasSub:
			return c.Send(youAreAlreadySubscribedMsg.In(userLang))
		case config.CodeUserWithTgIDNotExist:
			return c.Send(pressStartMsg.In(userLang))
		default:
			return c.Send(unexpectedErrorMsg(userLang, serr.Error()))
		}
	}

	return c.Send(youHaveSubscribedMsg.In(userLang))
}

func (h *Handler) unsubscribe(c telebot.Context) error {
	return c.Send("Unsub!")
}
