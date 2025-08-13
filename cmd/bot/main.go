package main

import (
	"bot/config"
	"bot/internal/controllers/telegram"
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

	tgBot, err := telebot.NewBot(telebot.Settings{
		Token:     cfg.Tg.BotToken,
		Poller:    &telebot.LongPoller{},
		ParseMode: telebot.ModeMarkdown,
	})
	if err != nil {
		log.Fatal("can't init tg bot: ", err)
	}

	tgHandler := telegram.NewHandler(cfg.Tg.AdminID)
	tgBot.Use(
		middleware.Recover(),
	)
	err = tgHandler.Use(tgBot)
	if err != nil {
		log.Fatal("can't use handler: ", err)
	}

	go func(tgBot *telebot.Bot) {
		interrupt := make(chan os.Signal, 1)
		signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
		<-interrupt
		tgBot.Stop()
	}(tgBot)

	log.Println("tg bot is running")
	tgBot.Start()
	log.Println("tg bot is shutted down")
}
