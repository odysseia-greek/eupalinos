package eupalinos

import "context"

type RabbitClient interface {
	PubSub() PubSub
}

type PubSub interface {
	Subscribe(sessions chan chan session, messages chan<- Message)
	Publish(sessions chan chan session, messages <-chan Message)
	Redial(ctx context.Context) chan chan session
}

type Rabbit struct {
	pubsub *PubSubImpl
}

type FakeRabbit struct {
	pubsub *PubSubMockImpl
}

func New(queueName, exchange string) (*Rabbit, error) {
	connection := CreateQueueConnectionStringFromEnv()
	pubsub, err := NewPubSubImpl(connection, queueName, exchange)
	if err != nil {
		return nil, err
	}

	return &Rabbit{
		pubsub: pubsub,
	}, nil

}

func NewMockRabbit(queueName, exchange string) (*FakeRabbit, error) {
	connection := CreateQueueConnectionStringFromEnv()
	pubsub, err := NewPubSubMockImpl(connection, queueName, exchange)
	if err != nil {
		return nil, err
	}

	return &FakeRabbit{
		pubsub: pubsub,
	}, nil

}

func (r *Rabbit) PubSub() PubSub {
	if r == nil {
		return nil
	}
	return r.pubsub
}

func (r *FakeRabbit) PubSub() PubSub {
	if r == nil {
		return nil
	}
	return r.pubsub
}
