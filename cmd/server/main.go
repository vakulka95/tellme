package main

import (
	"flag"
	"log"

	migr "github.com/golang-migrate/migrate/v4"
	gomgP "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"

	"gitlab.com/tellmecomua/tellme.api/app"
	"gitlab.com/tellmecomua/tellme.api/app/config"
)

var (
	migrate bool
	envFile string
)

func init() {
	flag.BoolVar(&migrate, "migrate", false, "Migrate database")
	flag.StringVar(&envFile, "env", "", "Env file")
	flag.Parse()
}

func main() {
	config.ApplyOptions(
		config.WithLocalEnvFile(envFile),
	)

	cfg, err := config.NewConfigFromEnv()
	if err != nil {
		log.Fatalf("failed to fetch config from env: %v", err)
	}

	if migrate {
		if err = migrateDatabase(cfg); err != nil {
			log.Fatalf("failed to migrate database: %v", err)
		}
	}

	apiserver := app.New(cfg)
	if err := apiserver.Run(); err != nil {
		log.Fatalf("failed to run apiserver: %v", err)
	}
}

func migrateDatabase(cfg *config.Environment) error {
	db, err := sqlx.Connect("postgres", cfg.GetPostgresSqlxDSN())
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}

	db.SetMaxOpenConns(cfg.PostgresMaxOpenConns)
	db.SetConnMaxLifetime(cfg.PostgresMaxConnLifetime)
	db.SetMaxIdleConns(cfg.PostgresMaxIdleConns)

	driver, err := gomgP.WithInstance(db.DB, &gomgP.Config{})
	if err != nil {
		return err
	}

	migrateInstance, err := migr.NewWithDatabaseInstance("file:"+cfg.MigrationFilesDir, cfg.PostgresDBName, driver)
	if err != nil {
		return err
	}

	if err := migrateInstance.Migrate(cfg.MigrateDatabaseVersion); err != nil && err != migr.ErrNoChange {
		return err
	}

	return nil
}
