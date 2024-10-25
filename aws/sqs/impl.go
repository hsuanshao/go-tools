package sqs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"

	"github.com/hsuanshao/go-tools/ctx"
)

// NewSQS to new a SQS Queue Service
func NewSQS(sqsRegion string) Queue {
	ctx := ctx.Background()

	awsConfig := &aws.Config{
		Region: aws.String(sqsRegion),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		ctx.WithField("err", err).Panic("initial aws session for aws sqs failed")
	}

	sqsapi := sqs.New(sess, aws.NewConfig().WithRegion(sqsRegion))

	return sqsapi
}

func NewSQSwithAssumeRole(sqsRegion string, assumeRole *credentials.Credentials) Queue {
	ctx := ctx.Background()

	awsConfig := &aws.Config{
		Region: aws.String(sqsRegion),
	}
	if assumeRole != nil {
		awsConfig.Credentials = assumeRole
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		ctx.WithField("err", err).Panic("initial aws session for aws sqs failed")
	}

	sqsapi := sqs.New(sess, aws.NewConfig().WithRegion(sqsRegion))

	return sqsapi
}
