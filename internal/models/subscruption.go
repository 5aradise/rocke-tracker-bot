package model

import (
	rocketleague "bot/internal/models/rocket-league"
	"time"
)

type Subscription struct {
	Players rocketleague.Players
	Mode    rocketleague.Mode
}

type Tournament struct {
	Type   Subscription
	Starts time.Time
}

type TgNotification struct {
	Tournament Subscription
	IDs        []int64
}
