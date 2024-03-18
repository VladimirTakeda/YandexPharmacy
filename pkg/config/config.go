package config

import "git.yandex-academy.ru/ooornament/code_architecture/pkg/db/postgresql"

type Config struct {
	Host string `config:"APP_HOST" yaml:"host"`
	Port string `config:"APP_PORT" yaml:"port"`

	Postgres postgresql.Config `config:"postgres"`
}
