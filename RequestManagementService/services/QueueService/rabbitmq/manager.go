package rabbitmq

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func (r *RabbitMQ) AddToQueue(queueName string, body string) error {
	q, err := r.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	err = r.Channel.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(body),
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf(" [x] Sent %s to queue %s", body, queueName)
	return nil
}

func (r *RabbitMQ) ReadFromQueue(queueName string, handler func(amqp.Delivery)) error {
	q, err := r.Channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	messages, err := r.Channel.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range messages {
			handler(d)
		}
	}()

	return nil
}
