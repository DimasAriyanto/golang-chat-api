package delivery

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DimasAriyanto/golang-chat-api/internal/domain"
	"github.com/DimasAriyanto/golang-chat-api/internal/usecase"
)

type ChatHandler struct {
	ChatUC *usecase.ChatUseCase
}

func NewChatHandler(chatUC *usecase.ChatUseCase) *ChatHandler {
	return &ChatHandler{ChatUC: chatUC}
}

func (h *ChatHandler) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	var chat domain.Chat
	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.ChatUC.SendMessage(chat); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Message sent successfully"})
}

func (h *ChatHandler) GetChatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.Atoi(r.URL.Query().Get("user_id"))
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	chats, err := h.ChatUC.GetChatHistory(userID)
	if err != nil {
		http.Error(w, "Failed to fetch chat history", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(chats)
}
