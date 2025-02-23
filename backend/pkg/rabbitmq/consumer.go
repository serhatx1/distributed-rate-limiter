package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"main/internals/repository"
	"main/pkg/models"
	"time"

	"github.com/streadway/amqp"
	"go.mongodb.org/mongo-driver/bson"
)

func Consumer(ch *amqp.Channel) {
	msgs, err := ch.Consume(
		"task_queue", // Queue name
		"",           // Consumer tag (empty means RabbitMQ assigns one)
		false,        // Auto-acknowledge (false ensures manual ACK)
		false,        // Exclusive
		false,        // No-local
		false,        // No-wait
		nil,          // Arguments
	)
	if err != nil {
		log.Fatalf("❌ Failed to start consuming: %v", err)
	}

	go func() {
		for msg := range msgs {
			fmt.Printf("✅ Received a message: %s and type is %s\n", msg.Body, msg.Type)

			if msg.Type == "setlimit" {
				var changeLimitReq models.ChangeLimitRequest
				err := json.Unmarshal(msg.Body, &changeLimitReq)
				if err != nil {
					log.Printf("❌ Error unmarshalling JSON: %v\n", err)
					_ = msg.Nack(false, false) // Reject message without requeueing
					continue
				}

				err = SetChange(&changeLimitReq)
				if err != nil {
					log.Printf("❌ Error processing request: %v\n", err)
					_ = msg.Nack(false, false) // Reject and requeue the message for retry
					continue
				}

				_ = msg.Ack(false)
			}
		}

	}()

}

func SetChange(changeLimitReq *models.ChangeLimitRequest) error {
	filter := bson.M{"endpoint": bson.M{"$eq": changeLimitReq.EndPoint}}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()
	collection := repository.GetRateLimiterCollection()
	if err := collection.FindOne(ctx, filter).Err(); err != nil {
		log.Printf("There is no endpoint called that.", err)
		return err
	}
	update := bson.M{
		"$set": bson.M{
			"ratelimit": changeLimitReq.Ratelimit,
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("wrong typed rate limit: %v", err)
		return err

	}
	fmt.Print("success")
	return nil
}
