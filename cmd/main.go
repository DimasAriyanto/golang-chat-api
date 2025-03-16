package main

import (
	"github.com/DimasAriyanto/golang-chat-api/config"
	"github.com/DimasAriyanto/golang-chat-api/internal/delivery"
	"github.com/DimasAriyanto/golang-chat-api/internal/middleware"
	"github.com/DimasAriyanto/golang-chat-api/internal/repository"
	"github.com/DimasAriyanto/golang-chat-api/internal/usecase"
	"github.com/DimasAriyanto/golang-chat-api/pkg/cache"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig()

	db := config.ConnectDB(cfg)
	redisCache := cache.NewRedisCache(cfg.GetRedisAddress(), cfg.RedisPassword, cfg.RedisDB)

	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := delivery.NewUserHandler(userUseCase, cfg)

	chatRepo := repository.NewChatRepository(db, redisCache, cfg)
	chatUseCase := usecase.NewChatUseCase(chatRepo)
	chatHandler := delivery.NewChatHandler(chatUseCase)

	delivery.InitWebSocket(redisCache, cfg.RabbitMQURL, cfg.JWTSecret)

	limiterMiddleware := middleware.RateLimiterMiddleware(1, 5)

	http.Handle("/register", limiterMiddleware(http.HandlerFunc(userHandler.RegisterUser)))
	http.Handle("/login", limiterMiddleware(http.HandlerFunc(userHandler.LoginUser)))
	http.HandleFunc("/ws", delivery.WsHandler)
	http.Handle("/send-message", middleware.AuthMiddleware(cfg)(limiterMiddleware(http.HandlerFunc(chatHandler.SendMessageHandler))))
	http.Handle("/chat-history", middleware.AuthMiddleware(cfg)(limiterMiddleware(http.HandlerFunc(chatHandler.GetChatHistoryHandler))))

	go delivery.StartConsumer(chatRepo, cfg.RabbitMQURL)

	serverAddr := ":" + cfg.ServerPort
	log.Printf("Server started on %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatal("ListenAndServe Error:", err)
	}
}
