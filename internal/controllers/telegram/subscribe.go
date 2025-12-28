package telegram

import (
	"bot/config"
	model "bot/internal/models"
	rocketleague "bot/internal/models/rocket-league"
	"context"
	"errors"
	"fmt"
	"strings"

	"gopkg.in/telebot.v4"
)

const maxSubBtnsInRow = 3 // more than 3 makes unreadable keyboard on small screens

var (
	errInvalidUnsubData = errors.New("invalid unsubscribe data")
)

var (
	subPlayersSelector = &telebot.ReplyMarkup{}
	subPlayers2x2Btn   = subPlayersSelector.Text("2x2")
	subPlayers3x3Btn   = subPlayersSelector.Text("3x3")

	subModeSelector      = &telebot.ReplyMarkup{OneTimeKeyboard: true}
	subModeSoccerBtn     = subModeSelector.Text("Soccer")
	subModePentathlonBtn = subModeSelector.Text("Pentathlon")
)

func init() {
	subPlayersSelector.Reply(
		subPlayersSelector.Row(subPlayers2x2Btn, subPlayers3x3Btn),
	)

	subModeSelector.Reply(
		subModeSelector.Row(subModeSoccerBtn, subModePentathlonBtn),
	)
}

func (*Handler) subscribe(c telebot.Context) error {
	user := c.Sender()
	userLang := userLanguage(user)

	return c.Send(choosePlayersModeMsg.In(userLang), subPlayersSelector)
}

func (h *Handler) onPlayersBtn(players rocketleague.Players) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		user := c.Sender()
		userLang := userLanguage(user)

		h.selectedPlayersMu.Lock()
		h.selectedPlayers[user.ID] = players
		h.selectedPlayersMu.Unlock()

		return c.Send(chooseGameModeMsg.In(userLang), subModeSelector)
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

	h.recentlyUnsubed.Delete(subscription{user.ID, sub})

	return c.Send(youHaveSubscribedMsg.In(userLang))
}

func (h *Handler) subscriptions(c telebot.Context) error {
	user := c.Sender()
	userLang := userLanguage(user)

	subs, serr := h.subs.ListTelegramUserSubscriptions(context.TODO(), user.ID)
	if !serr.IsZero() {
		return c.Send(unexpectedErrorMsg(userLang, serr.Error()))
	}

	if len(subs) == 0 {
		return c.Send(noSubscriptionsMsg.In(userLang))
	}
	return c.Send(subscriptionsList(userLang, subs))
}

func (h *Handler) unsubscribe(c telebot.Context) error {
	user := c.Sender()
	userLang := userLanguage(user)

	subs, serr := h.subs.ListTelegramUserSubscriptions(context.TODO(), user.ID)
	if !serr.IsZero() {
		return c.Send(unexpectedErrorMsg(userLang, serr.Error()))
	}

	if len(subs) == 0 {
		return c.Send(noSubscriptionsMsg.In(userLang))
	}

	return c.Send(selectUnsubMsg.In(userLang), unsubsSquareKeyboard(subs))
}

type subscription struct {
	tgID int64
	sub  model.Subscription
}

func (h *Handler) onSelectedUnsubBtn(c telebot.Context) error {
	user := c.Sender()
	userLang := userLanguage(user)
	data, _ := strings.CutPrefix(c.Data(), "\f") // i don't know where adds this prefix

	sub, err := subFromUnsubData(data)
	if err != nil {
		return c.Send(unexpectedErrorMsg(userLang, err.Error()))
	}

	_, unsubed := h.recentlyUnsubed.LoadOrStore(subscription{user.ID, sub}, struct{}{})
	if unsubed {
		return c.Send(youAreAlreadyUnsubscribedMsg.In(userLang))
	}

	serr := h.subs.UnsubscribeByTelegram(context.TODO(), user.ID, sub)
	if !serr.IsZero() {
		switch serr.Code {
		case config.CodeSubNotExist:
			return c.Send(youAreAlreadyUnsubscribedMsg.In(userLang))
		default:
			return c.Send(unexpectedErrorMsg(userLang, serr.Error()))
		}
	}

	return c.Send(youHaveUnsubscribedMsg.In(userLang))
}

func unsubsSquareKeyboard(subs []model.Subscription) *telebot.ReplyMarkup {
	keyboard := &telebot.ReplyMarkup{}
	btns := make([]telebot.Btn, 0, len(subs))
	for _, sub := range subs {
		btns = append(btns, keyboard.Data(
			tournamentMsg(sub),
			unsubDataFromSub(sub),
		))
	}
	keyboard.Inline(keyboard.Split(maxSubBtnsInRow, btns)...)
	return keyboard
}

func unsubDataFromSub(sub model.Subscription) string {
	var players, mode string
	switch sub.Players {
	case rocketleague.P2x2:
		players = "P2x2"
	case rocketleague.P3x3:
		players = "P3x3"
	default:
		panic("unknown players in subscription")
	}

	switch sub.Mode {
	case rocketleague.Soccer:
		mode = "Soccer"
	case rocketleague.Pentathlon:
		mode = "Pentathlon"
	default:
		panic("unknown mode in subscription")
	}

	return fmt.Sprintf("unsub:%s %s", players, mode)
}

func subFromUnsubData(data string) (model.Subscription, error) {
	data, found := strings.CutPrefix(data, "unsub:")
	if !found {
		return model.Subscription{}, fmt.Errorf("%w: no unsub prefix", errInvalidUnsubData)
	}

	parts := strings.SplitN(data, " ", 2)
	if len(parts) != 2 {
		return model.Subscription{}, fmt.Errorf("%w: bad format", errInvalidUnsubData)
	}

	var sub model.Subscription
	switch parts[0] {
	case "P2x2":
		sub.Players = rocketleague.P2x2
	case "P3x3":
		sub.Players = rocketleague.P3x3
	default:
		return model.Subscription{}, fmt.Errorf("%w: unknown players", errInvalidUnsubData)
	}
	switch parts[1] {
	case "Soccer":
		sub.Mode = rocketleague.Soccer
	case "Pentathlon":
		sub.Mode = rocketleague.Pentathlon
	default:
		return model.Subscription{}, fmt.Errorf("%w: unknown mode", errInvalidUnsubData)
	}
	return sub, nil
}
