package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	Telegram TelegramConfig
	Parser   ParserConfig
	JWT      JWTConfig
}
type JWTConfig struct {
	Secret string `env:"JWT_SECRET" env-required:"true"`
}

type PostgresConfig struct {
	Name     string `env:"PG_DB" env-required:"true"`
	Host     string `env:"PG_HOST" env-required:"true"`
	User     string `env:"PG_USER" env-required:"true"`
	Password string `env:"PG_PASSWORD" env-required:"true"`
	Port     string `env:"PG_PORT" env-default:"5432"`
}

type RedisConfig struct {
	Database int    `env:"REDIS_DB" env-required:"true"`
	Host     string `env:"REDIS_HOST" env-default:"127.0.0.1"`
	Port     string `env:"REDIS_PORT" env-default:"6379"`
	Password string `env:"REDIS_PASSWORD" env-required:"true"`
}

type ParserConfig struct {
	Interval time.Duration `env:"PARSER_INTERVAL" env-default:"30m"`
	Query    string        `env:"PARSER_QUERY" env-default:""`
}

type TelegramConfig struct {
	Token string `env:"TELEGRAM_TOKEN" env-required:"true"`
}

func MustLoad() *Config {
	cfg := Config{}
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}

func (c *PostgresConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		c.User, c.Password, c.Host, c.Port, c.Name)
}
