package postgres

import (
	"context"
	"database/sql"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
)

type userRepo struct {
	db *sql.DB
}

// NewUSerRepo create interface for user database
func NewUserRepo(db *sql.DB) *userRepo {
	return &userRepo{db: db}
}

// Create method used to created row object in with user
func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO users (id, name, email, created_at, updated_at)
         VALUES ($1, $2, $3, $4, $5)`,
		u.ID, u.Name, u.Email, u.CreatedAt, u.UpdatedAt,
	)
	return err
}

// GetByID used to return user by ID
func (r *userRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, created_at, updated_at FROM users WHERE id=$1`, id)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

// GetByEmail used to return user by email
func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, name, email, created_at, updated_at FROM users WHERE email=$1`, email)
	var u domain.User
	if err := row.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return &u, nil
}

// List used to return all user list
func (r *userRepo) List(ctx context.Context) ([]*domain.User, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, email, created_at, updated_at FROM users ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var u domain.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, &u)
	}
	return users, nil
}
