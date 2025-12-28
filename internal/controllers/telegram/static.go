package telegram

import (
	model "bot/internal/models"
	rocketleague "bot/internal/models/rocket-league"
	"bot/pkg/lang"
	"bot/pkg/md"
	"fmt"
	"strings"
)

var (
	unexpectedErrorMsgTmpl = lang.NewString(
		"An unexpected error has occurred:\n"+
			"%s",
		"Сталась непередбачувана помилка:\n"+
			"%s",
	) // (error message)

	greetingsMsgTmpl = lang.NewString(
		"Hello, @%s!\n"+
			"Here you can follow _Rocket League_ tournaments and maybe something else...\n"+
			"Check out the *menu* to see all the commands",
		"Привіт, @%s!\n"+
			"Тут ти зможеш відслідковувати турніри по грі _Rocket League_ та, можливо, щось ще...\n"+
			"Переглянь *меню*, щоб побачити всі команди",
	) // (username)

	tournamentMsgTmpl = "%s %s" // (mode, players)

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

	subscriptionsListHeaderMsg = lang.NewString(
		"Your subscriptions:",
		"Ваші підписки:",
	)
	noSubscriptionsMsg = lang.NewString(
		"You have no subscriptions",
		"У вас немає підписок",
	)

	selectUnsubMsg = lang.NewString(
		"Select the tournament you want to unsubscribe from:",
		"Оберіть, від якого турніру ви хочете відписатись:",
	)
	youHaveUnsubscribedMsg = lang.NewString(
		"You have unsubscribed from tournament!",
		"Ви відписались від турніру!",
	)
	youAreAlreadyUnsubscribedMsg = lang.NewString(
		"You are already unsubscribed from this tournament!",
		"Ви вже відписані від цього турніру!",
	)

	youAreInAdminModeMsg = lang.NewString(
		"You are in *admin mode*, each of your subsequent _text messages_ will be sent to the administration.\n"+
			"To exit, type /admin again",
		"Ви в *адміністраторському режимі*, кожне твоє наступне _текстове повідомлення_ буде надіслане адміністрації.\n"+
			"Щоб вийти напишіть знову /admin",
	)
	youAreNotInAdminModeMsg = lang.NewString(
		"You aren't in *admin mode*",
		"Ви вийшли з *адміністраторського моду*",
	)

	tournamentStartsInMsgTmpl = lang.NewString(
		"The tournament starts in *10 minutes*\n"+
			"Players: %s\n"+
			"Mode: %s",
		"Через *10 хвилин* турнір\n"+
			"Гравці: %s\n"+
			"Режим: %s",
	) // (players, mode)
)

func greetingsMsg(lang lang.Language, username string) string {
	return fmt.Sprintf(greetingsMsgTmpl.In(lang), md.Escape(username))
}

func unexpectedErrorMsg(lang lang.Language, errMsg string) string {
	return fmt.Sprintf(unexpectedErrorMsgTmpl.In(lang), md.Escape(errMsg))
}

func tournamentMsg(sub model.Subscription) string {
	players, mode := subscriptionStr(sub)
	return fmt.Sprintf(tournamentMsgTmpl, players, mode)
}

func tournamentStartsInMsg(lang lang.Language, sub model.Subscription) string {
	players, mode := subscriptionStr(sub)
	return fmt.Sprintf(tournamentStartsInMsgTmpl.In(lang), players, mode)
}

func subscriptionStr(sub model.Subscription) (players, mode string) {
	switch sub.Players {
	case rocketleague.P2x2:
		players = "2x2"
	case rocketleague.P3x3:
		players = "3x3"
	}
	switch sub.Mode {
	case rocketleague.Soccer:
		mode = "Soccer"
	case rocketleague.Pentathlon:
		mode = "Pentathlon"
	}
	return players, mode
}

func subscriptionsList(lang lang.Language, subscriptions []model.Subscription) string {
	msg := strings.Builder{}
	msg.WriteString(subscriptionsListHeaderMsg.In(lang))
	for _, sub := range subscriptions {
		msg.WriteString("\n- " + tournamentMsg(sub))
	}
	return msg.String()
}
