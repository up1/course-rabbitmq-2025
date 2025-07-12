# Working with Pub/Sub
* Fanout exchange

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