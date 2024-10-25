package dynamo

import (
	dyapi "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

// API provides DynamoDBAPI
type API interface {
	dyapi.DynamoDBAPI
}
