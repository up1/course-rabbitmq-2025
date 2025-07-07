# Workshop :: Create RabbitMQ cluster


## Services in Docker Compose

| Service     | Description               |
| ----------- | ------------------------- |
| `rabbitmq1` | RabbitMQ (cluster)        |
| `rabbitmq2` | RabbitMQ (cluster member) |
| `rabbitmq3` | RabbitMQ (cluster member) |
| `haproxy`   | Load Balancer             |

## Build Image of RabbitMQ
```
$docker compose build
```

Config of RabbitMQ
* [Erlang cookie](https://www.rabbitmq.com/docs/clustering#erlang-cookie)
  * Use for the nodes in cluster communicate with each other

Config of Load balancer [HA Proxy](http://www.haproxy.org/)
* File `configs/haproxy.cfg`


## Start RabbitMQ cluster and Load Balancer
```
$docker compose up -d haproxy

$docker compose ps
NAME                  IMAGE               COMMAND                  SERVICE     CREATED         STATUS                   PORTS
cluster-haproxy-1     haproxy:3.2.2       "docker-entrypoint.s…"   haproxy     2 minutes ago   Up About a minute        0.0.0.0:5672->5672/tcp, [::]:5672->5672/tcp, 0.0.0.0:15672->15672/tcp, [::]:15672->15672/tcp
cluster-rabbitmq1-1   cluster-rabbitmq1   "/usr/local/bin/clus…"   rabbitmq1   2 minutes ago   Up 2 minutes (healthy)   4369/tcp, 5671-5672/tcp, 15671-15672/tcp, 15691-15692/tcp, 25672/tcp
cluster-rabbitmq2-1   cluster-rabbitmq2   "/usr/local/bin/clus…"   rabbitmq2   2 minutes ago   Up 2 minutes (healthy)   4369/tcp, 5671-5672/tcp, 15671-15672/tcp, 15691-15692/tcp, 25672/tcp
cluster-rabbitmq3-1   cluster-rabbitmq3   "/usr/local/bin/clus…"   rabbitmq3   2 minutes ago   Up 2 minutes (healthy)   4369/tcp, 5671-5672/tcp, 15671-15672/tcp, 15691-15692/tcp, 25672/tcp
```

Access to RabbitMQ Management UI
* http://localhost:15672/
  * user=guest
  * password=guest

For AMQP client
* port=5672

## Working with application metric with Prometheus

Access to `rabbitmq1` and get metric data
```
$docker compose exec -it rabbitmq1 bash
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

## Delete all resources
```
$docker compose down
$docker volume prune
```