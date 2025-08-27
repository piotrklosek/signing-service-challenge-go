package database

import (
	"context"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
)

// MockUserRepo impelemnting UserRepository
type MockUserRepo struct {
	CreateFn     func(ctx context.Context, u *domain.User) error
	GetByIDFn    func(ctx context.Context, id string) (*domain.User, error)
	GetByEmailFn func(ctx context.Context, email string) (*domain.User, error)
	ListFn       func(ctx context.Context) ([]*domain.User, error)
}

// Create run func or return nil
func (m *MockUserRepo) Create(ctx context.Context, u *domain.User) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, u)
	}
	return nil
}

// GetByID run func or return nil
func (m *MockUserRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

// GetByEmail run func or return nil
func (m *MockUserRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	if m.GetByEmailFn != nil {
		return m.GetByEmailFn(ctx, email)
	}
	return nil, nil
}

// List run func or return empty list
func (m *MockUserRepo) List(ctx context.Context) ([]*domain.User, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx)
	}
	return []*domain.User{}, nil
}
