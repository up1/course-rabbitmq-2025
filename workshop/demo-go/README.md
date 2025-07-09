# Workshop with RabbitMQ and Go
* Work queue pattern

## Install RabbitMQ + Stream
```
$docker compose up -d rabbitmq
$docker compose ps
```

Enable stream plugin
```
$docker compose exec rabbitmq rabbitmq-plugins enable rabbitmq_stream rabbitmq_stream_management 
```

Go to RabbitMA management Admin UI
* http://localhost:15672
  * user=user
  * password=password
