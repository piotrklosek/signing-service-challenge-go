package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
	"github.com/piotrklosek/signing-service-challenge-go/internal/persistence"
	"github.com/piotrklosek/signing-service-challenge-go/internal/utils/jsonw"
)

// SignatureHandler represents object used to create signatures
type SignatureHandler struct {
	signatureRepo persistence.SignatureRepository
	deviceRepo    persistence.DeviceRepository
}

// NewSignatureHandler used to create signature handler
func NewSignatureHandler(signatureRepo persistence.SignatureRepository, deviceRepo persistence.DeviceRepository) *SignatureHandler {
	return &SignatureHandler{signatureRepo: signatureRepo, deviceRepo: deviceRepo}
}

type SignTransactionRequest struct {
	Data string `json:"data"`
}

// SignTransactionData and return signature and signed data, and updates devices details about signatures
func (h *SignatureHandler) SignTransactionData(w http.ResponseWriter, r *http.Request) error {
	deviceID := r.PathValue("id")
	device, err := h.deviceRepo.GetByID(r.Context(), deviceID)
	if err != nil {
		jsonw.Error(w, "device not found", nil, http.StatusNotFound)
		return err
	}

	var req SignTransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonw.Error(w, "invalid json", nil, http.StatusBadRequest)
		return err
	}

	// sign data using domain logic
	signedData, signature, err := domain.SignData(device, req.Data)
	if err != nil {
		jsonw.Error(w, "signing failed: "+err.Error(), nil, http.StatusInternalServerError)
		return err
	}

	// createding and saving signature record details
	record := &domain.SignatureRecord{
		ID:         uuid.NewString(),
		DeviceID:   device.ID,
		SignedData: signedData,
		Signature:  signature,
		CreatedAt:  time.Now(),
	}
	if err := h.signatureRepo.Create(r.Context(), record); err != nil {
		jsonw.Error(w, err.Error(), nil, http.StatusInternalServerError)
		return err
	}

	// update device counter and last signature
	device.SignatureCounter++
	device.LastSignature = signature
	device.UpdatedAt = time.Now()
	if err := h.deviceRepo.Update(r.Context(), device); err != nil {
		jsonw.Error(w, "failed to update device counter", nil, http.StatusInternalServerError)
		return err
	}

	jsonw.Success(w, map[string]string{
		"signature":   signature,
		"signed_data": signedData,
	}, http.StatusCreated)

	return nil
}

// ListSignatures used to list all signatures recorded in system
func (h *SignatureHandler) ListSignatures(w http.ResponseWriter, r *http.Request) error {
	deviceID := r.PathValue("id")

	records, err := h.signatureRepo.ListByDevice(r.Context(), deviceID)
	if err != nil {
		jsonw.Error(w, err.Error(), nil, http.StatusInternalServerError)
		return err
	}
	jsonw.Success(w, records, http.StatusOK)
	return nil
}
