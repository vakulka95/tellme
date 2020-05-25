package sqlx

import (
	"fmt"

	"gitlab.com/tellmecomua/tellme.api/app/persistence/model"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"gitlab.com/tellmecomua/tellme.api/app/config"
)

type Repository struct {
	e          *config.Environment
	cli        *sqlx.DB
	statements map[string]*sqlx.Stmt
}

func New(e *config.Environment) *Repository {
	return &Repository{e: e}
}

func (r *Repository) Connect() error {
	var err error

	r.cli, err = sqlx.Connect("postgres", r.e.GetPostgresSqlxDSN())
	if err != nil {
		return fmt.Errorf("failed to open db conn: %v", err)
	}

	r.cli.SetConnMaxLifetime(r.e.PostgresMaxConnLifetime)
	r.cli.SetMaxIdleConns(r.e.PostgresMaxIdleConns)
	r.cli.SetMaxOpenConns(r.e.PostgresMaxOpenConns)

	r.statements, err = prepareStatements(r.cli, model.GetRawQueries())
	if err != nil {
		return fmt.Errorf("failed to prepare statements: %v", err)
	}

	return nil
}

func (r *Repository) Disconnect() error {
	return r.cli.Close()
}

func (r *Repository) Stat() map[string]string {
	s := r.cli.Stats()

	return map[string]string{
		"idle":                 fmt.Sprint(s.Idle),
		"in_use":               fmt.Sprint(s.InUse),
		"max_idle_closed":      fmt.Sprint(s.MaxIdleClosed),
		"max_lifetime_closed":  fmt.Sprint(s.MaxLifetimeClosed),
		"max_open_connections": fmt.Sprint(s.MaxOpenConnections),
		"open_connections":     fmt.Sprint(s.OpenConnections),
		"wait_count":           fmt.Sprint(s.WaitCount),
		"wait_duration":        s.WaitDuration.String(),
	}
}

func (r *Repository) Name() string {
	return "sqlx"
}
