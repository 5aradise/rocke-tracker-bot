package telegram

import (
	"bot/internal/utils/md"
	"fmt"
)

var (
	greetingsMsgTmpl = langString{
		other: "Hello, @%s!\n" +
			"Here you can follow _Rocket League_ tournaments and maybe something else...\n" +
			"Check out the *menu* to see all the commands",
		uaru: "Привіт, @%s!\n" +
			"Тут ти зможеш відслідковувати турніри по грі _Rocket League_ та, можливо, щось ще...\n" +
			"Переглянь *меню*, щоб побачити всі команди",
	}

	youAreInAdminModeMsg = langString{
		other: "You are in *admin mode*, each of your subsequent _text messages_ " +
			"will be sent to the administration.\nTo exit, type /admin again",
		uaru: "Ви в *адміністраторському режимі*, кожне твоє наступне _текстове повідомлення_ " +
			"буде надіслане адміністрації.\nЩоб вийти напишіть знову /admin",
	}
	youAreNotInAdminModeMsg = langString{
		other: "You aren't in *admin mode*",
		uaru:  "Ви вийшли з *адміністраторського моду*",
	}
)

func greetingsMsg(lang string, username string) string {
	return fmt.Sprintf(greetingsMsgTmpl.in(lang), md.Escape(username))
}
