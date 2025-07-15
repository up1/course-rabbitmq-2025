package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"shared"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	shared.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	// Setup exchanges and queues for DLX and retry mechanism
	err = setup(ch)
	shared.FailOnError(err, "Failed to setup exchanges and queues")

	// Publish messages to the main exchange
	publish(ch)
}

// setup declares the main exchange and queues with Dead Letter Exchange (DLX) and retry mechanisms
func setup(ch *amqp.Channel) error {
	// 1) Final Dead-Letter Exchange & Queue
	if err := ch.ExchangeDeclare("dlx-exchange", "direct", true, false, false, false, nil); err != nil {
		return err
	}
	if _, err := ch.QueueDeclare("dlq", true, false, false, false, nil); err != nil {
		return err
	}
	if err := ch.QueueBind("dlq", "dlq-key", "dlx-exchange", false, nil); err != nil {
		return err
	}

	// 2) Retry Exchange & Queue (with TTL, then DLX back to main-exchange)
	retryArgs := amqp.Table{
		"x-dead-letter-exchange":    "main-exchange",
		"x-message-ttl":             int32(5_000), // 5 seconds TTL
		"x-dead-letter-routing-key": "main-key",
	}
	if err := ch.ExchangeDeclare("retry-exchange", "direct", true, false, false, false, nil); err != nil {
		return err
	}
	if _, err := ch.QueueDeclare("retry-queue", true, false, false, false, retryArgs); err != nil {
		return err
	}
	if err := ch.QueueBind("retry-queue", "retry-key", "retry-exchange", false, nil); err != nil {
		return err
	}

	// 3) Main Exchange & Queue (with DLX to dlx-exchange)
	mainArgs := amqp.Table{
		"x-dead-letter-exchange": "dlx-exchange",
	}
	if err := ch.ExchangeDeclare("main-exchange", "direct", true, false, false, false, nil); err != nil {
		return err
	}
	if _, err := ch.QueueDeclare("main-queue", true, false, false, false, mainArgs); err != nil {
		return err
	}
	if err := ch.QueueBind("main-queue", "main-key", "main-exchange", false, nil); err != nil {
		return err
	}

	return nil
}

func publish(ch *amqp.Channel) error {
	// Enable publisher confirms
	if err := ch.Confirm(false); err != nil {
		return fmt.Errorf("could not enable confirm mode: %w", err)
	}
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 1))

	// Listen for returned (unroutable) messages
	returns := ch.NotifyReturn(make(chan amqp.Return, 1))

	// Create message
	message := generateMessage()
	messageJson, err := json.Marshal(message)
	shared.FailOnError(err, "Failed to marshal message")

	// Publish with mandatory to catch unroutable
	err = ch.PublishWithContext(
		context.Background(),
		"main-exchange", // exchange
		"main-key",      // routing key
		true,            // mandatory
		false,           // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         messageJson,
			Headers:      amqp.Table{"message-id": uuid.NewString()},
		},
	)
	shared.FailOnError(err, "Failed to publish a message")
	log.Printf(" [x] Sent %s", message.String())

	select {
	case ret := <-returns:
		// Unroutable
		log.Printf("📪 Returned message: %s -> %s", ret.Exchange, ret.RoutingKey)
		// Retry by publishing to retry-exchange
		retryMessage := amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         messageJson,
			Headers:      amqp.Table{"message-id": uuid.NewString()},
		}
		if err := ch.PublishWithContext(
			context.Background(),
			"retry-exchange", // retry exchange
			"retry-key",      // routing key
			false,            // mandatory
			false,            // immediate
			retryMessage,
		); err != nil {
			return fmt.Errorf("failed to publish to retry-exchange: %w", err)
		}
		log.Printf("🔄 Retried message to retry-exchange")
	case conf := <-confirms:
		if !conf.Ack {
			log.Printf("❌ Message nacked by broker, retrying or DLQ")
		} else {
			log.Printf("✅ Message confirmed delivery")
		}
	case <-time.After(5 * time.Second):
		return errors.New("timeout waiting for confirm/return")
	}

	return nil
}

func generateMessage() shared.Message {
	return shared.Message{
		ID:      fmt.Sprintf("%d", time.Now().UTC().Unix()),
		Content: faker.Sentence(),
	}
}
