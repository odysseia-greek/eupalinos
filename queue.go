package eupalinos

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

type QueueImpl struct {
	channel *amqp.Channel
	queue   *amqp.Queue
}

func NewQueueImpl(conn *amqp.Connection, queueName string) (*QueueImpl, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	q, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return nil, err
	}

	return &QueueImpl{channel: channel, queue: &q}, nil
}

func (q *QueueImpl) Send(body []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return q.channel.PublishWithContext(ctx,
		"",           // exchange
		q.queue.Name, // routing key
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
}

func (q *QueueImpl) ReceiveOne() (*amqp.Delivery, error) {
	message, ok, err := q.channel.Get(
		q.queue.Name, // queue
		true,
	)

	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf("not ok")
	}

	return &message, nil
}
