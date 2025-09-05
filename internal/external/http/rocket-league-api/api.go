package rocketleagueapi

import (
	model "bot/internal/models"
	"encoding/json"
	"log"
	"net/http"
)

const tournamentsUrl = "https://rocket-league1.p.rapidapi.com/tournaments/"

type API struct {
	key    string
	url    string
	client *http.Client
}

type Options struct {
	Key    string
	Region string
	Client *http.Client
}

func New(opts Options) API {
	return API{
		key:    opts.Key,
		url:    tournamentsUrl + opts.Region,
		client: opts.Client,
	}
}

func (api API) Tournaments() []model.Tournament {
	req, err := http.NewRequest(http.MethodGet, api.url, nil)
	if err != nil {
		panic(err)
	}
	setHeaders(req.Header, api.key)

	log.Println("making request to api")
	res, err := api.client.Do(req)
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var resp Response
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		panic(err)
	}
	log.Println("received:", resp)

	return resp.ToModel()
}

func setHeaders(h http.Header, apiKey string) {
	h.Add("x-rapidapi-key", apiKey)
	h.Add("x-rapidapi-host", "rocket-league1.p.rapidapi.com")
	h.Add("User-Agent", "RockeTracker bot")
	h.Add("Accept-Encoding", "identity")
}
