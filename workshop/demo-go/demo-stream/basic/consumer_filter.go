package main

import (
	"bufio"
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
	streamName := "basic-stream-filter"
	err = env.DeclareStream(streamName,
		&stream.StreamOptions{
			MaxLengthBytes: stream.ByteCapacity{}.GB(2),
		},
	)
	shared.FailOnError(err, "Failed to declare stream")

	// Handle incoming messages
	messagesHandler := func(consumerContext stream.ConsumerContext, message *amqp.Message) {
		fmt.Printf("Stream: %s - Received message: %s with type: %s\n\n", consumerContext.Consumer.GetStreamName(),
			message.Data, message.ApplicationProperties["type"])
	}
	// Create a consumer with filtering
	type_name := os.Args[2]
	typeFilter := func(message *amqp.Message) bool {
		return message.ApplicationProperties["type"] == type_name
	}
	filter := stream.NewConsumerFilter([]string{type_name}, true, typeFilter)

	// Create a consumer with the filter
	name := "consumer-" + os.Args[1]
	consumer, err := env.NewConsumer(streamName, messagesHandler,
		stream.NewConsumerOptions().
			SetConsumerName(name).
			SetOffset(stream.OffsetSpecification{}.First()).
			SetFilter(filter))
	shared.FailOnError(err, "Failed to create consumer")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println(" [x] Waiting for messages. enter to close the consumer")
	_, _ = reader.ReadString('\n')
	err = consumer.Close()
	shared.FailOnError(err, "Failed to close consumer")
}
