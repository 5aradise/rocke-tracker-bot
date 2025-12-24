package telegram

import (
	model "bot/internal/models"
	"bot/pkg/lang"
	"log"

	"gopkg.in/telebot.v4"
)

func (*Handler) Notify(b *telebot.Bot, c <-chan model.TgNotification) {
	log.Println("start sending notifications via telegram")
	for ntf := range c {
		for _, id := range ntf.IDs {
			msg := tournamentStartsInMsg(lang.English, ntf.Tournament)
			go notifyUser(b, id, msg)
		}
	}
}

func notifyUser(b *telebot.Bot, id int64, msg string) {
	_, err := b.Send(newID(id), msg)
	if err != nil {
		log.Printf("can't send notification: id=%d, err=%s", id, err)
	}
}
