package config

import (
	"log"

	"github.com/spf13/viper"
)

// Config holds config struct
type Config struct {
	Port     string `env:"SIG_PORT"`
	DBType   string `env:"SIG_DB_TYPE"` // memory | postgres | mongo
	Postgres struct {
		DSN string `env:"SIG_DB_POSTGRES_DSN"`
	}
	Mongo struct {
		URI      string `env:"SIG_DB_MONGO_URI"`
		Database string `env:"SIG_DB_MONGO_DATABASE"`
	}
	InMemory struct {
		DBFilePath string `env:"SIG_DB_MEMORY_FILE"`
	}
}

// Load config values from env and config file
func Load() Config {
	var cfg Config

	cfg.Port = viper.GetString("port")
	cfg.DBType = viper.GetString("db.type")

	// database connectors configs
	cfg.Postgres.DSN = viper.GetString("db.postgres.dsn")

	cfg.Mongo.URI = viper.GetString("db.mongo.uri")
	cfg.Mongo.Database = viper.GetString("db.mongo.database")

	cfg.InMemory.DBFilePath = viper.GetString("db.inmemory.dbfilepath")

	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("error parsing config: %v", err)
	}

	return cfg
}
