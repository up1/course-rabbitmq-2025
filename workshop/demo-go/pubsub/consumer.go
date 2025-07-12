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

	q, err := ch.QueueDeclare(
		"",    // name with empty string to create a unique name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	shared.FailOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,        // queue name
		"",            // routing key
		"logs_fanout", // exchange name
		false,
		nil)
	shared.FailOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	shared.FailOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			message := shared.Message{}
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Printf("Failed to unmarshal message: %s", err)
				continue
			}
			log.Printf("Processing message: %s", message.String())
			log.Printf("Done")
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
