package ses

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
	"github.com/hsuanshao/go-tools/ctx"
)

func NewSES(sesRegion string) SimpleEmailService {
	ctx := ctx.Background()

	awsConfig := &aws.Config{
		Region: aws.String(sesRegion),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		ctx.WithField("err", err).Panic("initial aws simple email service failed")
	}

	svc := ses.New(sess, aws.NewConfig().WithRegion(sesRegion))
	return svc
}

func NewSESwithAssumeRole(sesRegion string, assumeRole *credentials.Credentials) SimpleEmailService {
	ctx := ctx.Background()

	awsConfig := &aws.Config{
		Region: aws.String(sesRegion),
	}
	if assumeRole != nil {
		awsConfig.Credentials = assumeRole
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		ctx.WithField("err", err).Panic("initial aws simple email service failed")
	}

	svc := ses.New(sess, aws.NewConfig().WithRegion(sesRegion))
	return svc
}
