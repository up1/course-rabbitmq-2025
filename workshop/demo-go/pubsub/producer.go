package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"shared"
	"time"

	"github.com/go-faker/faker/v4"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	shared.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_fanout", // name
		"fanout",      // type
		true,          // durable
		false,         // auto-deleted
		false,         // internal
		false,         // no-wait
		nil,           // arguments
	)
	shared.FailOnError(err, "Failed to declare an exchange")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	message := generateMessage()
	messageJson, err := json.Marshal(message)
	shared.FailOnError(err, "Failed to marshal message")
	err = ch.PublishWithContext(ctx,
		"logs_fanout", // exchange
		"",            // routing key
		false,         // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         messageJson,
		})
	shared.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", message.String())
}

func generateMessage() shared.Message {
	return shared.Message{
		ID:      fmt.Sprintf("%d", time.Now().UTC().Unix()),
		Content: faker.Sentence(),
	}
}
