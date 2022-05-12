package helper

import (
	"github.com/caarlos0/env"
	"github.com/joho/godotenv"
	"github.com/pkg/errors"
	"log"
)

func PrepareEnvConfig(appConfig interface{}) error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load config, error: ", err.Error())
		return errors.WithStack(err)
	}

	err = env.Parse(appConfig)
	if err != nil {
		log.Fatal("Failed to parse config, error: ", err.Error())
		return errors.WithStack(err)
	}
	return nil
}
