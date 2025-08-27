package database

import (
	"context"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
)

// MockDeviceRepo implement DeviceRepository
type MockDeviceRepo struct {
	CreateFn  func(ctx context.Context, d *domain.SignatureDevice) error
	GetByIDFn func(ctx context.Context, id string) (*domain.SignatureDevice, error)
	ListFn    func(ctx context.Context) ([]*domain.SignatureDevice, error)
	UpdateFn  func(ctx context.Context, d *domain.SignatureDevice) error
}

// Create runs func or return nil
func (m *MockDeviceRepo) Create(ctx context.Context, d *domain.SignatureDevice) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, d)
	}
	return nil
}

// GetByID runs func or return
func (m *MockDeviceRepo) GetByID(ctx context.Context, id string) (*domain.SignatureDevice, error) {
	if m.GetByIDFn != nil {
		return m.GetByIDFn(ctx, id)
	}
	return nil, nil
}

// List runs func or return nil
func (m *MockDeviceRepo) List(ctx context.Context) ([]*domain.SignatureDevice, error) {
	if m.ListFn != nil {
		return m.ListFn(ctx)
	}
	return []*domain.SignatureDevice{}, nil
}

// Update runs func or return nil
func (m *MockDeviceRepo) Update(ctx context.Context, d *domain.SignatureDevice) error {
	if m.UpdateFn != nil {
		return m.UpdateFn(ctx, d)
	}
	return nil
}
