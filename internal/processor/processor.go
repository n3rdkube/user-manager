package processor

import (
	"errors"
	"fmt"
	"time"

	ms "github.com/n3rdkube/user-manager/internal/messanger"
	"github.com/n3rdkube/user-manager/internal/storage"
	"github.com/sirupsen/logrus"
)

// Processor checks requests from manager and executes them
type Processor struct {
	mqProcessor   ms.MessageQueueRead
	mqCallback    ms.MessageQueueWrite
	notifications ms.MessageQueueWrite

	st storage.Storage
}

// NewProcessor initialize a Processor
func NewProcessor(st storage.Storage, mqCallback ms.MessageQueueWrite, mqNotificator ms.MessageQueueWrite, mqProcessor ms.MessageQueueRead) *Processor {
	return &Processor{
		mqProcessor:   mqProcessor,
		mqCallback:    mqCallback,
		notifications: mqNotificator,
		st:            st,
	}
}

// StartProcessing takes one by one all messages from the queue
func (p *Processor) StartProcessing() {
	logrus.Info("Starting processing messages")
	for {
		m, err := p.mqProcessor.GetMessageFromQueue()
		if err != nil {
			logrus.Warnf("error while retrieving message: %v", err)
			time.Sleep(time.Second)
			continue
		}

		err = p.processMessages(m)
		if err != nil {
			logrus.Warnf("error while processing message: %v", err)
			time.Sleep(time.Second)
		}
	}
}

func (p *Processor) processMessages(message *ms.Message) error {
	if message == nil {
		return errors.New("message is nil")
	}

	logrus.Infof("New message received %q", string(message.Data))
	switch message.MessageType {
	case ms.CreateUser:
		if err := p.addUser(message.Data); err != nil {
			return fmt.Errorf("adding user: %w", err)
		}

	case ms.DeleteUser:
		if err := p.deleteUser(message.Data); err != nil {
			return fmt.Errorf("deleting user: %w", err)
		}

	case ms.ListUser:
		if err := p.listUser(message.Data, message.CID); err != nil {
			return fmt.Errorf("listing user: %w", err)
		}

	case ms.UpdateUser:
		if err := p.updateUser(message.Data); err != nil {
			return fmt.Errorf("updating user: %w", err)
		}

	default:
		return fmt.Errorf("message type not recognised: %q, %q", message.MessageType, string(message.Data))
	}

	return nil
}
