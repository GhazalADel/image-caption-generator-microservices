package main

import (
	"ImageCaptioningService/handler"
	"ImageCaptioningService/services/QueueService/rabbitmq"
	"log"
)

func main() {
	rabbitMQ, err := rabbitmq.Connect()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	queueName := "requests"
	err = rabbitMQ.ReadFromQueue(queueName, handler.HandleRequest)
	if err != nil {
		log.Printf("Failed to read from queue: %v\n", err)
	}

	select {}
}
