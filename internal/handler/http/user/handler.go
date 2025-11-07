package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/w-h-a/demo-go/api/user"
	httphandler "github.com/w-h-a/demo-go/internal/handler/http"
	userservice "github.com/w-h-a/demo-go/internal/service/user"
)

// userHandler is the HTTP handler for user-related requests.
// It is "trivial" code, responsible for:
// 1. Decoding HTTP requests
// 2. Calling the service layer
// 3. Encoding HTTP responses (or errors)
type userHandler struct {
	service *userservice.Service
}

// CreateUser handles the HTTP POST /api/users request.
func (h *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var dto user.CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		httphandler.WrtErr(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	user, err := h.service.CreateUser(r.Context(), dto)
	if err != nil {
		if errors.Is(err, userservice.ErrEmailInUse) {
			httphandler.WrtErr(w, http.StatusConflict, err.Error())
			return
		}
		if errors.Is(err, userservice.ErrInvalidInput) {
			httphandler.WrtErr(w, http.StatusBadRequest, err.Error())
			return
		}
		httphandler.WrtErr(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	httphandler.WrtJSON(w, http.StatusCreated, user)
}

// GetUserByID handles the HTTP GET /api/users/{id} request.
func (h *userHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := h.service.GetUser(r.Context(), id)
	if err != nil {
		if errors.Is(err, userservice.ErrUserNotFound) {
			httphandler.WrtErr(w, http.StatusNotFound, "User not found")
			return
		}
		httphandler.WrtErr(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	httphandler.WrtJSON(w, http.StatusOK, user)
}

// GetAllUsers handles the HTTP GET /api/users request.
func (h *userHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers(r.Context())
	if err != nil {
		log.Printf("Internal server error on GetAllUsers: %v", err)
		httphandler.WrtErr(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	httphandler.WrtJSON(w, http.StatusOK, users)
}

func New(s *userservice.Service) *userHandler {
	return &userHandler{service: s}
}
