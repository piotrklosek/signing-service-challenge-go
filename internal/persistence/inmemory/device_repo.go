package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
)

type deviceRepo struct {
	mu         sync.RWMutex
	deviceData map[string]*domain.SignatureDevice
}

func NewDeviceRepo() *deviceRepo {
	return &deviceRepo{deviceData: make(map[string]*domain.SignatureDevice)}
}

func (r *deviceRepo) Create(ctx context.Context, d *domain.SignatureDevice) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.deviceData[d.ID]; exists {
		return errors.New("device already exists")
	}
	r.deviceData[d.ID] = d
	return nil
}

func (r *deviceRepo) GetByID(ctx context.Context, id string) (*domain.SignatureDevice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	d, ok := r.deviceData[id]
	if !ok {
		return nil, errors.New("device not found")
	}
	return d, nil
}

func (r *deviceRepo) List(ctx context.Context) ([]*domain.SignatureDevice, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*domain.SignatureDevice
	for _, d := range r.deviceData {
		list = append(list, d)
	}
	return list, nil
}

func (r *deviceRepo) Update(ctx context.Context, d *domain.SignatureDevice) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.deviceData[d.ID]; !ok {
		return errors.New("device not found")
	}
	r.deviceData[d.ID] = d
	return nil
}
