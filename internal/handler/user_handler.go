package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/Nuttachai-K/cafe-inventory-management/internal/model"
	"github.com/Nuttachai-K/cafe-inventory-management/internal/service"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// GetUserByID godoc
// @Summary Get a user by id
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} model.User
// @Failure 400 {object} model.ErrorResponse "invalid user id"
// @Failure 404 {object} model.ErrorResponse "data not found"
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		writeJSONError(
			w,
			"Invalid user id",
			http.StatusBadRequest,
		)
		return
	}

	user, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		writeError(w, err)
	}
}

// GetAllUser godoc
// @Summary Get all users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param limit query int false "Max results to return (default 20)"
// @Success 200 {array} model.User
// @Failure 400 {object} model.ErrorResponse "invalid limit"
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Router /api/v1/users [get]
func (h *UserHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	limit := 20
	if l := r.URL.Query().Get("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil {
			writeJSONError(w, "invalid limit", http.StatusBadRequest)
			return
		}
		limit = parsed
	}

	users, err := h.service.GetAll(
		r.Context(),
		limit,
	)
	if err != nil {
		writeError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		writeError(w, err)
	}
}

// CreateUser godoc
// @Summary Create a user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param user body model.UserCreate true "User payload"
// @Success 201 {object} object{id=int,message=string}
// @Failure 400 {object} model.ErrorResponse "invalid request body"
// @Failure 409 {object} model.ErrorResponse "email is already used"
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Router /api/v1/users [post]
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var user model.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		writeJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.Create(r.Context(), &user); err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(struct {
		ID      int    `json:"id"`
		Message string `json:"message"`
	}{
		ID:      user.ID,
		Message: "User created successfully",
	}); err != nil {
		writeError(w, err)
	}
}

// UpdateUser godoc
// @Summary Update a user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param user body model.UserUpdate true "User payload"
// @Success 200 {object} model.User
// @Failure 400 {object} model.ErrorResponse "invalid request body or invalid user id"
// @Failure 409 {object} model.ErrorResponse "email is already used"
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Router /api/v1/users/{id} [patch]
func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSONError(w, "invalid id", http.StatusBadRequest)
		return
	}

	var req model.UserUpdate

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := h.service.Update(r.Context(), id, &req)
	if err != nil {
		writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		writeError(w, err)
	}
}

// DeleteUser godoc
// @Summary Delete a user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 204
// @Failure 400 {object} model.ErrorResponse "invalid user id"
// @Failure 404 {object} model.ErrorResponse "data not found"
// @Failure 500 {object} model.ErrorResponse "internal server error"
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		writeJSONError(w, "invalid id", http.StatusBadRequest)
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
