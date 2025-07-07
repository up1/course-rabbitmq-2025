FROM rabbitmq:3-management
COPY ./configs/cluster-entrypoint.sh /usr/local/bin/
COPY ./configs/.erlang.cookie /var/lib/rabbitmq/.erlang.cookie
RUN chmod +x /usr/local/bin/cluster-entrypoint.sh