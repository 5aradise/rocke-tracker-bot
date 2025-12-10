package adapter

import rocketleague "bot/internal/models/rocket-league"

func PlayersToDB(p rocketleague.Players) string {
	switch p {
	case rocketleague.P2x2:
		return "2x2"
	case rocketleague.P3x3:
		return "3x3"
	}
	panic("unknown players mode")
}

func DBToPlayers(p string) rocketleague.Players {
	switch p {
	case "2x2":
		return rocketleague.P2x2
	case "3x3":
		return rocketleague.P3x3
	}
	panic("unknown players mode")
}

func ModeToDB(m rocketleague.Mode) string {
	switch m {
	case rocketleague.Soccer:
		return "soccer"
	case rocketleague.Pentathlon:
		return "pentathlon"
	}
	panic("unknown tournament mode")
}

func DBToMode(m string) rocketleague.Mode {
	switch m {
	case "soccer":
		return rocketleague.Soccer
	case "pentathlon":
		return rocketleague.Pentathlon
	}
	panic("unknown tournament mode")
}
