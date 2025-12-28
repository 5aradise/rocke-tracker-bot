package telegram

import (
	"bot/pkg/lang"

	"gopkg.in/telebot.v4"
)

const ( // telegram language codes
	englishLandCode   = "en"
	ukrainianLandCode = "ua"
	russianLandCode   = "ru"
)

func userLanguage(u *telebot.User) lang.Language {
	switch u.LanguageCode {
	case englishLandCode:
		return lang.English
	case ukrainianLandCode:
		return lang.Ukrainian
	case russianLandCode:
		return lang.Russian
	}
	return lang.Other
}
