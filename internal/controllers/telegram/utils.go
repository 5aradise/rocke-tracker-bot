package telegram

import (
	"bot/pkg/lang"

	"gopkg.in/telebot.v4"
)

func userLanguage(u *telebot.User) lang.Code {
	return lang.Code(u.LanguageCode)
}
