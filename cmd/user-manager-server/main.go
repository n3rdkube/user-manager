package main

import (
	"net/http"

	"github.com/n3rdkube/user-manager/internal/api"
	ms "github.com/n3rdkube/user-manager/internal/messanger"
	"github.com/sirupsen/logrus"
)

// This is the main of the service in charge of the REST api
// All calls are completely async but the list users.

const (
	rabbitMQURL      = "amqp://guest:guest@rabbitmq:5672/"
	listeningAddress = "0.0.0.0:35307"
)

func main() {

	logrus.Info("creating connection with the message broker")
	conn, err := ms.OpenConnection(rabbitMQURL)
	if err != nil {
		logrus.Fatalf("connecting to messaging system: %v", err)
	}
	defer conn.Close()

	logrus.Info("starting consumer to read callback answers")
	consumer, err := ms.NewRabbitMQ(conn, ms.CallbackQueue)
	if err != nil {
		logrus.Fatalf("creating rabbit mq instance: %v", err)
	}
	defer consumer.Close()

	logrus.Info("starting processor producers to process request async")
	producer, err := ms.NewRabbitMQ(conn, ms.ProcessorQueue)
	if err != nil {
		logrus.Fatalf("creating rabbit mq instance for callbacks: %v", err)
	}
	defer producer.Close()

	logrus.Info("creating mux that will handle all requests")
	mux := api.NewServerMux(consumer, producer, api.LogHandler)

	logrus.Infof("staring server")
	err = http.ListenAndServe(listeningAddress, mux)
	if err != nil {
		logrus.Errorf("starting serving on localhost: %v", err)
	}
}
