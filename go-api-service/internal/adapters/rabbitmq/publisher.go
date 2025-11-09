package rabbitmq

import (
	"context"
	"time"

	"github.com/AchilleasB/southpark-messaging/go-api-service/internal/ports"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue string
}

func NewPublisher(amqpURL, queue string) (ports.Publisher, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}
	_, err = ch.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}
	return &Publisher{conn: conn, ch: ch, queue: queue}, nil
}

func (p *Publisher) Publish(ctx context.Context, topic string, payload []byte) error {
	// Use topic param for future flexibility; we keep a single queue
	return p.ch.PublishWithContext(ctx,
		"",      // exchange
		p.queue, // routing key = queue
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         payload,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)
}

func (p *Publisher) Close() error {
	_ = p.ch.Close()
	return p.conn.Close()
}
