FROM confluentinc/cp-schema-registry:7.7.1

USER root

RUN chmod -R 777 /etc/schema-registry

USER 1000

CMD ["/etc/confluent/docker/run"]