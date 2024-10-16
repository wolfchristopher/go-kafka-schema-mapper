# thoth
Copyright (c) 2024 Christopher Wolf. All rights reserved.
This software is proprietary and confidential.
Unauthorized copying of this file, via any medium, is strictly prohibited.

Add Kafka config handler for the "/config" endpoint
```
curl -X POST http://localhost:8080/kafka_config -H "Content-Type: application/json" 
-d '{
    "brokers": "localhost:9092",
    "group_id": "client-consumer-group",
    "pre_service_topic": "input-topic",
    "post_service_topic": "output-topic",
    "security_protocol": "client-security-protocol",
    "sasl_mechanism": "sasl_mechanism",
    "sasl_username": "sasl_username",
    "sasl_password": "sasl_password",
    "sslc_a_location": "sslc_a_location",
    "auto_offset_reset": "earliest"
}'
```
Register the schema handler for the "/schema" endpoint
```
curl -X POST http://localhost:8080/schema 
	-H "Content-Type: application/json" 
	-d '{
		"fields": [
			{"name": "username", "type": "string", "required": true},
			{"name": "age", "type": "integer", "required": false}
			],
			"metadata": {
			"version": 1,
			"description": "User data schema"
			}
		}'
```