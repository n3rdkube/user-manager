package main

import (
	ms "github.com/n3rdkube/user-manager/internal/messanger"
	mp "github.com/n3rdkube/user-manager/internal/processor"
	"github.com/n3rdkube/user-manager/internal/storage"
	"github.com/sirupsen/logrus"
)

// This is the main of the service in charge of dispatching all messages received in the queue ms.ProcessorQueue.
// Once the message is received it connect to the DB, answer back if needed and notify on a different queue changes.

const (
	rabbitMQURL = "amqp://guest:guest@rabbitmq:5672/"
	mysqlURL    = "root:root@tcp(db:3306)/users3"
)

func main() {
	logrus.Info("creating connection with the message broker")
	conn, err := ms.OpenConnection(rabbitMQURL)
	if err != nil {
		logrus.Fatalf("opening rabbit mq connection: %v", err)
	}
	defer conn.Close()

	logrus.Info("creating consumer with message broker")
	mqConsumer, err := ms.NewRabbitMQ(conn, ms.ProcessorQueue)
	if err != nil {
		logrus.Fatalf("creating rabbit mq instance: %v", err)
	}
	defer mqConsumer.Close()

	logrus.Info("creating a callback producer with message broker")
	mqProducer, err := ms.NewRabbitMQ(conn, ms.CallbackQueue)
	if err != nil {
		logrus.Fatalf("creating rabbit mq instance for callbacks: %v", err)
	}
	defer mqProducer.Close()

	logrus.Info("creating a notification producer with message broker")
	mqNotification, err := ms.NewRabbitMQ(conn, ms.NotificationQueue)
	if err != nil {
		logrus.Fatalf("creating rabbit mq instance for callbacks: %v", err)
	}
	defer mqNotification.Close()

	logrus.Info("creating a mysql storage to save data")
	st, err := storage.NewMysqlStorageDB(mysqlURL)
	if err != nil {
		logrus.Fatal("creating storage instance")
	}

	logrus.Info("start processing messages")
	proc := mp.NewProcessor(st, mqProducer, mqNotification, mqConsumer)
	proc.StartProcessing()
}
