# Install Single RabbitMQ Server
* [Monitoring system](https://www.rabbitmq.com/docs/prometheus)
  * Prometheus
  * Grafana

## Build image
```
$docker compose build
```

## Start RabbitMQ server
```
$docker compose up -d rabbit_node1
$docker compose ps
```

Access to RabbitMQ Management UI
* http://localhost:15672/
  * user=user
  * password=password

For AMQP client
* port=5672

For Stream client
* port=5552

## Working with application metric with Prometheus
* https://www.rabbitmq.com/docs/prometheus

Access to `rabbit_node1` and get metric data
```
$docker compose exec -it rabbit_node1 bash
$curl -s localhost:15692/metrics | head -n 3

# TYPE erlang_mnesia_held_locks gauge
# HELP erlang_mnesia_held_locks Number of held locks.
erlang_mnesia_held_locks 0
```

Start prometheus server
```
$docker compose -f docker-compose-metric.yml up -d prometheus
$docker compose -f docker-compose-metric.yml ps
```

Access to Prometheus UI
* http://localhost:9090
* http://localhost:9090/targets

## Working with Grafana dashboard
* https://grafana.com/grafana/dashboards/10991-rabbitmq-overview/
* https://grafana.com/grafana/dashboards/14798-rabbitmq-stream/
```
$docker compose -f docker-compose-metric.yml up -d grafana
$docker compose -f docker-compose-metric.yml ps
```

Access to Grafana UI
* http://localhost:3000
  * user=admin
  * password=admin


## Delete all resources
```
$docker compose down
$docker compose -f docker-compose-metric.yml down

$docker volume prune
```