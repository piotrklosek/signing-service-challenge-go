package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/piotrklosek/signing-service-challenge-go/internal/persistence"
)

// PostgresStore used to store database reference
type PostgresStore struct {
	db *sql.DB
}

// NewStore create new postgress connection to use in repositories
func NewStore(dsn string) (*PostgresStore, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping postgres: %w", err)
	}

	return &PostgresStore{db: db}, nil
}

// Close use to close database connection
func (s *PostgresStore) Close() error {
	return s.db.Close()
}

// NewRepositories used to create repositories for postgres database connector
func NewRepositories(dsn string) (
	persistence.DeviceRepository,
	persistence.SignatureRepository,
	persistence.UserRepository,
	error,
) {
	store, err := NewStore(dsn)
	if err != nil {
		return nil, nil, nil, err
	}

	deviceRepo := NewDeviceRepo(store.db)
	signatureRepo := NewSignatureRepo(store.db)
	userRepo := NewUserRepo(store.db)

	return deviceRepo, signatureRepo, userRepo, nil
}

// TODO move into migrations to use golang migration tool
// RunMigrations used to create database schema at the start of app
func RunMigrations(dsn string) error {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS signature_devices (
			id UUID PRIMARY KEY,
			user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			algorithm TEXT NOT NULL,
			label TEXT,
			public_key TEXT,
			private_key TEXT,
			signature_counter BIGINT NOT NULL DEFAULT 0,
			last_signature TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);`,

		`CREATE TABLE IF NOT EXISTS signatures (
			id UUID PRIMARY KEY,
			device_id UUID NOT NULL REFERENCES signature_devices(id) ON DELETE CASCADE,
			signed_data TEXT NOT NULL,
			signature TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		);`,
	}

	for _, q := range queries {
		if _, err := db.ExecContext(context.Background(), q); err != nil {
			return fmt.Errorf("migration failed: %w", err)
		}
	}

	return nil
}
