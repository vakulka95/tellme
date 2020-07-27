package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/caarlos0/env"
)

//
// Config type - struct that contains
//	required data for an app and it's dependencies
//
type Environment struct {
	DomainName              string `env:"DOMAIN_NAME" envDefault:"127.0.0.1:8080"`
	ExpertDocumentsStoreDir string `env:"EXPERT_DOCUMENTS_STORE_DIR" envDefault:""`
	StaticFilesDir          string `env:"STATIC_FILES_DIR" envDefault:""`
	MigrationFilesDir       string `env:"MIGRATION_FILES_DIR" envDefault:""`
	MigrateDatabaseVersion  uint   `env:"MIGRATE_DATABASE_VERSION" envDefault:"132215"`

	ServePort string `env:"SERVE_PORT" envDefault:"8080"`

	PostgresDriver            string        `env:"POSTGRES_DRIVER" envDefault:"pgx"`
	PostgresHost              string        `env:"POSTGRES_HOST,required"`
	PostgresPort              string        `env:"POSTGRES_PORT,required"`
	PostgresUsername          string        `env:"POSTGRES_USERNAME,required"`
	PostgresPassword          string        `env:"POSTGRES_PASSWORD,required"`
	PostgresDBName            string        `env:"POSTGRES_DBNAME,required"`
	PostgresSSLMode           string        `env:"POSTGRES_SSLMODE" envDefault:"disable"`
	PostgresConnTimeout       int           `env:"POSTGRES_CONN_TIMEOUT" envDefault:"5"`
	PostgresMaxOpenConns      int           `env:"POSTGRES_MAX_OPEN_CONNS" envDefault:"10"`
	PostgresMaxIdleConns      int           `env:"POSTGRES_MAX_IDLE_CONNS" envDefault:"3"`
	PostgresMaxConnLifetime   time.Duration `env:"POSTGRES_MAX_CONN_LIFETIME" envDefault:"10s"`
	PostgresHealthCheckPeriod time.Duration `env:"POSTGRES_HEALTH_CHECK_PERIOD" envDefault:"5s"`

	AccessTokenDuration time.Duration `env:"ACCESS_TOKEN_DURATION"`
	JWTTokenSignKey     string        `env:"JWT_TOKEN_SIGN_KEY"`

	TurboSMSUsername string `env:"TURBO_SMS_USERNAME,required"`
	TurboSMSPassword string `env:"TURBO_SMS_PASSWORD,required"`

	GoogleCaptchaSecret string `env:"GOOGLE_CAPTCHA_SECRET,required"`

	NotProcessedRequisitionJobPeriod time.Duration `env:"NOT_PROCESSED_REQUISITION_JOB_PERIOD" envDefault:"1h"`
}

//
// GetPostgresSqlxDSN method - generate postgres dsn string for sqlx driver
//
func (e *Environment) GetPostgresSqlxDSN() string {
	var dsn string

	for param, value := range map[string]string{
		"host":            e.PostgresHost,
		"port":            e.PostgresPort,
		"dbname":          e.PostgresDBName,
		"user":            e.PostgresUsername,
		"password":        e.PostgresPassword,
		"connect_timeout": strconv.Itoa(e.PostgresConnTimeout),
		"sslmode":         e.PostgresSSLMode,
	} {
		if value != "" {
			dsn = fmt.Sprintf("%s %s", dsn, pairParams(param, value))
		}
	}

	return dsn
}

//
// GetPostgresPgxDSN method - generate postgres dsn string for pgx driver
//
func (e *Environment) GetPostgresPgxDSN() string {
	var dsn string

	for param, value := range map[string]string{
		"host":                     e.PostgresHost,
		"port":                     e.PostgresPort,
		"dbname":                   e.PostgresDBName,
		"user":                     e.PostgresUsername,
		"password":                 e.PostgresPassword,
		"connect_timeout":          strconv.Itoa(e.PostgresConnTimeout),
		"sslmode":                  e.PostgresSSLMode,
		"pool_max_conns":           strconv.Itoa(e.PostgresMaxOpenConns),
		"pool_max_conn_lifetime":   e.PostgresMaxConnLifetime.String(),
		"pool_health_check_period": e.PostgresHealthCheckPeriod.String(),
	} {
		if value != "" {
			dsn = fmt.Sprintf("%s %s", dsn, pairParams(param, value))
		}
	}

	return dsn
}

func pairParams(param, value string) string {
	return fmt.Sprintf("%s='%s'", param, value)
}

//
// ConfigFromEnv func - reads env by struct's fields 'env' annotation
//
func NewConfigFromEnv() (*Environment, error) {
	e := &Environment{}
	if err := env.Parse(e); err != nil {
		return nil, err
	}
	return e, nil
}
