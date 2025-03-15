package config

import "os"

type Config struct {
	ServerPort string
	RabbitMQURL string
}

func LoadConfig() Config {
	return Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		RabbitMQURL: getEnv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
	}
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
