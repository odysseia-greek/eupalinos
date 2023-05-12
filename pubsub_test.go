package eupalinos

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPubSub(t *testing.T) {
	name := "test"
	exchange := "pubsubtest"
	//url := "ampq://testrabbit@test"

	t.Run("Consumer", func(t *testing.T) {
		ctx, done := context.WithCancel(context.Background())
		testClient, err := NewMockRabbit(name, exchange)
		assert.Nil(t, err)
		assert.NotNil(t, testClient)

		lines := make(chan Message)

		go func() {
			testClient.PubSub().Publish(testClient.PubSub().Redial(ctx), lines)
			done()
		}()

		messages := []string{"hello", "world", "again", "sdfgsdf", "sadfsdf", "sdfgsg"}
		for _, msg := range messages {
			time.Sleep(time.Second * 1000)
			lines <- Message(msg)
		}
	})
}
