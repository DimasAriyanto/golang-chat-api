package delivery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/DimasAriyanto/golang-chat-api/internal/middleware"
	"github.com/DimasAriyanto/golang-chat-api/pkg/broker"
	"github.com/DimasAriyanto/golang-chat-api/pkg/cache"
	"github.com/gorilla/websocket"
	"golang.org/x/time/rate"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

var limiter = rate.NewLimiter(1, 5)
var redisCache = cache.NewRedisCache("localhost:6379", "", 0)

type Message struct {
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Message    string `json:"message"`
	Timestamp  string `json:"timestamp"`
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		http.Error(w, "Unauthorized - Token required", http.StatusUnauthorized)
		return
	}

	claims, err := middleware.ParseToken(tokenString)
	if err != nil {
		http.Error(w, "Unauthorized - Invalid token", http.StatusUnauthorized)
		return
	}

	userID, ok := (*claims)["user_id"].(float64)
	if !ok {
		http.Error(w, "Unauthorized - Invalid token data", http.StatusUnauthorized)
		return
	}
	log.Printf("User %d connected to WebSocket", int(userID))

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}
	defer conn.Close()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}

		var messageData Message
		if err := json.Unmarshal(msg, &messageData); err != nil {
			log.Println("Invalid message format:", err)
			continue
		}

		messageData.SenderID = int(userID)
		messageData.Timestamp = time.Now().Format(time.RFC3339)

		messageJSON, _ := json.Marshal(messageData)

		log.Printf("Received from User %d: %s", int(userID), string(messageJSON))

		err = redisCache.SetCache("latest_message", string(messageJSON), 10*time.Minute)
		if err != nil {
			log.Println("Error caching message:", err)
		}

		if err := broker.PublishToQueue(string(messageJSON)); err != nil {
			log.Println("Error publishing message to RabbitMQ:", err)
			continue
		}

		ackMsg := fmt.Sprintf("User %d says: %s", messageData.SenderID, messageData.Message)
		if err = conn.WriteMessage(websocket.TextMessage, []byte(ackMsg)); err != nil {
			log.Println("Error sending message:", err)
			break
		}
	}
}
