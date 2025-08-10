package telegram

import (
	"bot/internal/utils/md"
	"fmt"
)

var (
	greetingsTmpl = langString{
		uaru: "Привіт, @%s!\n" +
			"Тут ти зможеш відслідковувати турніри по грі *Rocket League* та, можливо, щось ще...\n" +
			"Для додаткової інформації напиши команду /info",
		other: "Hello, @%s!\n" +
			"Here you can follow *Rocket League* tournaments and maybe something else...\n" +
			"For additional information, type the command /info",
	}
)

func greetings(lang string, username string) string {
	return fmt.Sprintf(greetingsTmpl.in(lang), md.Escape(username))
}
