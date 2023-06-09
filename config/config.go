package config

import "github.com/AbsaOSS/env-binder/env"

var Configuration Config

type Config struct {
	Port         string `env:"PORT"`
	MongoUrl     string `env:"MONGO_URL"`
	DatabaseName string `env:"DATABASE_NAME"`
	DbUserName   string `env:"DB_USERNAME"`
	DbPassword   string `env:"DB_PASSWORD"`
	SecretToken  string `env:"SECRET_TOKEN"`
	LogPath      string `env:"LOG_PATH"`
}

func Init() {
	env.Bind(&Configuration)
}
