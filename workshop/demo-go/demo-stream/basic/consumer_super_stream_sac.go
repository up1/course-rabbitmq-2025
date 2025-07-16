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
	const EnableSingleActiveConsumer = true

	// Create a RabbitMQ Stream environment
	addresses := []string{os.Getenv("RABBITMQ_STREAM_URL")}
	env, err := stream.NewEnvironment(
		stream.NewEnvironmentOptions().
			SetMaxProducersPerClient(maxProducersPerClient).
			SetMaxConsumersPerClient(maxConsumersPerClient).
			SetUris(addresses))
	shared.FailOnError(err, "Failed to create environment")

	// Declare a stream or create a super stream
	superStreamName := "orders_stream"
	err = env.DeclareSuperStream(superStreamName,
		stream.NewPartitionsOptions(3).
			SetMaxLengthBytes(stream.ByteCapacity{}.GB(3)))
	shared.FailOnError(err, "Failed to declare super stream")

	// Create a super stream producer
	consumerName := "consumer-" + os.Args[1]
	sac := stream.NewSingleActiveConsumer(
		func(partition string, isActive bool) stream.OffsetSpecification {
			// This function is called when the consumer is promoted to active
			// or not active anymore
			restart := stream.OffsetSpecification{}.First()
			offset, err := env.QueryOffset(consumerName, partition)
			if err == nil {
				restart = stream.OffsetSpecification{}.Offset(offset + 1)
			}

			addInfo := fmt.Sprintf("The consumer is now active ....Restarting from offset: %s", restart)
			if !isActive {
				addInfo = "The consumer is not active anymore for this partition."
			}

			fmt.Printf("[%s] - Consumer update for: %s. %s\n", time.Now().Format(time.TimeOnly),
				partition, addInfo)

			return restart
		},
	)

	// Create a consumer
	messagesHandler := func(consumerContext stream.ConsumerContext, message *amqp.Message) {
		// Processing the message
		fmt.Printf("Stream: %s - Received message: %s\n", consumerContext.Consumer.GetStreamName(), message.Data)

		// Store the offset after processing the message :: bad practice
		err := consumerContext.Consumer.StoreOffset()
		shared.FailOnError(err, "Failed to store offset")
	}

	// Create a consumer stream
	superStreamConsumer, err := env.NewSuperStreamConsumer(superStreamName, messagesHandler,
		stream.NewSuperStreamConsumerOptions().
			SetSingleActiveConsumer(sac.SetEnabled(EnableSingleActiveConsumer)).
			SetConsumerName(consumerName).
			SetOffset(stream.OffsetSpecification{}.First()))
	shared.FailOnError(err, "Failed to create consumer")

	reader := bufio.NewReader(os.Stdin)
	fmt.Println(" [x] Waiting for messages. enter to close the consumer")
	_, _ = reader.ReadString('\n')
	err = superStreamConsumer.Close()
	shared.FailOnError(err, "Failed to close consumer")
}
