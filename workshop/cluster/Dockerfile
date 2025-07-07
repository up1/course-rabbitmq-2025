FROM rabbitmq:3-management
COPY ./configs/cluster-entrypoint.sh /usr/local/bin/
COPY ./configs/.erlang.cookie /var/lib/rabbitmq/.erlang.cookie
RUN rabbitmq-plugins enable rabbitmq_prometheus
RUN apt update && apt install -y curl
RUN chmod +x /usr/local/bin/cluster-entrypoint.sh