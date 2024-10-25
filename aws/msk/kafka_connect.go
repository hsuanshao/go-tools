package msk

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kafkaconnect"
	"github.com/aws/aws-sdk-go/service/kafkaconnect/kafkaconnectiface"

	"github.com/hsuanshao/go-tools/aws/entities"
	"github.com/hsuanshao/go-tools/ctx"
)

// MSKConnectConfig defines kafka
type MSKConnectConfig struct {
	IAMUser       *entities.IAMUser `json:"iam_user,omitempty"`
	AssumeRoleARN string            `json:"assume_role_arn"`
	Region        string            `json:"region"`
}

// NewKafkaConnect to new kafkaConnet API
func NewKafkaConnect(ctx ctx.CTX, conectConf *MSKConnectConfig) (mskSrv KafkaConnect) {
	sess := session.Must(session.NewSession(&aws.Config{Region: &conectConf.Region}))
	creds := stscreds.NewCredentials(sess, conectConf.AssumeRoleARN)
	mskAPI := kafkaconnect.New(sess, &aws.Config{Credentials: creds})
	mskSrv = kafkaconnectiface.KafkaConnectAPI(mskAPI)

	return mskSrv
}
