package main

import (
	"log"

	migr "github.com/golang-migrate/migrate/v4"
	gomgP "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

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

	db, err := sqlx.Connect("postgres", cfg.GetPostgresSqlxDSN())
	if err != nil {
		log.Fatalf("failed to connect postgres: %v", err)
	}

	db.SetMaxOpenConns(cfg.PostgresMaxOpenConns)
	db.SetConnMaxLifetime(cfg.PostgresMaxConnLifetime)
	db.SetMaxIdleConns(cfg.PostgresMaxIdleConns)

	driver, err := gomgP.WithInstance(db.DB, &gomgP.Config{})
	if err != nil {
		log.Fatalf("failed to init db driver: %v", err)
	}
	migrateInstance, err := migr.NewWithDatabaseInstance("file:/"+cfg.MigrationFilesDir, cfg.PostgresDBName, driver)
	if err != nil {
		log.Fatalf("failed to init migrate instance: %v", err)
	}

	if err := migrateInstance.Migrate(cfg.MigrateDatabaseVersion); err != nil && err != migr.ErrNoChange {
		log.Fatalf("failed to migrate db: %v", err)
	}

	log.Printf("db was successfully migrated")
}
