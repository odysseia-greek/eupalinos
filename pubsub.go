// Command pubsub is an example of a fanout exchange with dynamic reliable
// membership, reading from stdin, writing to stdout.
//
// This example shows how to implement reconnect logic independent from a
// publish/subscribe loop with bridges to application types.

package eupalinos

import (
	"context"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type PubSubImpl struct {
	url       string
	queueName string
	exchange  string
}

func NewPubSubImpl(url, queueName, exchange string) (*PubSubImpl, error) {
	return &PubSubImpl{url: url, queueName: queueName, exchange: exchange}, nil
}

type Message []byte

// session composes an amqp.Connection with an amqp.Channel
type session struct {
	*amqp.Connection
	*amqp.Channel
}

// Close tears the connection down, taking the channel with it.
func (s session) Close() error {
	if s.Connection == nil {
		return nil
	}
	return s.Connection.Close()
}

// Redial continually connects to the URL, exiting the program when no longer possible
func (p *PubSubImpl) Redial(ctx context.Context) chan chan session {
	sessions := make(chan chan session)

	go func() {
		sess := make(chan session)
		defer close(sessions)

		for {
			select {
			case sessions <- sess:
			case <-ctx.Done():
				log.Println("shutting down session factory")
				return
			}

			conn, err := amqp.Dial(p.url)
			if err != nil {
				log.Fatalf("cannot (re)dial: %v: %q", err, p.url)
			}

			ch, err := conn.Channel()
			if err != nil {
				log.Fatalf("cannot create channel: %v", err)
			}

			if err := ch.ExchangeDeclare(p.exchange, "fanout", false, true, false, false, nil); err != nil {
				log.Fatalf("cannot declare fanout exchange: %v", err)
			}

			select {
			case sess <- session{conn, ch}:
			case <-ctx.Done():
				log.Println("shutting down new session")
				return
			}
		}
	}()

	return sessions
}

// Publish publishes messages to a reconnecting session to a fanout exchange.
// It receives from the application specific source of messages.
func (p *PubSubImpl) Publish(sessions chan chan session, messages <-chan Message) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	for session := range sessions {
		var (
			running bool
			reading = messages
			pending = make(chan Message, 1)
			confirm = make(chan amqp.Confirmation, 1)
		)

		pub := <-session

		// publisher confirms for this channel/connection
		if err := pub.Confirm(false); err != nil {
			log.Printf("publisher confirms not supported")
			close(confirm) // confirms not supported, simulate by always nacking
		} else {
			pub.NotifyPublish(confirm)
		}

		log.Printf("publishing...")

	Publish:
		for {
			var body Message
			select {
			case confirmed, ok := <-confirm:
				if !ok {
					break Publish
				}
				if !confirmed.Ack {
					log.Printf("nack message %d, body: %q", confirmed.DeliveryTag, string(body))
				}
				reading = messages

			case body = <-pending:
				routingKey := "ignored for fanout exchanges, application dependent for other exchanges"
				err := pub.PublishWithContext(ctx, p.exchange, routingKey, false, false, amqp.Publishing{
					Body: body,
				})
				// Retry failed delivery on the next session
				if err != nil {
					pending <- body
					pub.Close()
					break Publish
				}

			case body, running = <-reading:
				// all messages consumed
				if !running {
					return
				}
				// work on pending delivery until ack'd
				pending <- body
				reading = nil
			}
		}
	}
}

func (p *PubSubImpl) Subscribe(sessions chan chan session, messages chan<- Message) {
	for session := range sessions {
		sub := <-session

		if _, err := sub.QueueDeclare(p.queueName, false, true, true, false, nil); err != nil {
			log.Printf("cannot consume from exclusive queue: %q, %v", p.queueName, err)
			return
		}

		routingKey := "application specific routing key for fancy topologies"
		if err := sub.QueueBind(p.queueName, routingKey, p.exchange, false, nil); err != nil {
			log.Printf("cannot consume without a binding to exchange: %q, %v", p.exchange, err)
			return
		}

		deliveries, err := sub.Consume(p.queueName, "", false, true, false, false, nil)
		if err != nil {
			log.Printf("cannot consume from: %q, %v", p.queueName, err)
			return
		}

		log.Printf("subscribed...")

		for msg := range deliveries {
			messages <- msg.Body
			sub.Ack(msg.DeliveryTag, false)
		}
	}
}
