// File path: /internal/usecase/chat_usecase.go
package usecase

import (
	"github.com/DimasAriyanto/golang-chat-api/internal/domain"
	"github.com/DimasAriyanto/golang-chat-api/internal/repository"
	"errors"
	"time"
)

type ChatUseCase struct {
	ChatRepo repository.ChatRepository
}

func NewChatUseCase(chatRepo repository.ChatRepository) *ChatUseCase {
	return &ChatUseCase{ChatRepo: chatRepo}
}

func (uc *ChatUseCase) SendMessage(chat domain.Chat) error {
	chat.Timestamp = time.Now()

	if chat.Sender == "" || chat.Message == "" {
		return errors.New("sender and message are required")
	}

	return uc.ChatRepo.PublishMessage(chat)
}

func (uc *ChatUseCase) GetLatestMessage() (domain.Chat, error) {
	return uc.ChatRepo.GetLatestMessage()
}
