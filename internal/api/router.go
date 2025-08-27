package api

import (
	"net/http"

	"github.com/piotrklosek/signing-service-challenge-go/internal/api/handlers"
	"github.com/piotrklosek/signing-service-challenge-go/internal/persistence"
	"github.com/piotrklosek/signing-service-challenge-go/internal/utils/logger"
)

// TODO move handlers to separeted packages, so every handler can have different set of errors
// NewRouter creates HTTP router with endpoints
func NewRouter(
	deviceRepo persistence.DeviceRepository,
	signatureRepo persistence.SignatureRepository,
	userRepo persistence.UserRepository,
) http.Handler {

	mux := http.NewServeMux()
	apiLogger := logger.CreateLogger("api")

	// Handlers
	healthHandler := handlers.NewHealthHandler()
	deviceHandler := handlers.NewDeviceHandler(deviceRepo, userRepo)
	signatureHandler := handlers.NewSignatureHandler(signatureRepo, deviceRepo)
	userHandler := handlers.NewUserHandler(userRepo)

	// Devices
	mux.Handle("POST /api/v1/devices", middleware(apiLogger, deviceHandler.CreateDevice))
	mux.Handle("GET /api/v1/devices", middleware(apiLogger, deviceHandler.ListDevices))
	mux.Handle("GET /api/v1/devices/{id}", middleware(apiLogger, deviceHandler.GetDevice))

	// Signatures
	mux.Handle("POST /api/v1/devices/{id}/sign", middleware(apiLogger, signatureHandler.SignTransactionData))
	mux.Handle("GET /api/v1/devices/{id}/signatures", middleware(apiLogger, signatureHandler.ListSignatures))

	// Users as an extra for user management
	mux.Handle("POST /api/v1/users", middleware(apiLogger, userHandler.CreateUser))
	mux.Handle("GET /api/v1/users", middleware(apiLogger, userHandler.ListUsers))
	mux.Handle("GET /api/v1/users/{id}", middleware(apiLogger, userHandler.GetUser))

	// Health
	mux.HandleFunc("GET /health", healthHandler.Health)

	// TODO wrap endpoints with auth middleware for only admin access

	return mux
}
