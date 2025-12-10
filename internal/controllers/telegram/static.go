package telegram

import (
	model "bot/internal/models"
	"bot/pkg/lang"
	"bot/pkg/md"
	"fmt"
	"strings"
)

const ( // telegram language codes
	english   = "en"
	ukrainian = "ua"
	russian   = "ru"
)

var (
	unexpectedErrorMsgTmpl = lang.NewString(
		"An unexpected error has occurred:\n%s",
		"Сталась непередбачувана помилка:\n%s",
	)

	greetingsMsgTmpl = lang.NewString(
		"Hello, @%s!\n"+
			"Here you can follow _Rocket League_ tournaments and maybe something else...\n"+ //nolint
			"Check out the *menu* to see all the commands",
		"Привіт, @%s!\n"+
			"Тут ти зможеш відслідковувати турніри по грі _Rocket League_ та, можливо, щось ще...\n"+ //nolint
			"Переглянь *меню*, щоб побачити всі команди",
	)

	noSubscriptionsMsg = lang.NewString(
		"You have no subscriptions",
		"У вас немає підписок",
	)
	subscriptionsListHeaderMsg = lang.NewString(
		"Your subscriptions:",
		"Ваші підписки:",
	)
	tournamentMsgTmpl = lang.NewString(
		"Players: %s; Mode: %s",
		"Гравці: %s; Режим: %s",
	)
	choosePlayersModeMsg = lang.NewString(
		"Choose players mode:",
		"Виберіть режим гравців:",
	)
	chooseGameModeMsg = lang.NewString(
		"Choose game mode:",
		"Оберіть режим гри:",
	)
	youHaveSubscribedMsg = lang.NewString(
		"You have subscribed for tournament!",
		"Ви підписались на турнір!",
	)
	youAreAlreadySubscribedMsg = lang.NewString(
		"You are already subscribed for this tournament!",
		"Ви вже підписані на цей турнір!",
	)
	pressStartMsg = lang.NewString(
		"Press /start",
		"Натисніть /start",
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

	tournamentStartsInMsgTmpl = lang.NewString(
		"The tournament starts in *10 minutes*\nPlayers: %s\nMode: %s",
		"Через *10 хвилин* турнір\nГравці: %s\nРежим: %s",
	)
)

func greetingsMsg(langCode lang.Code, username string) string {
	return fmt.Sprintf(greetingsMsgTmpl.In(langCode), md.Escape(username))
}

func unexpectedErrorMsg(langCode lang.Code, errMsg string) string {
	return fmt.Sprintf(unexpectedErrorMsgTmpl.In(langCode), md.Escape(errMsg))
}

func tournamentStartsInMsg(langCode lang.Code, players, mode string) string {
	return fmt.Sprintf(tournamentStartsInMsgTmpl.In(langCode), players, mode)
}

func subscriptionsList(langCode lang.Code, subscriptions []model.Subscription) string {
	tournamentMsgTmpl := tournamentMsgTmpl.In(langCode)

	msg := strings.Builder{}
	msg.WriteString(subscriptionsListHeaderMsg.In(langCode))
	for _, sub := range subscriptions {
		msg.WriteByte('\n')
		msg.WriteString(fmt.Sprintf(tournamentMsgTmpl, sub.Players, sub.Mode))
	}
	return msg.String()
}
