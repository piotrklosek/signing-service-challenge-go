package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
)

type userRepo struct {
	mu       sync.RWMutex
	userData map[string]*domain.User
}

func NewUserRepo() *userRepo {
	return &userRepo{userData: make(map[string]*domain.User)}
}

func (r *userRepo) Create(ctx context.Context, u *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.userData[u.ID]; exists {
		return errors.New("user already exists")
	}
	r.userData[u.ID] = u
	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.userData[id]
	if !ok {
		return nil, errors.New("user not found")
	}
	return u, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, u := range r.userData {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("user not found")
}

func (r *userRepo) List(ctx context.Context) ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.User
	for _, u := range r.userData {
		list = append(list, u)
	}
	return list, nil
}
