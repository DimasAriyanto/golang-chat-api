// File path: /pkg/broker/rabbitmq_broker.go
package broker

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func PublishToQueue(msg string) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return fmt.Errorf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("Failed to open a channel: %v", err)
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
		return fmt.Errorf("Failed to declare a queue: %v", err)
	}

	err = ch.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
	if err != nil {
		return fmt.Errorf("Failed to publish a message: %v", err)
	}

	log.Printf("[x] Sent: %s", msg)
	return nil
}
