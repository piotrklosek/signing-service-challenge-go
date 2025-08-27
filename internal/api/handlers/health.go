package handlers

import (
	"net/http"

	"github.com/piotrklosek/signing-service-challenge-go/internal/utils/jsonw"
)

type HealthHandler struct{}

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health evaluates the health of the service and writes a standardized response.
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	health := HealthResponse{
		Status:  "pass",
		Version: "v1",
	}
	jsonw.Success(w, health, http.StatusOK)

}
