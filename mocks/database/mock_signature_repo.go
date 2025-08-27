package database

import (
	"context"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
)

// MockSignatureRepo implement SignatureRepository
type MockSignatureRepo struct {
	CreateFn       func(ctx context.Context, s *domain.SignatureRecord) error
	ListByDeviceFn func(ctx context.Context, deviceID string) ([]*domain.SignatureRecord, error)
}

// Create run func or return nil
func (m *MockSignatureRepo) Create(ctx context.Context, s *domain.SignatureRecord) error {
	if m.CreateFn != nil {
		return m.CreateFn(ctx, s)
	}
	return nil
}

// ListByDevice run func or return empty list
func (m *MockSignatureRepo) ListByDevice(ctx context.Context, deviceID string) ([]*domain.SignatureRecord, error) {
	if m.ListByDeviceFn != nil {
		return m.ListByDeviceFn(ctx, deviceID)
	}
	return []*domain.SignatureRecord{}, nil
}
