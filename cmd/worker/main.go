package main

import (
	"immersion-go/internal/queue"
	"os"
	"time"
)

func main() {
	qcfg := queue.RabbitMQConfig{
		URL:       os.Getenv("RABBITMQ_URL"),
		TopicName: os.Getenv("RABBITMQ_TOPIC_NAME"),
		Timeout:   time.Now(),
	}
}
