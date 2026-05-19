package rabbitmq

import (
	config "RequestManagementService/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQ struct {
	Conn    *amqp.Connection
	Channel *amqp.Channel
}

func Connect() (*RabbitMQ, error) {
	cfg, err := config.NewConfig()
	if err != nil {
		return nil, err
	}

	amqpURL := cfg.RabbitQueue.Url
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		closeChannelErr := conn.Close()
		if closeChannelErr != nil {
			return nil, closeChannelErr
		}
		return nil, err
	}

	return &RabbitMQ{
		Conn:    conn,
		Channel: ch,
	}, nil
}

func (r *RabbitMQ) Close() {
	if r.Channel != nil {
		err := r.Channel.Close()
		if err != nil {
			return
		}
	}
	if r.Conn != nil {
		err := r.Conn.Close()
		if err != nil {
			return
		}
	}
}
