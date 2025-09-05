package main

import (
	"bot/config"
	"bot/internal/controllers/telegram"
	rocketleagueapi "bot/internal/external/http/rocket-league-api"
	model "bot/internal/models"
	subService "bot/internal/services/subscriptions"
	userService "bot/internal/services/users"
	subStorage "bot/internal/storage/subscriptions"
	userStorage "bot/internal/storage/users"
	"net/http"

	"bot/pkg/sqlite"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/telebot.v4"
	"gopkg.in/telebot.v4/middleware"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal("can't load config: ", err)
	}

	db, err := sqlite.New(cfg.DB.File)
	if err != nil {
		log.Fatal("can't init db: ", err)
	}

	tgNotificationCh := make(chan model.TgNotification)

	// storages
	userStor := userStorage.New(db)
	subStor := subStorage.New(db)

	// external
	rlApi := rocketleagueapi.New(rocketleagueapi.Options{
		Key:    cfg.API.Key,
		Region: cfg.API.Region,
		Client: &http.Client{},
	})

	// services
	userServ := userService.New(userStor)
	subServ := subService.New(rlApi, subStor)

	// controllers
	tgHandler := telegram.New(userServ, subServ, cfg.Tg.AdminID)

	tgBot, err := telebot.NewBot(telebot.Settings{
		Token:     cfg.Tg.BotToken,
		Poller:    &telebot.LongPoller{},
		ParseMode: telebot.ModeMarkdown,
	})
	if err != nil {
		log.Fatal("can't init tg bot: ", err)
	}

	tgBot.Use(
		middleware.Recover(),
	)
	err = tgHandler.Use(tgBot)
	if err != nil {
		log.Fatal("can't use handler: ", err)
	}

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	go tgHandler.Notify(tgBot, tgNotificationCh)
	go subServ.RunNotifications(tgNotificationCh)
	go tgBot.Start()
	log.Println("tg bot is running")

	<-interrupt
	tgBot.Stop()
	log.Println("tg bot is shutted down")
}
