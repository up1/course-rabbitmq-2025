# Working with Dead Letter Queue(DLQ)
* Detect failures (consumer errors, publish nacks, unroutable messages, connection drops)
* Retry up to N times (with headers and/or retry queues)
* Guarantee no messages are lost or skipped
* Log or route to a Dead-Letter Queue (DLQ) after exceeding retry attempts

## Config Producer
* Direct exchange
  * main-exchange => main-queue
  * retry-exchange => retry-queue with TTL and back to main-exchange (max 3 times)
  * dlx-exchange => dlq (retry more than 3 times)

Run
```
$go run producer.go
```

Results
* Exchange = direct
  * dlx-exchange
  * main-exchange
  * retry-exchange
* Queue
  * dlq
  * main-queue
  * retry-queue

## Start Consumer for main-exchange 
```
$go run consumer.go
```

## Start Consumer for dlx-exchange 
```
$go run consumer_dlq.go
```


