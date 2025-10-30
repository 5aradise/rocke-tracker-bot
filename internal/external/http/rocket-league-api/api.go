package rocketleagueapi

import (
	model "bot/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

const tournamentsURL = "https://rocket-league1.p.rapidapi.com/tournaments/"

var (
	ErrRequestLimitExceeded = errors.New("request limit exceeded")
)

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
		url:    tournamentsURL + opts.Region,
		client: opts.Client,
	}
}

func (api API) Tournaments() ([]model.Tournament, error) {
	req, err := http.NewRequest(http.MethodGet, api.url, nil)
	if err != nil {
		panic(err)
	}
	setHeaders(req.Header, api.key)

	log.Println("making request to api")
	res, err := api.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("making request: %w", err)
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			log.Printf("can't close response body: %s\n", err)
		}
	}()

	switch res.StatusCode {
	case http.StatusOK:
	case http.StatusTooManyRequests:
		return nil, ErrRequestLimitExceeded
	default:
		return nil, fmt.Errorf("bad response status: %s", res.Status)
	}

	var body Response
	err = json.NewDecoder(res.Body).Decode(&body)
	if err != nil {
		panic(err)
	}
	log.Println("received:", body)

	return body.ToModel(), nil
}

func setHeaders(h http.Header, apiKey string) {
	h.Add("X-Rapidapi-Key", apiKey)
	h.Add("X-Rapidapi-Host", "rocket-league1.p.rapidapi.com")
	h.Add("User-Agent", "RockeTracker bot")
	h.Add("Accept-Encoding", "identity")
}
