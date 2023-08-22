# producer
kafka-console-producer.sh \
--producer.config /opt/bitnami/kafka/config/producer.properties \
--bootstrap-server kafka:9092 --topic my-topic

# consumer
kafka-console-consumer.sh \
--consumer.config /opt/bitnami/kafka/config/consumer.properties \
--bootstrap-server kafka:9092 --topic my-topic --from-beginning

docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"