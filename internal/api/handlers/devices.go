package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
	"github.com/piotrklosek/signing-service-challenge-go/internal/persistence"
	"github.com/piotrklosek/signing-service-challenge-go/internal/utils/jsonw"
	"github.com/piotrklosek/signing-service-challenge-go/internal/validation"
)

type DeviceHandler struct {
	deviceRepo persistence.DeviceRepository
	userRepo   persistence.UserRepository
}

func NewDeviceHandler(deviceRepo persistence.DeviceRepository, userRepo persistence.UserRepository) *DeviceHandler {
	return &DeviceHandler{deviceRepo: deviceRepo, userRepo: userRepo}
}

type CreateDeviceRequest struct {
	UserID    string `json:"user_id" validate:"required,uuid4"`
	Algorithm string `json:"algorithm" validate:"required,algorithm"`
	Label     string `json:"label" validate:"omitempty,min=3,max=100"`
}

func (h *DeviceHandler) CreateDevice(w http.ResponseWriter, r *http.Request) error {
	var req CreateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonw.Error(w, "invalid json", nil, http.StatusBadRequest)
		return fmt.Errorf("%v - %v", ErrInvalidJson, err)
	}
	if err := validation.ValidateStruct(&req); err != nil {
		jsonw.Error(w, err.Error(), nil, http.StatusBadRequest)
		return err
	}

	// TODO move to middelware layer
	// check if user exists
	if _, err := h.userRepo.GetByID(r.Context(), req.UserID); err != nil {
		jsonw.Error(w, "user not found", nil, http.StatusBadRequest)
		return fmt.Errorf("%v - %v", ErrUserNotFounc, err)
	}

	device := &domain.SignatureDevice{
		ID:               uuid.NewString(),
		UserID:           req.UserID,
		Algorithm:        domain.AlgorithmType(req.Algorithm),
		Label:            req.Label,
		SignatureCounter: 0,
		LastSignature:    "", // fill during signing document
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := device.GenerateKeys(); err != nil {
		jsonw.Error(w, err.Error(), nil, http.StatusInternalServerError)
		return fmt.Errorf("%v - %v", ErrGenerateKeys, err)
	}

	if err := h.deviceRepo.Create(r.Context(), device); err != nil {
		jsonw.Error(w, err.Error(), nil, http.StatusInternalServerError)
		return fmt.Errorf("%v - %v", ErrCreatingDevice, err)
	}
	jsonw.Success(w, device, http.StatusCreated)
	return nil
}

func (h *DeviceHandler) ListDevices(w http.ResponseWriter, r *http.Request) error {
	devices, err := h.deviceRepo.List(r.Context())
	if err != nil {
		jsonw.Error(w, err.Error(), nil, http.StatusInternalServerError)
		return fmt.Errorf("%v - %v", ErrListDevices, err)
	}
	jsonw.Success(w, devices, http.StatusOK)
	return nil
}

func (h *DeviceHandler) GetDevice(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	device, err := h.deviceRepo.GetByID(r.Context(), id)
	if err != nil {
		jsonw.Error(w, "device not found", nil, http.StatusNotFound)
		return fmt.Errorf("%v - %v", ErrListDevices, err)
	}
	jsonw.Success(w, device, http.StatusOK)
	return nil
}
