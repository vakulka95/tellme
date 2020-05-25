package config

import (
	"flag"
	"log"

	"github.com/joho/godotenv"
)

type Option func()

func ApplyOptions(opts ...Option) {
	for _, opt := range opts {
		opt()
	}
}

func WithLocalEnvFile() Option {
	return loadEnvFile
}

var envFile = flag.String("env", "", "Env config file")

func loadEnvFile() {
	flag.Parse()

	if *envFile != "" {
		err := godotenv.Load(*envFile)

		if err != nil {
			log.Fatalf("failed to load env file [%s]: %v", *envFile, err)
		}
	}
}
