package db

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/silinternational/slack-bot-sms-reminder"
	"os"
	"testing"
)

const ENV_DYNAMO_ENDPOINT = "AWS_DYNAMODB_ENDPOINT"

var db *dynamodb.DynamoDB

func GetDb() *dynamodb.DynamoDB {
	if db == nil {
		dynamoEndpoint := os.Getenv(ENV_DYNAMO_ENDPOINT)
		fmt.Fprintf(os.Stdout, "dynamodb endpoint: %s\n", dynamoEndpoint)
		db = dynamodb.New(session.New(), aws.NewConfig().WithRegion("us-east-1").WithEndpoint(dynamoEndpoint))
	}
	return db
}

func PutItem(item interface{}) error {
	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		reminderbot.ServerError(fmt.Errorf("failed to DynamoDB marshal Record, %v", err))
	}
	input := &dynamodb.PutItemInput{
		TableName: aws.String(reminderbot.GetDbTableName(reminderbot.DynamoDBTableName)),
		Item:      av,
	}

	db := GetDb()
	_, err = db.PutItem(input)
	return err
}

func DeleteItem(id string) (bool, error) {

	returnOldValues := "ALL_OLD"
	// Prepare the input for the query.
	input := &dynamodb.DeleteItemInput{
		TableName: aws.String(reminderbot.GetDbTableName(reminderbot.DynamoDBTableName)),
		Key: map[string]*dynamodb.AttributeValue{
			"ID": {
				S: aws.String(id),
			},
		},
		ReturnValues: &returnOldValues,
	}

	db := GetDb()
	// Delete the item from DynamoDB. I
	resp, err := db.DeleteItem(input)
	if err != nil {
		return false, err
	}

	// resp.Attributes contains attribute of old record before deletion, if empty the original item was not found
	if len(resp.Attributes) == 0 {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func ListMessages() ([]reminderbot.SmsMessage, error) {

	tableName := reminderbot.GetDbTableName(reminderbot.DynamoDBTableName)
	input := &dynamodb.ScanInput{
		TableName: &tableName,
	}

	db := GetDb()
	var results []map[string]*dynamodb.AttributeValue
	err := db.ScanPages(input,
		func(page *dynamodb.ScanOutput, lastPage bool) bool {
			results = append(results, page.Items...)
			return !lastPage
		})

	if err != nil {
		return []reminderbot.SmsMessage{}, err
	}

	var messages []reminderbot.SmsMessage
	err = dynamodbattribute.UnmarshalListOfMaps(results, &messages)
	if err != nil {
		return []reminderbot.SmsMessage{}, err
	}

	return messages, nil
}

func FlushTables(t *testing.T) {
	tables := []string{reminderbot.DynamoDBTableName}
	db := GetDb()

	for _, tableName := range tables {
		input := &dynamodb.ScanInput{
			TableName: &tableName,
		}

		var results []map[string]*dynamodb.AttributeValue
		err := db.ScanPages(input,
			func(page *dynamodb.ScanOutput, lastPage bool) bool {
				results = append(results, page.Items...)
				return !lastPage
			})

		if err != nil {
			t.Error(err)
			t.Fail()
		}

		for _, item := range results {

			var keyCriteria map[string]*dynamodb.AttributeValue
			keyCriteria = map[string]*dynamodb.AttributeValue{
				"ID": {
					S: aws.String(*item["ID"].S),
				},
			}

			deleteInput := &dynamodb.DeleteItemInput{
				TableName: aws.String(tableName),
				Key:       keyCriteria,
			}

			_, err := db.DeleteItem(deleteInput)
			if err != nil {
				t.Errorf("Unable to delete item ID: %s, from table %s. Error: %s", *item["ID"].S, tableName, err.Error())
				t.Fail()
			}
		}
	}
}
