package rocketleagueapi

import (
	model "bot/internal/models"
	rocketleague "bot/internal/models/rocket-league"
	"fmt"
	"log"
	"time"
)

type Response struct {
	Tournaments []struct {
		Players int       `json:"players"`
		Starts  time.Time `json:"starts"`
		Mode    string    `json:"mode"`
	} `json:"tournaments"`
}

func (resp Response) ToModel() []model.Tournament {
	ms := make([]model.Tournament, 0, len(resp.Tournaments))
	for _, t := range resp.Tournaments {
		p, err := players(t.Players)
		if err != nil {
			log.Printf("can't parse players: %s\n", err)
			continue
		}
		m, err := mode(t.Mode)
		if err != nil {
			log.Printf("can't parse mode: %s\n", err)
			continue
		}
		ms = append(ms, model.Tournament{
			Type: model.Subscription{
				Players: p,
				Mode:    m,
			},
			Starts: t.Starts,
		})
	}
	return ms
}

func players(n int) (rocketleague.Players, error) {
	switch n {
	case 2: //nolint
		return rocketleague.P2x2, nil
	case 3: //nolint
		return rocketleague.P3x3, nil
	default:
		return 0, fmt.Errorf("unknown number of players mode: %d", n)
	}
}

func mode(s string) (rocketleague.Mode, error) {
	switch s {
	case "Soccer":
		return rocketleague.Soccer, nil
	case "": // bad api
		return rocketleague.Pentathlon, nil
	default:
		return 0, fmt.Errorf("unknown game mode: %s", s)
	}
}
