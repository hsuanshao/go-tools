package sns

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"

	"github.com/hsuanshao/go-tools/ctx"
)

func NewSNS(snsRegion string) Notification {
	ctx := ctx.Background()
	awsConfig := &aws.Config{
		Region: aws.String(snsRegion),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		ctx.WithField("err", err).Panic("initial aws session for secret manager failed")
	}

	svc := sns.New(sess, aws.NewConfig().WithRegion(snsRegion))
	return svc
}

func NewSNSwithAssumeRole(snsRegion string, assumeRole *credentials.Credentials) Notification {
	ctx := ctx.Background()
	awsConfig := &aws.Config{
		Region: aws.String(snsRegion),
	}
	if assumeRole != nil {
		awsConfig.Credentials = assumeRole
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		ctx.WithField("err", err).Panic("initial aws session for secret manager failed")
	}

	svc := sns.New(sess, aws.NewConfig().WithRegion(snsRegion))
	return svc
}
