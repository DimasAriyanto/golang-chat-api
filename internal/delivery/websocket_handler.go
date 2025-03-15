package delivery

import (
	"fmt"
	"log"
	"net/http"
	"time"

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

func wsHandler(w http.ResponseWriter, r *http.Request) {
	if !limiter.Allow() {
		http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
		return
	}

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

		log.Printf("Received: %s", msg)

		// Simpan pesan ke cache Redis
		err = redisCache.SetCache("latest_message", string(msg), 10*time.Minute)
		if err != nil {
			log.Println("Error caching message:", err)
		}

		if err := broker.PublishToQueue(string(msg)); err != nil {
			log.Println("Error publishing message to RabbitMQ:", err)
			continue
		}

		if err = conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("Message received at %v", time.Now()))); err != nil {
			log.Println("Error sending message:", err)
			break
		}
	}
}
