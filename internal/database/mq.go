package database

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

var MQConn *amqp.Connection
var MQChannel *amqp.Channel

const OrderQueueName = "surge_orders_queue"

func InitializeMessageQueue() {
	var err error

	MQConn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal("Failed to connect to RabbitMQ : ", err)
	}

	MQChannel, err = MQConn.Channel()
	if err != nil {
		log.Fatal("Failed to open a MQ Channel : ", err)
	}

	/**
	 * parameters: (name, durable, autoDelete, exclusive, noWait, args)
	 */
	_, err = MQChannel.QueueDeclare(OrderQueueName, true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare the MQ queue : ", err)
	}

	fmt.Println("RabbitMQ Connected Successfully.")
}
