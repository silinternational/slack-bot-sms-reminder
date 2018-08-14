#!/usr/bin/env ash

# Create data table
aws dynamodb create-table --table-name reminderbot_queue --attribute-definitions AttributeName=ID,AttributeType=S --key-schema AttributeName=ID,KeyType=HASH --provisioned-throughput ReadCapacityUnits=50,WriteCapacityUnits=50 --endpoint-url http://dynamo:8000
