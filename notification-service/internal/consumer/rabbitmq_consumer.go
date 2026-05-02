package consumer

import (
	"encoding/json"
	"log"

	"github.com/rabbitmq/amqp091-go"

	"notification-service/internal/domain"
	"notification-service/internal/usecase"
)

type RabbitMQConsumer struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
	uc      *usecase.NotificationUsecase
}

func NewRabbitMQConsumer(rabbitURL string, uc *usecase.NotificationUsecase) (*RabbitMQConsumer, error) {
	conn, err := amqp091.Dial(rabbitURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(
		"payment.completed",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer{
		conn:    conn,
		channel: ch,
		uc:      uc,
	}, nil
}

func (c *RabbitMQConsumer) Start() error {
	msgs, err := c.channel.Consume(
		"payment.completed",
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	log.Println("Notification service is waiting for payment events...")

	for msg := range msgs {
		var event domain.PaymentEvent

		if err := json.Unmarshal(msg.Body, &event); err != nil {
			log.Println("failed to parse message:", err)
			msg.Nack(false, false)
			continue
		}

		if err := c.uc.SendNotification(event); err != nil {
			log.Println("failed to send notification:", err)
			msg.Nack(false, true)
			continue
		}

		msg.Ack(false)
	}

	return nil
}

func (c *RabbitMQConsumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
