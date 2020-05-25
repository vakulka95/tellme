package main

import (
	"log"

	"gitlab.com/tellmecomua/tellme.api/app"
	"gitlab.com/tellmecomua/tellme.api/app/config"
)

func main() {
	config.ApplyOptions(
		config.WithLocalEnvFile(),
	)

	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to fetch config from env: %v", err)
	}

	apiserver := app.New(cfg)
	if err := apiserver.Run(); err != nil {
		log.Fatalf("failed to run apiserver: %v", err)
	}
}
