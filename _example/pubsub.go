package _example

import (
	"context"
	"github.com/odysseia-greek/eupalinos"
	"log"
	"time"
)

func Example() {
	ctx, done := context.WithCancel(context.Background())
	lines := make(chan eupalinos.Message)

	channelName := "thisisthechannelname"
	exchangeName := "thisistheexchangename"

	p, err := eupalinos.New(channelName, exchangeName)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		p.PubSub().Publish(p.PubSub().Redial(ctx), lines)
		done()
	}()

	received := make(chan eupalinos.Message)
	go func() {
		p.PubSub().Subscribe(p.PubSub().Redial(ctx), received)
		done()
	}()

	go func() {
		for line := range received {
			log.Print(line)
		}
	}()

	messages := []string{"hello", "world", "again", "sdfgsdf", "sadfsdf", "sdfgsg"}
	for _, msg := range messages {
		time.Sleep(time.Second * 1000)
		lines <- eupalinos.Message(msg)
	}

	<-ctx.Done()
}
