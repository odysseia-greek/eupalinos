package eupalinos

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSendMessage(t *testing.T) {
	name := "test"

	t.Run("Created", func(t *testing.T) {
		connection := CreateQueueConnectionStringFromEnv()
		testClient, err := New(connection, name)
		assert.Nil(t, err)

		sut := []byte("hello")
		err = testClient.Queue().Send(sut)
		assert.Nil(t, err)
		msg, err := testClient.Queue().ReceiveOne()
		assert.Equal(t, sut, msg.Body)
	})
}
