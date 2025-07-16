package main

import (
	"bufio"
	"errors"
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

	// Create a consumer
	messagesHandler := func(consumerContext stream.ConsumerContext, message *amqp.Message) {
		// Processing the message
		fmt.Printf("Stream: %s - Received message: %s\n", consumerContext.Consumer.GetStreamName(),
			message.Data)

		// Store the offset after processing the message :: bad practice
		err := consumerContext.Consumer.StoreOffset()
		shared.FailOnError(err, "Failed to store offset")
	}
	// Create a consumer with last offset
	var offsetSpecification stream.OffsetSpecification
	name := "consumer-offset-" + os.Args[1]
	storedOffset, err := env.QueryOffset(name, streamName)
	if errors.Is(err, stream.OffsetNotFoundError) {
		offsetSpecification = stream.OffsetSpecification{}.First() // If no offset is stored, start from the beginning
	} else {
		offsetSpecification = stream.OffsetSpecification{}.Offset(storedOffset + 1)
	}

	// Create a consumer with offset
	consumer, err := env.NewConsumer(streamName, messagesHandler,
		stream.NewConsumerOptions().
			SetConsumerName(name).
			SetManualCommit().
			SetOffset(offsetSpecification))
	shared.FailOnError(err, "Failed to create consumer")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println(" [x] Waiting for messages. enter to close the consumer")
	_, _ = reader.ReadString('\n')
	err = consumer.Close()
	shared.FailOnError(err, "Failed to close consumer")
}
