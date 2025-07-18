# Working with RabbitMQ Stream PerfTest
* https://rabbitmq.github.io/rabbitmq-stream-perf-test/stable/htmlsingle/


## 1. Run with Java
* [Binary release](https://github.com/rabbitmq/rabbitmq-stream-perf-test/releases)
```
$wget https://github.com/rabbitmq/rabbitmq-java-tools-binaries-dev/releases/download/v-stream-perf-test-latest/stream-perf-test-latest.jar
```

Run
```
$java -jar stream-perf-test-latest.jar --uris rabbitmq-stream://user:password@localhost:5552

$java -jar stream-perf-test-latest.jar --uris rabbitmq-stream://user:password@localhost:5552 --producers 1 --consumers 5 --rate 10000

$java -jar stream-perf-test-latest.jar --uris rabbitmq-stream://user:password@localhost:5552 --producers 1 --consumers 5 --rate 10000 --prometheus
```