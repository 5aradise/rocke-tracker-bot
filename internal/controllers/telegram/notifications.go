package telegram

import (
	model "bot/internal/models"
	rocketleague "bot/internal/models/rocket-league"
	"bot/pkg/lang"
	"log"

	"gopkg.in/telebot.v4"
)

func (h *Handler) Notify(b *telebot.Bot, c <-chan model.TgNotification) {
	log.Println("start sending notifications via telegram")
	for ntf := range c {
		players, mode := subStr(ntf.Tournament)
		for _, id := range ntf.IDs {
			go notifyUser(b, id, lang.English, players, mode)
		}
	}
}

func notifyUser(b *telebot.Bot, id int64, userLang lang.Code, players, mode string) {
	_, err := b.Send(newID(id), tournamentStartsInMsg(userLang, players, mode))
	if err != nil {
		log.Printf("can't send notification: id=%d, err=%s", id, err)
	}
}

func subStr(sub model.Subscription) (players, mode string) {
	switch sub.Players {
	case rocketleague.P2x2:
		players = "2x2"
	case rocketleague.P3x3:
		players = "3x3"
	default:
		panic("unknown players mode")
	}
	switch sub.Mode {
	case rocketleague.Soccer:
		mode = "Soccer"
	case rocketleague.Pentathlon:
		mode = "Pentathlon"
	default:
		panic("unknown game mode")
	}
	return
}
