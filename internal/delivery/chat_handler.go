package delivery

import (
	"encoding/json"
	"net/http"

	"github.com/DimasAriyanto/golang-chat-api/internal/domain"
	"github.com/DimasAriyanto/golang-chat-api/internal/usecase"
	"github.com/DimasAriyanto/golang-chat-api/internal/middleware"
)

type ChatHandler struct {
	ChatUC *usecase.ChatUseCase
}

func NewChatHandler(chatUC *usecase.ChatUseCase) *ChatHandler {
	return &ChatHandler{ChatUC: chatUC}
}

func (h *ChatHandler) SendMessageHandler(w http.ResponseWriter, r *http.Request) {
	senderID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized - Invalid user data", http.StatusUnauthorized)
		return
	}

	var chat domain.Chat
	if err := json.NewDecoder(r.Body).Decode(&chat); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	chat.SenderID = senderID

	if chat.ReceiverID == 0 {
		http.Error(w, "Receiver ID is required", http.StatusBadRequest)
		return
	}

	isValidReceiver, err := h.ChatUC.IsValidUser(chat.ReceiverID)
	if err != nil || !isValidReceiver {
		http.Error(w, "Invalid receiver ID", http.StatusBadRequest)
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
	userID, ok := r.Context().Value(middleware.UserIDKey).(int)
	if !ok {
		http.Error(w, "Unauthorized - Invalid user data", http.StatusUnauthorized)
		return
	}

	chats, err := h.ChatUC.GetChatHistory(userID)
	if err != nil {
		http.Error(w, "Failed to fetch chat history", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(chats)
}
