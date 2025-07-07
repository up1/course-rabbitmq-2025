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