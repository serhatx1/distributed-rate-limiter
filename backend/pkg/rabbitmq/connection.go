package rabbitmq

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

var (
	channel *amqp.Channel
	conn    *amqp.Connection
	once    sync.Once
)

func RabbitMQInit() {
	once.Do(func() {
		var err error
		conn, err = ConnectRabbitMQ()
		if err != nil {
			log.Fatalf("Failed to connect to RabbitMQ after multiple attempts: %v", err)
		}
		fmt.Println("✅ Connected to RabbitMQ")

		channel, err = conn.Channel()
		if err != nil {
			log.Fatalf("Failed to open a channel: %v", err)
		}

		DeclareQueue(channel)
		Consumer(channel)

	})
}
func ConnectRabbitMQ() (*amqp.Connection, error) {
	var conn *amqp.Connection
	var err error
	for i := 0; i < 3; i++ {
		conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err == nil {
			return conn, nil
		}
		reTryTime := time.Duration(1<<uint(i+1)) * time.Second
		log.Printf("Failed to connect to RabbitMQ: %v. Retrying in %v seconds...", err, reTryTime)

		time.Sleep(reTryTime)
	}

	return conn, err
}
func DeclareQueue(ch *amqp.Channel) {
	_, err := ch.QueueDeclare(
		"task_queue", // queue name
		true,         // durable (persists after restartss)
		false,        // auto-delete (won't delete automatically)
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	// Set QoS (Prefetch Count for Rate Limiting)
	err = ch.Qos(
		5,     // prefetchCount (Number of messages the worker can process at once)
		0,     // prefetchSize (Not used)
		false, // global (false = applied to this channel only)
	)
	if err != nil {
		log.Fatalf("Failed to set QoS: %v", err)
	}

}
func Channel() *amqp.Channel {
	if channel == nil {
		log.Fatalf("Channel is nil, you need to call RabbitMQInit() first.")
	}
	return channel
}
