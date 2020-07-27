package config

import (
	"log"

	"github.com/joho/godotenv"
)

type Option func()

func ApplyOptions(opts ...Option) {
	for _, opt := range opts {
		opt()
	}
}

func WithLocalEnvFile(envFile string) Option {
	return func() {
		if envFile != "" {
			err := godotenv.Load(envFile)

			if err != nil {
				log.Fatalf("failed to load env file [%s]: %v", envFile, err)
			}
		}
	}
}
