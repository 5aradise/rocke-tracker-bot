package config

import "github.com/kelseyhightower/envconfig"

type (
	Config struct {
		Tg Tg
		DB DB
	}

	Tg struct {
		BotToken string `envconfig:"TG_BOT_TOKEN" required:"true"`
		AdminID  int64  `envconfig:"TG_ADMIN_ID" required:"true"`
	}

	DB struct {
		File string `envconfig:"DB_FILE" required:"true"`
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
