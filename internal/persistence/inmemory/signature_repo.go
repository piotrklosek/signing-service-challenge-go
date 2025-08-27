package inmemory

import (
	"context"
	"errors"
	"sync"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
)

type signatureRepo struct {
	mu             sync.RWMutex
	signaturesData map[string][]*domain.SignatureRecord // deviceID -> list of signatures
}

func NewSignatureRepo() *signatureRepo {
	return &signatureRepo{signaturesData: make(map[string][]*domain.SignatureRecord)}
}

func (r *signatureRepo) Create(ctx context.Context, s *domain.SignatureRecord) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if s.DeviceID == "" {
		return errors.New("signature must have device id")
	}
	r.signaturesData[s.DeviceID] = append(r.signaturesData[s.DeviceID], s)
	return nil
}

func (r *signatureRepo) ListByDevice(ctx context.Context, deviceID string) ([]*domain.SignatureRecord, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.signaturesData[deviceID], nil
}
