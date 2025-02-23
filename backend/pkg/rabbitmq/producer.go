package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func PublishMessage(msgType string, ch *amqp.Channel, body []byte) error {
	err := ch.Publish(
		"",           // Exchange (default)
		"task_queue", // Routing key (queue name)
		false,        // Mandatory
		false,        // Immediate
		amqp.Publishing{
			Type:         msgType,
			DeliveryMode: amqp.Persistent, // Ensure message persistence
			ContentType:  "application/json",
			Body:         body,
		})
	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	fmt.Println("✅ Message published:", body)
	return nil
}

func CheckQueue(ch *amqp.Channel) {
	// Inspect the queue
	queue, err := ch.QueueInspect("task_queue")
	if err != nil {
		log.Fatalf("Failed to inspect queue: %v", err)
	}

	fmt.Printf("Queue Name: %s\n", queue.Name)
	fmt.Printf("Messages in queue: %d\n", queue.Messages)
	fmt.Printf("Consumers in queue: %d\n", queue.Consumers)

}
