# Workshop with RabbitMQ

## Manage prefetch count
* Receive up to N unacknowledged message
* Classic queue = QoS = AMQP protocol(port=5672)
  * prefetch_count = N
  * consumer model = Push
  * use case = tasks
* Stream queue (port=5552)
  * SetInitialCredits(N)
  * consumer model = pull (credit-based)
  * use case = high-throughput logs

## 1. Class queue


## 2. Stream queue with Single Active Consumer (SAC)

Run producer
```
$export RABBITMQ_STREAM_URL=rabbitmq-stream://user:password@localhost:5552
$go run producer_offset.go
```

Run consumers with offset
```
$export RABBITMQ_STREAM_URL=rabbitmq-stream://user:password@localhost:5552

$go run consumer_offset.go A1
$go run consumer_offset.go A2
```

Run consumers with Single Active Consumer (SAC)
```
$export RABBITMQ_STREAM_URL=rabbitmq-stream://user:password@localhost:5552

$go run consumer_sac.go A1
$go run consumer_sac.go A1
$go run consumer_sac.go A1
```
