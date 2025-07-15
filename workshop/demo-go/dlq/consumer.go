package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"shared"

	amqp "github.com/rabbitmq/amqp091-go"
)

const maxRetries = 3

func main() {
	conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
	shared.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	shared.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	shared.FailOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		"main-queue", // queue
		"",           // consumer
		false,        // auto-ack
		false,        // exclusive
		false,        // no-local
		false,        // no-wait
		nil,          // args
	)
	shared.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			// Detect redelivery
			if d.Redelivered {
				log.Printf("↩️ redelivered, headers=%v", d.Headers)
			}

			// Extract retry count from header
			retries := 0
			if v, ok := d.Headers["x-retries"].(int32); ok {
				retries = int(v)
			}

			// Process the message
			message := shared.Message{}
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Printf("Failed to unmarshal message: %s", err)
				d.Nack(false, false)
				continue
			}

			// Process the message
			if err := processMessage(message); err != nil {
				retries++

				if retries < maxRetries {
					log.Printf("⚠️ processing failed (attempt %d), retrying... with message: %s", retries, message.String())

					// Create new headers with updated retry count
					newHeaders := make(amqp.Table)
					for k, v := range d.Headers {
						newHeaders[k] = v
					}
					newHeaders["x-retries"] = int32(retries)

					// Republish message with updated headers for retry
					err = ch.Publish(
						"retry-exchange", // exchange
						"retry-key",      // routing key
						false,            // mandatory
						false,            // immediate
						amqp.Publishing{
							ContentType: "application/json",
							Body:        d.Body,
							Headers:     newHeaders,
						},
					)
					if err != nil {
						log.Printf("Failed to republish message for retry: %s", err)
					}

					// Acknowledge the original message since we've republished it
					d.Ack(false)
				} else {
					log.Printf("❌ exceeded retries (%d), sending to DLQ with message: %s", retries, message.String())
					// Nack → DLX on main-queue routes to dlx-exchange → dlq
					err = ch.Publish(
						"dlx-exchange", // exchange
						"dlq-key",      // routing key
						false,          // mandatory
						false,          // immediate
						amqp.Publishing{
							ContentType: "application/json",
							Body:        d.Body,
							Headers:     amqp.Table{"x-retries": int32(retries)},
						},
					)
					if err != nil {
						log.Printf("Failed to publish to DLQ: %s", err)
					}
					// Acknowledge the original message
					d.Ack(false)
				}
			} else {
				log.Printf("✅ processed successfully: %s", message.String())
				// Acknowledge successful processing
				d.Ack(false)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func processMessage(message shared.Message) error {
	// Simulate processing logic
	if rand.Intn(10) < 5 { // 90% chance of failure
		return fmt.Errorf("simulated processing error")
	}
	log.Printf("Processing message: %s", message.String())
	return nil
}
