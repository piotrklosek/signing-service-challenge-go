package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/piotrklosek/signing-service-challenge-go/internal/domain"
	"github.com/piotrklosek/signing-service-challenge-go/internal/persistence"
	"github.com/piotrklosek/signing-service-challenge-go/internal/utils/jsonw"
	"github.com/piotrklosek/signing-service-challenge-go/internal/validation"
)

type UserHandler struct {
	userRepo persistence.UserRepository
}

func NewUserHandler(userRepo persistence.UserRepository) *UserHandler {
	return &UserHandler{userRepo: userRepo}
}

type CreateUserRequest struct {
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}

type ListUserRequest struct{}
type ListUsersResponse struct {
	Users []domain.User `json:"users"`
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonw.Error(w, "invalid json", nil, http.StatusBadRequest)
		return err
	}
	if err := validation.ValidateStruct(&req); err != nil {
		jsonw.Error(w, err.Error(), nil, http.StatusBadRequest)
		return err
	}

	user := &domain.User{
		ID:        uuid.NewString(),
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := h.userRepo.Create(r.Context(), user); err != nil {
		jsonw.Error(w, err.Error(), nil, http.StatusInternalServerError)
		return err
	}

	jsonw.Success(w, user, http.StatusCreated)
	return nil
}

func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) error {
	users, err := h.userRepo.List(r.Context())
	if err != nil {
		jsonw.Error(w, err.Error(), nil, http.StatusInternalServerError)
		return err
	}
	jsonw.Success(w, users, http.StatusOK)
	return nil
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) error {
	id := r.PathValue("id")
	user, err := h.userRepo.GetByID(r.Context(), id)
	if err != nil {
		jsonw.Error(w, "user not found", nil, http.StatusNotFound)
		return err
	}
	jsonw.Success(w, user, http.StatusOK)
	return nil
}
