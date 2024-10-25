package mq

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/mq"
	awsmq "github.com/aws/aws-sdk-go/service/mq/mqiface"
	"github.com/sirupsen/logrus"

	"github.com/hsuanshao/go-tools/ctx"
)

// NewMQ is the constructor to new AWS MQ API,
func NewMQ(ctx ctx.CTX, mqRegion string) (mqAPI API) {
	awsConfig := &aws.Config{
		Region: aws.String(mqRegion),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": mqRegion}).Error("new an aws config session failed")
		return
	}

	mqsess := mq.New(sess)
	mqSrv := awsmq.MQAPI(mqsess)

	return mqSrv
}

func NewMQwithAssumeRole(ctx ctx.CTX, mqRegion string, assumeRole *credentials.Credentials) (mqAPI API) {
	awsConfig := &aws.Config{
		Region: aws.String(mqRegion),
	}
	if assumeRole != nil {
		awsConfig.Credentials = assumeRole
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": mqRegion}).Error("new an aws config session failed")
		return
	}

	mqsess := mq.New(sess)
	mqSrv := awsmq.MQAPI(mqsess)

	return mqSrv
}
