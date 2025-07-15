# Use case :: Re-delivery
* Consumer crash
* Consumer that will re-read (redeliver) any messages that weren’t ACK’d before the crash
* Consumer crashes and redeliver unacknowledged messages


## Config Producer
* Uses Persistent delivery mode so messages survive broker restarts
* Declares the queue as durable

## Config Consumer
* autoAck=false ensures you ACK only after processing
* If the process crashes before d.Ack, RabbitMQ will requeue the message
* The d.Redelivered flag tells you if this is a retry

## Start consumer
```
$export RABBITMQ_URL=amqp://user:password@localhost:5672
$go run consumer.go
```

## Start producer
```
$export RABBITMQ_URL=amqp://user:password@localhost:5672
$go run producer.go
```