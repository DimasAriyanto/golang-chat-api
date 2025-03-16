package usecase

import (
	"errors"
	"time"

	"github.com/DimasAriyanto/golang-chat-api/internal/domain"
	"github.com/DimasAriyanto/golang-chat-api/internal/repository"
)

type ChatUseCase struct {
	ChatRepo *repository.ChatRepository
}

func NewChatUseCase(chatRepo *repository.ChatRepository) *ChatUseCase {
	return &ChatUseCase{ChatRepo: chatRepo}
}

func (uc *ChatUseCase) SendMessage(chat domain.Chat) error {
	if chat.Message == "" {
		return errors.New("message cannot be empty")
	}
	if chat.SenderID == 0 {
		return errors.New("invalid sender ID")
	}

	chat.Timestamp = time.Now()

	if err := uc.ChatRepo.SaveMessage(chat); err != nil {
		return err
	}

	return uc.ChatRepo.PublishMessage(chat)
}

func (uc *ChatUseCase) GetChatHistory(userID int) ([]domain.Chat, error) {
	chats, err := uc.ChatRepo.GetCachedMessages(userID)
	if err == nil && len(chats) > 0 {
		return chats, nil
	}

	return uc.ChatRepo.GetMessagesByUserID(userID)
}
