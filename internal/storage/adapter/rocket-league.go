package adapter

import rocketleague "bot/internal/models/rocket-league"

func PlayersToDB(p rocketleague.Players) string {
	switch p {
	case rocketleague.P2x2:
		return "2x2"
	case rocketleague.P3x3:
		return "3x3"
	default:
		panic("unknown players mode")
	}
}

func ModeToDB(m rocketleague.Mode) string {
	switch m {
	case rocketleague.Soccer:
		return "soccer"
	case rocketleague.Pentathlon:
		return "pentathlon"
	default:
		panic("unknown tournament mode")
	}
}
