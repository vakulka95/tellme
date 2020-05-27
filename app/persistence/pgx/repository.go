package pgx

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"

	"gitlab.com/tellmecomua/tellme.api/app/config"
)

type Repository struct {
	e   *config.Environment
	cli *pgxpool.Pool
}

func New(e *config.Environment) *Repository {
	return &Repository{e: e}
}

func (r *Repository) Connect() error {
	cfg, err := pgxpool.ParseConfig(r.e.GetPostgresPgxDSN())
	if err != nil {
		return err
	}

	cfg.MaxConns = int32(r.e.PostgresMaxOpenConns)
	cfg.MaxConnLifetime = r.e.PostgresMaxConnLifetime
	cfg.HealthCheckPeriod = r.e.PostgresHealthCheckPeriod

	r.cli, err = pgxpool.ConnectConfig(context.TODO(), cfg)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) Disconnect() error {
	r.cli.Close()
	return nil
}

func (r *Repository) Stat() map[string]string {
	s := r.cli.Stat()

	return map[string]string{
		"max_conns":              fmt.Sprint(s.MaxConns()),
		"acquire_count":          fmt.Sprint(s.AcquireCount()),
		"acquired_conns":         fmt.Sprint(s.AcquiredConns()),
		"acquire_duration":       s.AcquireDuration().String(),
		"canceled_acquire_count": fmt.Sprint(s.CanceledAcquireCount()),
		"constructing_conns":     fmt.Sprint(s.ConstructingConns()),
		"empty_acquire_count":    fmt.Sprint(s.EmptyAcquireCount()),
		"idle_conns":             fmt.Sprint(s.IdleConns()),
		"rotal_conns":            fmt.Sprint(s.TotalConns()),
	}
}

func (r *Repository) Name() string {
	return "pgx"
}
