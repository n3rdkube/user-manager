package messanger

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

const (
	// ProcessorQueue name
	ProcessorQueue = "processorQueue"
	// CallbackQueue name
	CallbackQueue = "callbackQueue"
	// NotificationQueue name
	NotificationQueue = "notificationQueue"

	// An empty string will cause
	//the library to generate a unique identity
	consumerName = ""

	defaultDurable    = false
	defaultAutodelete = false
	defaultExclusive  = false
	defaultNoWait     = false
	defaultAutoAck    = false
	defaultNoLocal    = false
	defaultMandatory  = false
	defaultImmediate  = false
)

// RabbitMQ stub implementation of a messaging system, if needed we could have different structs for producers and consumers
// having different interfaces was considered enough for a stub
type RabbitMQ struct {
	conn  *amqp.Connection
	queue string
}

// NewRabbitMQ setup new RabbitMQ opening connections and creating the queue if needed
func NewRabbitMQ(conn *amqp.Connection, queueName string) (*RabbitMQ, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("creating channel to rabbitMQ: %w", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(queueName, defaultDurable, defaultAutodelete, defaultExclusive, defaultNoWait, nil)
	if err != nil {
		return nil, fmt.Errorf("defining %q: %w", queueName, err)
	}

	return &RabbitMQ{
		conn:  conn,
		queue: queueName,
	}, nil
}

// OpenConnection opens a rabbitmq connection
func OpenConnection(url string) (*amqp.Connection, error) {
	conn, err := amqp.Dial(url)

	if err != nil {
		return nil, fmt.Errorf("connecting to rabbitMQ: %w", err)
	}

	return conn, nil
}

// Close closes a rabbitmq connection
func (mq RabbitMQ) Close() error {
	err := mq.conn.Close()
	if err != nil {
		return fmt.Errorf("closing channel")
	}

	return nil
}

// GetMessageFromQueue get one message from the queue
// it could have been implemented returning a channel
func (mq RabbitMQ) GetMessageFromQueue() (*Message, error) {
	ch, err := mq.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("creating channel to rabbitMQ: %w", err)
	}
	defer ch.Close()

	qCh, err := ch.Consume(mq.queue, consumerName, defaultAutoAck, defaultExclusive, defaultNoLocal, defaultNoWait, nil)
	if err != nil {
		return nil, fmt.Errorf("consuming %q: %w", mq.queue, err)
	}

	logrus.Infof("reading message from %q", mq.queue)
	ms := <-qCh
	err = ms.Ack(false)
	if err != nil {
		return nil, fmt.Errorf("failed to ack message %w", err)
	}

	return &Message{
		Data:        ms.Body,
		MessageType: ms.Type,
		CID:         ms.CorrelationId,
	}, nil
}

// GetMessageFromQueueWithCID get a message from the queue with a specific CID requeueing if needed
func (mq RabbitMQ) GetMessageFromQueueWithCID(cID string) (*Message, error) {
	ch, err := mq.conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("creating channel to rabbitMQ: %w", err)
	}
	defer ch.Close()

	qCh, err := ch.Consume(mq.queue, consumerName, defaultAutoAck, defaultExclusive, defaultNoLocal, defaultNoWait, nil)
	if err != nil {
		return nil, fmt.Errorf("consuming %q: %w", mq.queue, err)
	}

	logrus.Infof("reading message from %q with cid %q", mq.queue, cID)
	for {
		select {
		case ms := <-qCh:
			if ms.CorrelationId == cID {
				err = ms.Ack(false)
				if err != nil {
					return nil, fmt.Errorf("failed to ack message %w", err)
				}
				return &Message{
					Data:        ms.Body,
					MessageType: ms.Type,
					CID:         ms.CorrelationId,
				}, nil
			}
			err = ms.Nack(false, true)
			if err != nil {
				return nil, fmt.Errorf("failed to nack message %w", err)
			}
			logrus.Infof("message with a different CID was received: %q !=  %q", ms.CorrelationId, cID)

		case <-time.After(time.Minute):
			return nil, fmt.Errorf("deadline exceedeed waiting for message %q", cID)
		}
	}

}

// PostMessageToQueue post a message to a queue
func (mq RabbitMQ) PostMessageToQueue(message Message) error {
	ch, err := mq.conn.Channel()
	if err != nil {
		return fmt.Errorf("creating channel to rabbitMQ: %w", err)
	}
	defer ch.Close()

	logrus.Infof("posting message to %q with cid %q", mq.queue, message.CID)
	err = ch.Publish("", mq.queue, defaultMandatory, defaultImmediate,
		amqp.Publishing{
			ContentType:   "application/json",
			Body:          message.Data,
			Type:          message.MessageType,
			CorrelationId: message.CID,
		},
	)
	if err != nil {
		return fmt.Errorf("publishing data: %w", err)
	}

	return nil
}
