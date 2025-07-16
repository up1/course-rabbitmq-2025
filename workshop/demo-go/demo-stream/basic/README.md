# Basic of Stream queue
* Fanout patterns
  * Offset and consumer name
  * Filter message or routing key
* Super stream
  * Partition the stream
* RabbitMQ stream
  * Port=5552
* [Stream client with Go](https://github.com/rabbitmq/rabbitmq-stream-go-client)


## Install
```
$go mod tidy
```

## Run basic producer
```
$export RABBITMQ_STREAM_URL=rabbitmq-stream://user:password@localhost:5552
$go run producer.go 
```

## Run basic consumer
```
$export RABBITMQ_STREAM_URL=rabbitmq-stream://user:password@localhost:5552
$go run consumer.go A
$go run consumer.go B
```

## Run consumer with offset tracking
* consumer name
* store offset
```
$export RABBITMQ_STREAM_URL=rabbitmq-stream://user:password@localhost:5552
$go run consumer_offset.go A
$go run consumer_offset.go B
```

## Filter messages by consumer
* Producer
  * Add message type = "type1", "type2", "type3"
* Consumer by type

Run producer
```
$export RABBITMQ_STREAM_URL=rabbitmq-stream://user:password@localhost:5552
$go run producer_filter.go 
```

Run consumer per type
```
$export RABBITMQ_STREAM_URL=rabbitmq-stream://user:password@localhost:5552
$go run consumer_filter.go A type1
$go run consumer_filter.go B type2
$go run consumer_filter.go C type3
```

## Single Active Consumer (SAC)

Run producer
```
$export RABBITMQ_STREAM_URL=rabbitmq-stream://user:password@localhost:5552
$go run producer_sac.go
```

Run consumers
```
$export RABBITMQ_STREAM_URL=rabbitmq-stream://user:password@localhost:5552

$go run consumer_sac.go A1
$go run consumer_sac.go A1
$go run consumer_sac.go A1
```

## Super stream

### Run producer to create super stream
* Exchange=`orders_stream` => direct
* Stream
  * orders_stream-0
  * orders_stream-1
  * orders_stream-2

```
$go run producer_super_stream.go

Publishing messages to super stream 'orders_stream'...
Published 1 messages, 0 confirmed, 0 failed
Press enter to close the producer
Message with key: key_0 stored in partition orders_stream-0, total: 1
Message with key: key_3 stored in partition orders_stream-0, total: 2
Message with key: key_4 stored in partition orders_stream-0, total: 3
Message with key: key_7 stored in partition orders_stream-0, total: 4
Message with key: key_2 stored in partition orders_stream-2, total: 5
Message with key: key_5 stored in partition orders_stream-2, total: 6
Message with key: key_6 stored in partition orders_stream-2, total: 7
Message with key: key_1 stored in partition orders_stream-1, total: 8
Message with key: key_8 stored in partition orders_stream-1, total: 9
Message with key: key_9 stored in partition orders_stream-1, total: 10
```

