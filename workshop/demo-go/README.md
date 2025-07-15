# Workshop with RabbitMQ and Go
* Work queue pattern

## Install RabbitMQ + Stream
```
$docker compose up -d rabbitmq
$docker compose ps
```

Enable stream plugin
```
$docker compose exec rabbitmq rabbitmq-plugins enable rabbitmq_stream rabbitmq_stream_management 
```

Go to RabbitMA management Admin UI
* http://localhost:15672
  * user=user
  * password=password

## Basic use cases
* Work queue
* Pub/Sub with fanout exchange
* Routing message
  * Direct exchange
  * Topic exchange

## Failure use cases
* [Re-deliver](https://github.com/up1/course-rabbitmq-2025/tree/main/workshop/demo-go/redeliver)
* Detect failures
  * Retry
  * Guarantee no messages are lost or skipped
  * Dead-Letter Queue (DLQ)

## RabbitMQ Stream
