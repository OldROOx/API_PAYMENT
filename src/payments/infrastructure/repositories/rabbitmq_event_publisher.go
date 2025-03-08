package repositories

import (
	"app.payment/src/payments/domain/entities"
	"encoding/json"
	"github.com/streadway/amqp"
)

type RabbitMQEventPublisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
}

func NewRabbitMQEventPublisher(amqpURL, exchange string) (*RabbitMQEventPublisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// Declare the exchange
	err = ch.ExchangeDeclare(
		exchange, // name
		"topic",  // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQEventPublisher{
		conn:     conn,
		channel:  ch,
		exchange: exchange,
	}, nil
}

func (p *RabbitMQEventPublisher) PublishEvent(event entities.Event) error {
	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		p.exchange, // exchange
		event.Type, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
}

func (p *RabbitMQEventPublisher) Close() error {
	if err := p.channel.Close(); err != nil {
		return err
	}
	return p.conn.Close()
}
