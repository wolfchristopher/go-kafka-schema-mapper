#!/bin/bash

# Wait until Kafka is ready
sleep 10

# Create the 'transactions' topic if it doesn't exist
kafka-topics.sh --create \
  --bootstrap-server localhost:9092 \
  --replication-factor 1 \
  --partitions 1 \
  --topic transactions || echo "Topic 'transactions' already exists"
