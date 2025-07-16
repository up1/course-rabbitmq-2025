package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"shared"
	"sync/atomic"

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

	// Declare a stream or create a super stream
	superStreamName := "orders_stream"
	err = env.DeclareSuperStream(superStreamName,
		stream.NewPartitionsOptions(3).
			SetMaxLengthBytes(stream.ByteCapacity{}.GB(3)))
	shared.FailOnError(err, "Failed to declare super stream")

	// Create a super stream producer
	superStreamProducer, err := env.NewSuperStreamProducer(superStreamName,
		stream.NewSuperStreamProducerOptions(
			stream.NewHashRoutingStrategy(func(message message.StreamMessage) string {
				return message.GetMessageProperties().MessageID.(string)
			})).SetClientProvidedName("my-super-stream-producer"))
	shared.FailOnError(err, "Failed to create producer")

	// Start the producer and handle confirmations
	confirmed, failed := publishMessages(superStreamProducer, 10)
	fmt.Printf("Published %d messages, %d confirmed, %d failed\n", 10, atomic.LoadInt32(&confirmed), atomic.LoadInt32(&failed))

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Press enter to close the producer")
	_, _ = reader.ReadString('\n')
	err = superStreamProducer.Close()
	shared.FailOnError(err, "Failed to close producer")
	fmt.Println("Super stream producer closed.")
}

func publishMessages(producer *stream.SuperStreamProducer, count int) (int32, int32) {
	var confirmed, failed int32

	go handleConfirmations(producer.NotifyPublishConfirmation(1), &confirmed, &failed)

	fmt.Printf("Publishing messages to super stream...\n")
	for i := 0; i < count; i++ {
		message := shared.Message{
			ID:      fmt.Sprintf("msg-%d", i),
			Content: fmt.Sprintf("This is message %d", i),
		}
		messageJson, err := json.Marshal(message)
		shared.FailOnError(err, "Failed to marshal message")

		msg := amqp.NewMessage(messageJson)
		msg.Properties = &amqp.MessageProperties{
			MessageID: fmt.Sprintf("key_%d", i),
		}
		err = producer.Send(msg)
		shared.FailOnError(err, "Failed to send message")
	}

	return confirmed, failed
}

func handleConfirmations(ch <-chan stream.PartitionPublishConfirm, confirmed *int32, failed *int32) {
	for superStreamPublishConfirm := range ch {
		for _, confirm := range superStreamPublishConfirm.ConfirmationStatus {
			if confirm.IsConfirmed() {
				fmt.Printf("Message with key: %s stored in partition %s, total: %d\n",
					confirm.GetMessage().GetMessageProperties().MessageID,
					superStreamPublishConfirm.Partition,
					atomic.AddInt32(confirmed, 1))
			} else {
				atomic.AddInt32(failed, 1)
				fmt.Printf("Message failed to be stored in partition %s\n", superStreamPublishConfirm.Partition)
			}
		}
	}
}
