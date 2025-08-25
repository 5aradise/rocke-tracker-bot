package model

import rocketleague "bot/internal/models/rocket-league"

type Subscription struct {
	Players rocketleague.Players
	Mode    rocketleague.Mode
}
