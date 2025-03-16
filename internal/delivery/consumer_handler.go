package delivery

import (
	"encoding/json"
	"log"

	"github.com/DimasAriyanto/golang-chat-api/internal/domain"
	"github.com/DimasAriyanto/golang-chat-api/internal/repository"
	"github.com/streadway/amqp"
)

func StartConsumer(chatRepo *repository.ChatRepository, rabbitmqURL string) {
	conn, err := amqp.Dial(rabbitmqURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"chat_queue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Println("Consumer started. Waiting for messages...")

	go func() {
		for msg := range msgs {
			log.Printf("Received a message from RabbitMQ: %s", msg.Body)

			var chat domain.Chat
			if err := json.Unmarshal(msg.Body, &chat); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			if err := chatRepo.SaveMessage(chat); err != nil {
				log.Printf("Failed to save message to DB: %v", err)
			} else {
				log.Println("Message successfully saved to database.")
			}

			var wsMsg Message
			if err := json.Unmarshal(msg.Body, &wsMsg); err == nil {
				DeliverMessageToUser(wsMsg)
			}
		}
	}()

	select {}
}