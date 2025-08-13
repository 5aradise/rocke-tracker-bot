package config

import "github.com/kelseyhightower/envconfig"

type (
	Config struct {
		Tg Tg
	}

	Tg struct {
		BotToken string `envconfig:"TG_BOT_TOKEN" required:"true"`
		AdminID  int64  `envconfig:"TG_ADMIN_ID" required:"true"`
	}
)

func New() (Config, error) {
	var c Config

	err := envconfig.Process("", &c)
	if err != nil {
		return Config{}, err
	}

	return c, nil
}
