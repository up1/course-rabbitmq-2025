package main

import (
	"encoding/json"
	"log"
	"os"
	"shared"

	amqp "github.com/rabbitmq/amqp091-go"
)

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
		"dlq", // queue
		"",    // consumer
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	shared.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			log.Printf("Dead Letter message, headers=%v", d.Headers)

			// Process the message
			message := shared.Message{}
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Printf("Failed to unmarshal message: %s", err)
				d.Nack(false, false)
				continue
			}

			log.Printf("✅ processed successfully: %s", message.String())
			// Acknowledge successful processing
			d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
