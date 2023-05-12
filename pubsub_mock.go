package eupalinos

import (
	"context"
	"log"
)

type PubSubMockImpl struct {
	url       string
	queueName string
	exchange  string
}

func NewPubSubMockImpl(url, queueName, exchange string) (*PubSubMockImpl, error) {
	return &PubSubMockImpl{url: url, queueName: queueName, exchange: exchange}, nil
}

func (p *PubSubMockImpl) Redial(ctx context.Context) chan chan session {
	return nil
}

func (p *PubSubMockImpl) Publish(sessions chan chan session, messages <-chan Message) {
	log.Print(messages)
}

func (p *PubSubMockImpl) Subscribe(sessions chan chan session, messages chan<- Message) {

}
