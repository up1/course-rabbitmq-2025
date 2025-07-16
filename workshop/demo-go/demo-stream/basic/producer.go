package main

import (
	"encoding/json"
	"fmt"
	"os"
	"shared"

	"github.com/rabbitmq/rabbitmq-stream-go-client/pkg/amqp"
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
	streamName := "basic-stream"
	err = env.DeclareStream(streamName,
		&stream.StreamOptions{
			MaxLengthBytes: stream.ByteCapacity{}.GB(2),
		},
	)
	shared.FailOnError(err, "Failed to declare stream")

	// Create a producer
	producer, err := env.NewProducer(streamName, stream.NewProducerOptions())
	shared.FailOnError(err, "Failed to create producer")

	// Publish messages with confirmation
	messageCount := 10
	chPublishConfirm := producer.NotifyPublishConfirmation()
	ch := make(chan bool)
	handlePublishConfirm(chPublishConfirm, messageCount, ch)

	fmt.Printf("Publishing %d messages...\n", messageCount)
	for i := 0; i < messageCount; i++ {
		message := shared.Message{
			ID:      fmt.Sprintf("msg-%d", i),
			Content: fmt.Sprintf("This is message %d", i),
		}
		messageJson, err := json.Marshal(message)
		shared.FailOnError(err, "Failed to marshal message")

		// Send the message
		err = producer.Send(amqp.NewMessage(messageJson))
		shared.FailOnError(err, "Failed to send message")
	}
	_ = <-ch
	fmt.Println("Messages confirmed.")

	err = producer.Close()
	shared.FailOnError(err, "Failed to close producer")
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
