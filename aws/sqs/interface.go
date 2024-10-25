package sqs

import (
	"github.com/aws/aws-sdk-go/service/sqs/sqsiface"
)

// Queue is SQS service interface
type Queue interface {
	sqsiface.SQSAPI
}
