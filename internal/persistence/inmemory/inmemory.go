package inmemory

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"sync"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
)

// MemoryStore represents inmemory layer object to enable saving/loading data
type MemoryStore struct {
	UserRepo      *userRepo
	DeviceRepo    *deviceRepo
	SignatureRepo *signatureRepo

	mu     sync.RWMutex
	dbFile string
}

// NewMemoryStore create a new inmemory layer, if previous dbfile exist loads data
func NewMemoryStore(dbFile string) (*MemoryStore, error) {
	store := &MemoryStore{
		UserRepo:      NewUserRepo(),
		DeviceRepo:    NewDeviceRepo(),
		SignatureRepo: NewSignatureRepo(),
		dbFile:        dbFile,
	}

	if dbFile != "" {
		if err := store.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	return store, nil
}

// Save dump inmemory layer into json file
func (s *MemoryStore) Save() error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	dump := struct {
		Users      map[string]*domain.User              `json:"users"`
		Devices    map[string]*domain.SignatureDevice   `json:"devices"`
		Signatures map[string][]*domain.SignatureRecord `json:"signatures"`
	}{
		Users:      s.UserRepo.userData,
		Devices:    s.DeviceRepo.deviceData,
		Signatures: s.SignatureRepo.signaturesData,
	}

	data, err := json.MarshalIndent(dump, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(s.dbFile, data, 0644)
}

// Load dumped data into inmemory layer
func (s *MemoryStore) Load() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := os.ReadFile(s.dbFile)
	if err != nil {
		return err
	}

	var dump struct {
		Users      map[string]*domain.User              `json:"users"`
		Devices    map[string]*domain.SignatureDevice   `json:"devices"`
		Signatures map[string][]*domain.SignatureRecord `json:"signatures"`
	}
	if err := json.Unmarshal(data, &dump); err != nil {
		return err
	}

	s.UserRepo.userData = dump.Users
	s.DeviceRepo.deviceData = dump.Devices
	s.SignatureRepo.signaturesData = dump.Signatures

	return nil
}

// SaveOnShutdown triggers saving data of inmemory persistence layer
func (s *MemoryStore) SaveOnShutdown(ctx context.Context) {
	<-ctx.Done()
	_ = s.Save()
}
