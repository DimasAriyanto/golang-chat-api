package main

import (
	"github.com/DimasAriyanto/golang-chat-api/config"
	"github.com/DimasAriyanto/golang-chat-api/internal/delivery"
	"github.com/DimasAriyanto/golang-chat-api/internal/repository"
	"github.com/DimasAriyanto/golang-chat-api/internal/usecase"
	"log"
	"net/http"
)

func main() {
	db := config.ConnectDB()

	userRepo := repository.NewUserRepository(db)
	userUseCase := usecase.NewUserUseCase(userRepo)
	userHandler := delivery.NewUserHandler(userUseCase)

	http.HandleFunc("/register", userHandler.RegisterUser)
	http.HandleFunc("/login", userHandler.LoginUser)

	log.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe Error:", err)
	}
}
