package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"shared"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/message"
	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/stream"
)

func main() {
	// Tuning the parameters to test the reliability
	const maxProducersPerClient = 1
	const maxConsumersPerClient = 1

	// Create a RabbitMQ Stream environment
	addresses := []string{os.Getenv("RABBITMQ_STREAM_URL")}
	env, err := stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			SetMaxProducersPerClient(maxProducersPerClient).
			SetMaxConsumersPerClient(maxConsumersPerClient).
			SetUris(addresses))
	shared.FailOnError(err, "Failed to create environment")

	// Declare a stream
	streamName := "basic-stream-filter"
	err = env.DeclareStream(streamName,
		&stream.StreamOptions{
			MaxLengthBytes: stream.ByteCapacity{}.GB(2),
		},
	)
	shared.FailOnError(err, "Failed to declare stream")

	// Create a producer with filter by "type" header
	// This filter will only allow messages with a specific "type" header to be sent
	producer, err := env.NewProducer(streamName, stream.NewProducerOptions().SetFilter(
		stream.NewProducerFilter(func(message message.StreamMessage) string {
			return fmt.Sprintf("%s", message.GetApplicationProperties()["type"])
		})))
	shared.FailOnError(err, "Failed to create producer")

	// Publish messages with confirmation
	messageCount := 1
	chPublishConfirm := producer.NotifyPublishConfirmation()
	ch := make(chan bool)
	handlePublishConfirm(chPublishConfirm, messageCount, ch)

	fmt.Printf("Publishing %d messages...\n", messageCount)
	for i := 0; i < messageCount; i++ {
		message := shared.Message{
			ID:      fmt.Sprintf("%d", time.Now().UTC().Unix()),
			Content: faker.Sentence(),
		}
		messageJson, err := json.Marshal(message)
		shared.FailOnError(err, "Failed to marshal message")

		// Create a new message with application properties
		msg := amqp.NewMessage(messageJson)
		msg.ApplicationProperties = map[string]interface{}{
			"type": randomType(), // Randomly assign a type to the message
		}
		// The filter will only allow messages with a specific "type" header to be sent
		shared.FailOnError(err, "Failed to set message properties")

		// Send the message
		err = producer.Send(msg)
		fmt.Printf("Sent message: %s with type: %s\n", message.ID, msg.ApplicationProperties["type"])
		shared.FailOnError(err, "Failed to send message")
	}
	_ = <-ch
	fmt.Println("Messages confirmed.")

	err = producer.Close()
	shared.FailOnError(err, "Failed to close producer")
}

func randomType() string {
	types := []string{"type1", "type2", "type3"}
	return types[rand.Intn(len(types))]
}

func handlePublishConfirm(confirms stream.ChannelPublishConfirm, messageCount int, ch chan bool) {
	go func() {
		confirmedCount := 0
		for confirmed := range confirms {
			for _, msg := range confirmed {
				if msg.IsConfirmed() {
					confirmedCount++
					if confirmedCount == messageCount {
						ch <- true
					}
				}
			}
		}
	}()
}
