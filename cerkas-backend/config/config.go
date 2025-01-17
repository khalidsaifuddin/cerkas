package config

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	DBType string `envconfig:"DB_TYPE" default:"postgres"`

	HTTPPort    string `envconfig:"HTTP_PORT" default:"8080"`
	Environment string `envconfig:"ENVIRONMENT" default:"staging"`

	Host            string `envconfig:"PPDB_BACKEND_PGSQL_HOST" default:"103.41.206.233"`
	Port            string `envconfig:"PPDB_BACKEND_PGSQL_PORT" default:"5432"`
	Username        string `envconfig:"PPDB_BACKEND_PGSQL_USERNAME" default:"postgres"`
	Password        string `envconfig:"PPDB_BACKEND_PGSQL_PASSWORD" default:"r00t_db"`
	DBName          string `envconfig:"PPDB_BACKEND_PGSQL_DBNAME" default:"simdik"`
	LogMode         bool   `envconfig:"DB_LOG_MODE" default:"true"`
	MaxIdleConns    int    `envconfig:"DB_MAX_IDLE_CONNS" default:"5"`
	MaxOpenConns    int    `envconfig:"DB_MAX_OPEN_CONNS" default:"10"`
	ConnMaxLifetime int    `envconfig:"DB_CONN_MAX_LIFETIME" default:"10"`

	RedisHost     string `envconfig:"REDIS_HOST" default:"127.0.0.1"`
	RedisPort     string `envconfig:"REDIS_PORT" default:"6379"`
	RedisPassword string `envconfig:"REDIS_PASSWORD" default:""`
	RedisMaxIdle  int    `envconfig:"REDIS_MAX_IDLE" default:"10"`
	DefaultTTL    int64  `envconfig:"DEFAULT_TTL" default:"3600"`

	PSPInternalSecretKey string `envconfig:"PSP_INTERNAL_SECRET_KEY" default:"PSP_INTERNAL_SECRET_KEY"`
}

func Get() Config {
	cfg := Config{}
	envconfig.MustProcess("", &cfg)
	return cfg
}
