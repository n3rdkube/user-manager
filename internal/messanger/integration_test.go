//go:build integration

package messanger_test

import (
	"testing"
	"time"

	"github.com/n3rdkube/user-manager/internal/messanger"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"

	"github.com/avast/retry-go"
	"github.com/streadway/amqp"
)

const (
	rabbitMQURL     = "amqp://guest:guest@:5672/"
	testQueue       = "test"
	testData        = "testData"
	testCid         = "testCid"
	testMessageType = "testMessageType"
)

func Test_Messages_Are_Delivered(t *testing.T) {
	conn := connect(t)
	producer, err := messanger.NewRabbitMQ(conn, testQueue)
	require.NoError(t, err)
	defer conn.Close()

	consumer, err := messanger.NewRabbitMQ(conn, testQueue)
	require.NoError(t, err)

	m := messanger.Message{
		Data:        []byte(testData),
		MessageType: testMessageType,
		CID:         testCid,
	}

	err = producer.PostMessageToQueue(m)
	require.NoError(t, err)

	mess, err := consumer.GetMessageFromQueue()
	require.NoError(t, err)
	assert.Equal(t, m, *mess)
}

func Test_Messages_Are_Delivered_with_cid(t *testing.T) {

	conn := connect(t)
	producer, err := messanger.NewRabbitMQ(conn, testQueue)
	require.NoError(t, err)
	defer conn.Close()

	consumer, err := messanger.NewRabbitMQ(conn, testQueue)
	require.NoError(t, err)

	m := messanger.Message{
		Data:        []byte(testData),
		MessageType: testMessageType,
		CID:         testCid,
	}

	mNoise := messanger.Message{
		Data:        []byte(testData),
		MessageType: testMessageType,
		CID:         "noise",
	}

	err = producer.PostMessageToQueue(mNoise)
	require.NoError(t, err)
	err = producer.PostMessageToQueue(mNoise)
	require.NoError(t, err)
	err = producer.PostMessageToQueue(mNoise)
	require.NoError(t, err)
	err = producer.PostMessageToQueue(m)
	require.NoError(t, err)
	err = producer.PostMessageToQueue(mNoise)
	require.NoError(t, err)
	err = producer.PostMessageToQueue(mNoise)
	require.NoError(t, err)
	err = producer.PostMessageToQueue(mNoise)
	require.NoError(t, err)

	mess, err := consumer.GetMessageFromQueueWithCID(testCid)
	require.NoError(t, err)
	assert.Equal(t, m, *mess)

	t.Run("and_are_requeued", func(t *testing.T) {
		for i := 0; i < 6; i++ {
			mess, err = consumer.GetMessageFromQueue()
			require.NoError(t, err)
			assert.Equal(t, mNoise, *mess)
		}
	})
}

func connect(t *testing.T) *amqp.Connection {
	var conn *amqp.Connection
	var err error

	err = retry.Do(
		func() error {
			conn, err = messanger.OpenConnection(rabbitMQURL)
			return err
		},
		retry.Delay(time.Second*2),
	)
	require.NoError(t, err)

	return conn
}
