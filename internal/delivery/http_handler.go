package delivery

import (
	"github.com/DimasAriyanto/golang-chat-api/internal/domain"
	"github.com/DimasAriyanto/golang-chat-api/internal/usecase"
	"github.com/DimasAriyanto/golang-chat-api/internal/middleware"
	"github.com/DimasAriyanto/golang-chat-api/config"
	"encoding/json"
	"net/http"
)

type UserHandler struct {
    UserUC *usecase.UserUseCase
    Config config.Config
}

func NewUserHandler(userUC *usecase.UserUseCase, cfg config.Config) *UserHandler {
    return &UserHandler{UserUC: userUC, Config: cfg}
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

    token, _ := middleware.GenerateToken(user.ID, h.Config.JWTSecret)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
