package config

import "time"

type Config struct {
	Postgres PostgresConfig
	Redis    RedisConfig
	Telegram TelegramConfig
	Parser   ParserConfig
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
}

type TelegramConfig struct {
	Token string `env:"TELEGRAM_TOKEN" env-required:"true"`
}
