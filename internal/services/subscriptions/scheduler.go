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

	tryCount = 3
	tryDelay = 10 * time.Second
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
		if ts := s.fetchWithRetries(); ts != nil {
			ntf.scheduleNotifications(ts)
		} else {
			log.Printf("can't fetch tournaments: next update in %s day\n", updateDelay)
		}
		<-tick
	}
}

func (s *Service) fetchWithRetries() []model.Tournament {
	ts, err := s.api.Tournaments()
	if err == nil {
		return ts
	}
	log.Printf("can't fetch tournaments: %s\n", err)
	for i := range tryCount - 1 {
		time.Sleep(tryDelay)

		ts, err := s.api.Tournaments()
		if err == nil {
			return ts
		}
		log.Printf("(try %d) can't fetch tournaments: %s\n", i+2, err)
	}
	return nil
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
