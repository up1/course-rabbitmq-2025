# Basic of Stream queue
* Fanout patterns

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

