package telegram

import (
	"bot/internal/utils/lang"
	"bot/internal/utils/md"
	"fmt"
)

const ( // telegram language codes
	english   = "en"
	ukrainian = "ua"
	russian   = "ru"
)

var (
	greetingsMsgTmpl = lang.NewString(
		"Hello, @%s!\n"+
			"Here you can follow _Rocket League_ tournaments and maybe something else...\n"+
			"Check out the *menu* to see all the commands",
		"Привіт, @%s!\n"+
			"Тут ти зможеш відслідковувати турніри по грі _Rocket League_ та, можливо, щось ще...\n"+
			"Переглянь *меню*, щоб побачити всі команди",
	)

	youAreInAdminModeMsg = lang.NewString(
		"You are in *admin mode*, each of your subsequent _text messages_ "+
			"will be sent to the administration.\nTo exit, type /admin again",
		"Ви в *адміністраторському режимі*, кожне твоє наступне _текстове повідомлення_ "+
			"буде надіслане адміністрації.\nЩоб вийти напишіть знову /admin",
	)
	youAreNotInAdminModeMsg = lang.NewString(
		"You aren't in *admin mode*",
		"Ви вийшли з *адміністраторського моду*",
	)
)

func greetingsMsg(langCode string, username string) string {
	return fmt.Sprintf(greetingsMsgTmpl.In(lang.Code(langCode)), md.Escape(username))
}
