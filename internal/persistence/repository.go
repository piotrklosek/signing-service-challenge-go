package persistence

import (
	"context"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, u *domain.User) error
	GetByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	List(ctx context.Context) ([]*domain.User, error)
}

type DeviceRepository interface {
	Create(ctx context.Context, d *domain.SignatureDevice) error
	GetByID(ctx context.Context, id string) (*domain.SignatureDevice, error)
	List(ctx context.Context) ([]*domain.SignatureDevice, error)
	Update(ctx context.Context, d *domain.SignatureDevice) error
}

type SignatureRepository interface {
	Create(ctx context.Context, s *domain.SignatureRecord) error
	ListByDevice(ctx context.Context, deviceID string) ([]*domain.SignatureRecord, error)
}
