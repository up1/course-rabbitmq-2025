# Workshop :: Create RabbitMQ cluster


## Services in Docker Compose

| Service     | Description               |
| ----------- | ------------------------- |
| `rabbitmq1` | RabbitMQ (cluster)        |
| `rabbitmq2` | RabbitMQ (cluster member) |
| `rabbitmq3` | RabbitMQ (cluster member) |
| `haproxy`   | Load Balancer             |