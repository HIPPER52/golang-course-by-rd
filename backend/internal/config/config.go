package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type Config struct {
	Env  string `env:"APP_ENV" envDefault:"development"`
	Port string `env:"PORT" envDefault:"8081"`

	MongoURI    string `env:"MONGO_URI" envDefault:"mongodb://root:root@localhost:27017"`
	MongoDbName string `env:"MONGO_DB_NAME" envDefault:"course_project"`

	AuthSecret      string `env:"AUTH_SECRET" envDefault:"supersecret"`
	TokenTTLMinutes string `env:"TOKEN_TTL_MINUTES" envDefault:"300"`

	RabbitURL      string `env:"RABBITMQ_URL" envDefault:"amqp://guest:guest@localhost:5672/"`
	RabbitExchange string `env:"RABBITMQ_EXCHANGE" envDefault:"messages"`

	LogLevel string `env:"LOG_LEVEL" envDefault:"info"`
}

func NewConfigFromEnv() (*Config, error) {
	cfg := &Config{}
	err := env.Parse(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config from env: %w", err)
	}

	return cfg, nil
}
