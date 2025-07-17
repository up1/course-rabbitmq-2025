# Working with Work queue
* Direct exchange

## QoS (Quality of Service) in RabbitMQ
* It refers to the mechanism that allows you to manage how messages are delivered to consumers
* how many unacknowledged messages are allowed on a channel at any given time ?

## Start consumer A
```
$export RABBITMQ_URL=amqp://user:password@localhost:5672
$go run consumer.go
```

## Start consumer B
```
$export RABBITMQ_URL=amqp://user:password@localhost:5672
$go run consumer.go
```

## Start producer
```
$export RABBITMQ_URL=amqp://user:password@localhost:5672
$go run producer.go
```