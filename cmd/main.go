package main

import (
	"github.com/DimasAriyanto/golang-chat-api/config"
	"github.com/DimasAriyanto/golang-chat-api/internal/delivery"
	"github.com/DimasAriyanto/golang-chat-api/internal/repository"
	"github.com/DimasAriyanto/golang-chat-api/internal/usecase"
	"github.com/DimasAriyanto/golang-chat-api/pkg/cache"
	"github.com/DimasAriyanto/golang-chat-api/internal/middleware"
	"log"
	"net/http"
)

func main() {
	db := config.ConnectDB()

	redisCache := cache.NewRedisCache("localhost:6379", "", 0)

	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := delivery.NewUserHandler(userUseCase)
	chatRepo := repository.NewChatRepository(db, redisCache)
	chatUseCase := usecase.NewChatUseCase(chatRepo)
	chatHandler := delivery.NewChatHandler(chatUseCase)

	http.HandleFunc("/register", userHandler.RegisterUser)
	http.HandleFunc("/login", userHandler.LoginUser)
	http.HandleFunc("/ws", delivery.WsHandler)
	http.Handle("/send-message", middleware.AuthMiddleware(http.HandlerFunc(chatHandler.SendMessageHandler)))
	http.Handle("/chat-history", middleware.AuthMiddleware(http.HandlerFunc(chatHandler.GetChatHistoryHandler)))

	go delivery.StartConsumer(chatRepo)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe Error:", err)
	}
}
