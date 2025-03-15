// File path: /internal/delivery/http_handler.go
package delivery

import (
	"github.com/DimasAriyanto/golang-chat-api/internal/domain"
	"github.com/DimasAriyanto/golang-chat-api/internal/usecase"
	"github.com/DimasAriyanto/golang-chat-api/internal/middleware"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
	UserUC *usecase.UserUseCase
}

func NewUserHandler(userUC *usecase.UserUseCase) *UserHandler {
	return &UserHandler{UserUC: userUC}
}

func (h *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user domain.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.UserUC.RegisterUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := h.UserUC.LoginUser(credentials.Username, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	token, _ := middleware.GenerateToken(user.Username)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
