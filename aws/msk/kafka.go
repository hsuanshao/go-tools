package msk

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kafka"
	"github.com/aws/aws-sdk-go/service/kafka/kafkaiface"

	"github.com/hsuanshao/go-tools/ctx"
)

// MSKConfig is the contract to New Kafka management
type MSKConfig struct {
	AssumeRoleARN string `json:"assume_role_arn"`
	Region        string `json:"region"`
}

// NewKafkaMaster provide aws-go-sdk-v1 test able solution
func NewKafkaMaster(ctx ctx.CTX, mskConfig *MSKConfig) (mskSrv Kafka) {
	sess := session.Must(session.NewSession(&aws.Config{Region: &mskConfig.Region}))
	creds := stscreds.NewCredentials(sess, mskConfig.AssumeRoleARN)
	mskAPI := kafka.New(sess, &aws.Config{Credentials: creds})
	mskSrv = kafkaiface.KafkaAPI(mskAPI)

	return mskSrv
}
