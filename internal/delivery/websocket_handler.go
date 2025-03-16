package delivery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
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

var activeConnections = struct {
	sync.RWMutex
	conns map[int]*websocket.Conn
}{conns: make(map[int]*websocket.Conn)}

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
	userIDInt := int(userID)
	log.Printf("User %d connected to WebSocket", userIDInt)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	activeConnections.Lock()
	if prevConn, exists := activeConnections.conns[userIDInt]; exists {
		prevConn.Close()
	}
	activeConnections.conns[userIDInt] = conn
	activeConnections.Unlock()

	defer func() {
		conn.Close()
		activeConnections.Lock()
		delete(activeConnections.conns, userIDInt)
		activeConnections.Unlock()
		log.Printf("User %d disconnected from WebSocket", userIDInt)
	}()

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

		messageData.SenderID = userIDInt
		messageData.Timestamp = time.Now().Format(time.RFC3339)

		messageJSON, _ := json.Marshal(messageData)

		log.Printf("Received from User %d to User %d: %s", userIDInt, messageData.ReceiverID, messageData.Message)

		err = redisCache.SetCache("latest_message", string(messageJSON), 10*time.Minute)
		if err != nil {
			log.Println("Error caching message:", err)
		}

		if err := broker.PublishToQueue(string(messageJSON)); err != nil {
			log.Println("Error publishing message to RabbitMQ:", err)
			continue
		}

		ackMsg := fmt.Sprintf("Message sent to User %d: %s", messageData.ReceiverID, messageData.Message)
		if err = conn.WriteMessage(websocket.TextMessage, []byte(ackMsg)); err != nil {
			log.Println("Error sending ack message:", err)
			break
		}

		DeliverMessageToUser(messageData)
	}
}

func DeliverMessageToUser(msg Message) {
	activeConnections.RLock()
	recipientConn, exists := activeConnections.conns[msg.ReceiverID]
	activeConnections.RUnlock()

	if exists {
		formattedMsg := fmt.Sprintf("New message from User %d: %s", msg.SenderID, msg.Message)
		if err := recipientConn.WriteMessage(websocket.TextMessage, []byte(formattedMsg)); err != nil {
			log.Printf("Error delivering message to User %d: %v", msg.ReceiverID, err)
		} else {
			log.Printf("Message delivered to User %d in real-time", msg.ReceiverID)
		}
	} else {
		log.Printf("User %d is not connected, message queued for later delivery", msg.ReceiverID)
	}
}
