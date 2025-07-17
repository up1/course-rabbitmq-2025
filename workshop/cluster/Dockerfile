FROM rabbitmq:4-management
COPY ./configs/cluster-entrypoint.sh /usr/local/bin/
COPY ./configs/.erlang.cookie /var/lib/rabbitmq/.erlang.cookie
RUN chmod 400 /var/lib/rabbitmq/.erlang.cookie
RUN rabbitmq-plugins enable rabbitmq_prometheus rabbitmq_management rabbitmq_stream rabbitmq_stream_management
RUN apt update && apt install -y curl
RUN chmod +x /usr/local/bin/cluster-entrypoint.sh