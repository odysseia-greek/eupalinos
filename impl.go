package eupalinos

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitClient interface {
	Queue() Queue
}

type Queue interface {
	Send(body []byte) error
	ReceiveOne() (*amqp.Delivery, error)
}

type Rabbit struct {
	queue *QueueImpl
}

func New(connection, queueName string) (*Rabbit, error) {
	conn, err := amqp.Dial(connection)

	queue, err := NewQueueImpl(conn, queueName)
	if err != nil {
		return nil, err
	}

	return &Rabbit{
		queue: queue,
	}, nil

}

func (r *Rabbit) Queue() Queue {
	if r == nil {
		return nil
	}
	return r.queue
}
