package dynamo

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/sirupsen/logrus"

	"github.com/hsuanshao/go-tools/aws/entities"
	"github.com/hsuanshao/go-tools/ctx"
)

// NewDynamoDB to launch a dynamoAPI
func NewDynamoDB(ctx ctx.CTX, dynamoRegion string) (dynamoAPI API) {
	awsConfig := &aws.Config{
		Region: aws.String(dynamoRegion),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "dynamo region": dynamoRegion}).Error("new a aws session failed")
		return nil
	}

	dynamoSess := dynamodb.New(sess)
	dynamoAPI = dynamodbiface.DynamoDBAPI(dynamoSess)

	return dynamoAPI
}

type Config struct {
	IAMUser    *entities.IAMUser `json:"iam_user,omitempty"`
	AssumeRole *string           `json:"assume_role_arn,omitempty"`
	Region     string            `json:"region"`
}

// NewDynamoDBWithAssumeRole ....
func NewDynamoDBWithAssumeRole(ctx ctx.CTX, conf *Config) API {
	sess, err := session.NewSession(&aws.Config{Region: &conf.Region})
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "dynamo region": conf.Region}).Error("new a aws session failed")
		return nil
	}
	dynamoAPI := dynamodb.New(sess)

	if conf.AssumeRole != nil && strings.TrimSpace(*conf.AssumeRole) != "" {
		creds := stscreds.NewCredentials(sess, *conf.AssumeRole)
		dynamoAPI = dynamodb.New(sess, &aws.Config{Credentials: creds})
	}

	dynamoService := dynamodbiface.DynamoDBAPI(dynamoAPI)
	return dynamoService
}
