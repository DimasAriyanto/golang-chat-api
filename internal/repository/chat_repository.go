package repository

import (
	"github.com/DimasAriyanto/golang-chat-api/internal/domain"
	"github.com/DimasAriyanto/golang-chat-api/pkg/broker"
	"github.com/DimasAriyanto/golang-chat-api/pkg/cache"
	"encoding/json"
	"fmt"
)

type ChatRepository struct {
	Cache cache.RedisCache
}

func NewChatRepository(cache cache.RedisCache) *ChatRepository {
	return &ChatRepository{Cache: cache}
}

func (r *ChatRepository) PublishMessage(chat domain.Chat) error {
	data, err := json.Marshal(chat)
	if err != nil {
		return fmt.Errorf("failed to marshal chat: %v", err)
	}

	return broker.PublishToQueue(string(data))
}

func (r *ChatRepository) GetLatestMessage() (domain.Chat, error) {
	data, err := r.Cache.GetCache("latest_message")
	if err != nil {
		return domain.Chat{}, fmt.Errorf("failed to get latest message from cache: %v", err)
	}

	var chat domain.Chat
	if err := json.Unmarshal([]byte(data), &chat); err != nil {
		return domain.Chat{}, fmt.Errorf("failed to unmarshal chat: %v", err)
	}

	return chat, nil
}
