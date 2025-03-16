package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/DimasAriyanto/golang-chat-api/internal/domain"
	"github.com/DimasAriyanto/golang-chat-api/pkg/broker"
	"github.com/DimasAriyanto/golang-chat-api/pkg/cache"
)

type ChatRepository struct {
	DB    *sql.DB
	Cache *cache.RedisCache
}

func NewChatRepository(db *sql.DB, cache *cache.RedisCache) *ChatRepository {
	return &ChatRepository{DB: db, Cache: cache}
}

func (r *ChatRepository) SaveMessage(chat domain.Chat) error {
	query := `INSERT INTO chat (sender_id, receiver_id, group_id, message, timestamp)
              VALUES (?, ?, ?, ?, ?)`
	_, err := r.DB.Exec(query, chat.SenderID, chat.ReceiverID, chat.GroupID, chat.Message, chat.Timestamp)
	return err
}

func (r *ChatRepository) PublishMessage(chat domain.Chat) error {
	data, err := json.Marshal(chat)
	if err != nil {
		return fmt.Errorf("failed to marshal chat: %v", err)
	}

	return broker.PublishToQueue(string(data))
}

func (r *ChatRepository) GetCachedMessages(userID int) ([]domain.Chat, error) {
	cacheKey := fmt.Sprintf("user:%d:chat_history", userID)
	cachedData, err := r.Cache.GetCache(cacheKey)
	if err != nil {
		return nil, err
	}

	var chats []domain.Chat
	if err := json.Unmarshal([]byte(cachedData), &chats); err != nil {
		return nil, err
	}

	return chats, nil
}

func (r *ChatRepository) GetMessagesByUserID(userID int) ([]domain.Chat, error) {
	query := `SELECT id, sender_id, receiver_id, message, timestamp
              FROM chat WHERE sender_id = ? OR receiver_id = ? ORDER BY timestamp DESC`

	rows, err := r.DB.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chats []domain.Chat
	for rows.Next() {
		var chat domain.Chat
		if err := rows.Scan(&chat.ID, &chat.SenderID, &chat.ReceiverID, &chat.Message, &chat.Timestamp); err != nil {
			return nil, err
		}
		chats = append(chats, chat)
	}

	jsonData, err := json.Marshal(chats)
	if err == nil {
		r.Cache.SetCache(fmt.Sprintf("user:%d:chat_history", userID), string(jsonData), 10*time.Minute)
	}

	return chats, nil
}

func (r *ChatRepository) GetUserByID(userID int) (domain.User, error) {
	var user domain.User
	query := `SELECT id, username, email FROM users WHERE id = ?`
	row := r.DB.QueryRow(query, userID)

	err := row.Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

