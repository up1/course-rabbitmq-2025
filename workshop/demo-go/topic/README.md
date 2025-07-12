# Working with Routing key with wildcard
* Topic exchange

## Start consumer A :: Receive all messages
```
$export RABBITMQ_URL=amqp://user:password@localhost:5672
$go run consumer.go "#"
```

## Start consumer B :: Receive message with A
```
$export RABBITMQ_URL=amqp://user:password@localhost:5672
$go run consumer.go "A.*"
```

## Start consumer C :: Receive message with A and B
```
$export RABBITMQ_URL=amqp://user:password@localhost:5672
$go run consumer.go "A.*" "*.B"
```

## Start producer
* With routing key
```
$export RABBITMQ_URL=amqp://user:password@localhost:5672

$go run producer.go A
$go run producer.go B
```