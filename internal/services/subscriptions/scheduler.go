package subService

import (
	model "bot/internal/models"
	"context"
	"log"
	"time"
)

const (
	updateDelay     = 24 * time.Hour
	tournamentDelay = 10 * time.Minute
	timeForDB       = 2 * time.Minute
)

type notifier struct {
	s   *Service
	tgC chan model.TgNotification
}

func (s *Service) RunNotifications(tgC chan model.TgNotification) {
	ntf := notifier{
		s:   s,
		tgC: tgC,
	}
	tick := time.Tick(updateDelay)
	log.Println("start scheduling notifications")
	for {
		ts := s.api.Tournaments()
		ntf.scheduleNotifications(ts)
		<-tick
	}
}

func (ntf notifier) scheduleNotifications(ts []model.Tournament) {
	for _, t := range ts {
		preparationStart := t.Starts.Add(-(tournamentDelay + timeForDB))
		log.Printf("scheduled preparing at %s", preparationStart)
		atFunc(preparationStart, func() {
			ntf.prepareNotification(t.Starts, t.Type)
		})
	}
}

func (ntf notifier) prepareNotification(tourStart time.Time, tourType model.Subscription) {
	log.Printf("preparing notification for %v, that starts at %s\n", tourType, tourStart)
	tgIDs, err := ntf.s.subs.ListTelegramIDsBySubscription(context.Background(), tourType)
	if err != nil {
		log.Printf("can't load telegram ids from db: %s\n", err)
		return
	}
	log.Printf("fetched ids for %v: %v", tourType, tgIDs)

	if len(tgIDs) > 0 {
		notificationStart := tourStart.Add(-tournamentDelay)
		atFunc(notificationStart, func() {
			ntf.tgC <- model.TgNotification{
				Tournament: tourType,
				IDs:        tgIDs,
			}
		})
	}
}

func atFunc(when time.Time, f func()) {
	time.AfterFunc(time.Until(when), f)
}
