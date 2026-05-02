package publisher

import (
	"encoding/json"

	"github.com/rabbitmq/amqp091-go"
)

type PaymentEvent struct {
	EventID       string `json:"event_id"`
	OrderID       string `json:"order_id"`
	Amount        int64  `json:"amount"`
	CustomerEmail string `json:"customer_email"`
	Status        string `json:"status"`
}

type RabbitMQPublisher struct {
	conn    *amqp091.Connection
	channel *amqp091.Channel
}

func NewRabbitMQPublisher(rabbitURL string) (*RabbitMQPublisher, error) {
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

	return &RabbitMQPublisher{
		conn:    conn,
		channel: ch,
	}, nil
}

func (p *RabbitMQPublisher) Publish(event PaymentEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"",
		"payment.completed",
		false,
		false,
		amqp091.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (p *RabbitMQPublisher) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}
