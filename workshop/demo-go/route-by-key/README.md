# Working with Routing key
* Direct exchange

## Start consumer A
```
$export RABBITMQ_URL=amqp://user:password@localhost:5672
$go run consumer.go A
```

## Start consumer B
```
$export RABBITMQ_URL=amqp://user:password@localhost:5672
$go run consumer.go A B
```

## Start producer
* With routing key
```
$export RABBITMQ_URL=amqp://user:password@localhost:5672

$go run producer.go A
$go run producer.go B
```