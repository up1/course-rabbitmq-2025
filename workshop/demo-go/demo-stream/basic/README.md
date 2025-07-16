# Basic of Stream queue
* Fanout patterns

## Install
```
$go mod tidy
```

## Run producer
```
$export RABBITMQ_STREAM_URL=rabbitmq-stream://user:password@localhost:5552
$go run producer.go 
```

## Run consumer
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

