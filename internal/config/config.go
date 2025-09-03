package config

import (
	"os"

	"github.com/kelseyhightower/envconfig"
	"github.com/stackus/dotenv"
)

type AppConfig struct {
	GOLWIDTH         int
	GOLHEIGHT        int
	BATTLESHIPWIDTH  int
	BATTLESHIPHEIGHT int
}

func InitConfig() (cfg AppConfig, err error) {
	if err = dotenv.Load(dotenv.EnvironmentFiles(os.Getenv("ENVIRONMENT"))); err != nil {
		return
	}

	err = envconfig.Process("", &cfg)

	return
}
