package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
	"github.com/piotrklosek/signing-service-challenge-go/mocks/database"
)

func TestCreateDevice_Success(t *testing.T) {
	userID := uuid.NewString()

	userRepo := &database.MockUserRepo{
		GetByIDFn: func(ctx context.Context, id string) (*domain.User, error) {
			return &domain.User{ID: userID, Name: "Alice"}, nil
		},
	}

	deviceRepo := &database.MockDeviceRepo{
		CreateFn: func(ctx context.Context, d *domain.SignatureDevice) error {
			return nil
		},
	}

	h := NewDeviceHandler(deviceRepo, userRepo)

	body := []byte(`{"user_id":"` + userID + `","algorithm":"RSA","label":"test_device"}`)
	req := httptest.NewRequest(http.MethodPost, "/devices", bytes.NewReader(body))
	w := httptest.NewRecorder()

	err := h.CreateDevice(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateDevice_InvalidJSON(t *testing.T) {
	h := NewDeviceHandler(&database.MockDeviceRepo{}, &database.MockUserRepo{})

	req := httptest.NewRequest(http.MethodPost, "/devices", bytes.NewReader([]byte("{bad-json")))
	w := httptest.NewRecorder()

	// expect error
	err := h.CreateDevice(w, req)
	if !strings.ContainsAny(err.Error(), ErrBadRequest.Error()) {
		t.Errorf("expected error[%v], got[%v]", ErrBadRequest.Error(), err.Error())
	}

}

func TestCreateDevice_UserNotFound(t *testing.T) {
	userRepo := &database.MockUserRepo{
		GetByIDFn: func(ctx context.Context, id string) (*domain.User, error) {
			return nil, errors.New("user not found")
		},
	}
	h := NewDeviceHandler(&database.MockDeviceRepo{}, userRepo)

	body := []byte(`{"user_id":"` + uuid.NewString() + `","algorithm":"RSA"}`)
	req := httptest.NewRequest(http.MethodPost, "/devices", bytes.NewReader(body))
	w := httptest.NewRecorder()

	// expect error
	err := h.CreateDevice(w, req)
	if !strings.ContainsAny(err.Error(), ErrUserNotFounc.Error()) {
		t.Errorf("expected error[%v], got[%v]", ErrUserNotFounc.Error(), err.Error())
	}
}

func TestCreateDevice_DeviceRepoError(t *testing.T) {
	userID := uuid.NewString()
	userRepo := &database.MockUserRepo{
		GetByIDFn: func(ctx context.Context, id string) (*domain.User, error) {
			return &domain.User{ID: id}, nil
		},
	}
	deviceRepo := &database.MockDeviceRepo{
		CreateFn: func(ctx context.Context, d *domain.SignatureDevice) error {
			return errors.New("db error")
		},
	}

	h := NewDeviceHandler(deviceRepo, userRepo)

	body := []byte(`{"user_id":"` + userID + `","algorithm":"RSA"}`)
	req := httptest.NewRequest(http.MethodPost, "/devices", bytes.NewReader(body))
	w := httptest.NewRecorder()

	err := h.CreateDevice(w, req)
	if !strings.ContainsAny(err.Error(), ErrCreatingDevice.Error()) {
		t.Errorf("expected error[%v], got[%v]", ErrCreatingDevice.Error(), err.Error())
	}

}

func TestListDevices_RepoError(t *testing.T) {
	deviceRepo := &database.MockDeviceRepo{
		ListFn: func(ctx context.Context) ([]*domain.SignatureDevice, error) {
			return nil, errors.New("db error")
		},
	}

	h := NewDeviceHandler(deviceRepo, &database.MockUserRepo{})

	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	w := httptest.NewRecorder()

	err := h.ListDevices(w, req)
	if !strings.ContainsAny(err.Error(), ErrListDevices.Error()) {
		t.Errorf("expected error[%v], got[%v]", ErrListDevices.Error(), err.Error())
	}

}

// helper type to decode response and check if list contains elements
type deviceListResponse struct {
	Status string                    `json:"status"`
	Data   []*domain.SignatureDevice `json:"data"`
}

func TestListDevices_Success(t *testing.T) {
	device := &domain.SignatureDevice{
		ID:        uuid.NewString(),
		UserID:    uuid.NewString(),
		Algorithm: domain.AlgorithmRSA,
		Label:     "Test Device",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	deviceRepo := &database.MockDeviceRepo{
		ListFn: func(ctx context.Context) ([]*domain.SignatureDevice, error) {
			return []*domain.SignatureDevice{device}, nil
		},
	}

	h := NewDeviceHandler(deviceRepo, &database.MockUserRepo{})

	req := httptest.NewRequest(http.MethodGet, "/devices", nil)
	w := httptest.NewRecorder()

	err := h.ListDevices(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	var resp deviceListResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("expected status=success, got %s", resp.Status)
	}
	if len(resp.Data) != 1 {
		t.Errorf("expected 1 device, got %d", len(resp.Data))
	}
}

// helper type to decode response and check device
type getDeviceResponse struct {
	Status string                  `json:"status"`
	Data   *domain.SignatureDevice `json:"data"`
}

func TestGetDevice_Success(t *testing.T) {
	device := &domain.SignatureDevice{
		ID:        "dev-1",
		UserID:    uuid.NewString(), // currently we are not checking if user exist
		Algorithm: domain.AlgorithmRSA,
		Label:     "My Device",
	}

	deviceRepo := &database.MockDeviceRepo{
		GetByIDFn: func(ctx context.Context, id string) (*domain.SignatureDevice, error) {
			return device, nil
		},
	}

	h := NewDeviceHandler(deviceRepo, &database.MockUserRepo{})

	req := httptest.NewRequest(http.MethodGet, "/devices/dev-1", nil)
	req.SetPathValue("id", "dev-1")

	w := httptest.NewRecorder()
	err := h.GetDevice(w, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", w.Result().StatusCode)
	}

	var resp getDeviceResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if resp.Status != "success" {
		t.Errorf("expected status=success, got %s", resp.Status)
	}

	if resp.Data.ID != device.ID && resp.Data.Label != device.Label && resp.Data.UserID != device.UserID {
		t.Errorf("expected device[%v]\n, got %v", device, resp.Data)
	}
}

func TestGetDevice_NotFound(t *testing.T) {
	deviceRepo := &database.MockDeviceRepo{
		GetByIDFn: func(ctx context.Context, id string) (*domain.SignatureDevice, error) {
			return nil, errors.New("device not found")
		},
	}

	h := NewDeviceHandler(deviceRepo, &database.MockUserRepo{})

	req := httptest.NewRequest(http.MethodGet, "/devices/xyz", nil)
	req.SetPathValue("id", "xyz")

	w := httptest.NewRecorder()
	err := h.GetDevice(w, req)
	if !strings.ContainsAny(err.Error(), ErrDeviceNotFounc.Error()) {
		t.Errorf("expected error[%v], got[%v]", ErrDeviceNotFounc.Error(), err.Error())
	}
	if w.Result().StatusCode != http.StatusNotFound {
		t.Errorf("expected 404, got %d", w.Result().StatusCode)
	}
}
