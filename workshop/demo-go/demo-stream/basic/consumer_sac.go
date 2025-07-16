package main

import (
	"bufio"
	"fmt"
	"os"
	"shared"
	"time"

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
	streamName := "single-active-consumer-stream"
	err = env.DeclareStream(streamName,
		&stream.StreamOptions{
			MaxLengthBytes: stream.ByteCapacity{}.GB(2),
		},
	)
	shared.FailOnError(err, "Failed to declare stream")

	// Create a consumer
	messagesHandler := func(consumerContext stream.ConsumerContext, message *amqp.Message) {
		// Processing the message
		fmt.Printf("Stream: %s - Received message: %s\n", consumerContext.Consumer.GetStreamName(), message.Data)

		// Store the offset after processing the message :: bad practice
		err := consumerContext.Consumer.StoreOffset()
		shared.FailOnError(err, "Failed to store offset")
	}
	// Create a consumer with last offset
	consumerName := "consumer-" + os.Args[1]

	// Consumer update function
	consumerUpdate := func(streamName string, isActive bool) stream.OffsetSpecification {
		fmt.Printf("[%s] - Consumer promoted for: %s. Active status: %t\n", time.Now().Format(time.TimeOnly),
			streamName, isActive)

		offset, err := env.QueryOffset(consumerName, streamName)
		if err != nil {
			// If the offset is not found, we start from the beginning
			return stream.OffsetSpecification{}.First()
		}

		// If the offset is found, we start from the last offset
		// we add 1 to the offset to start from the next message
		return stream.OffsetSpecification{}.Offset(offset + 1)
	}

	// Create a consumer with offset
	consumer, err := env.NewConsumer(streamName, messagesHandler,
		stream.NewConsumerOptions().
			SetConsumerName(consumerName).
			SetSingleActiveConsumer(stream.NewSingleActiveConsumer(consumerUpdate)))
	shared.FailOnError(err, "Failed to create consumer")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println(" [x] Waiting for messages. enter to close the consumer")
	_, _ = reader.ReadString('\n')
	err = consumer.Close()
	shared.FailOnError(err, "Failed to close consumer")
}
