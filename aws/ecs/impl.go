package ecs

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
	"github.com/sirupsen/logrus"

	"github.com/hsuanshao/go-tools/ctx"
)

func NewECS(ctx ctx.CTX, ecsRegion string) (ecsapi API) {
	awsConfig := &aws.Config{
		Region: aws.String(ecsRegion),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": ecsRegion}).Error("new an aws config session failed")
		return
	}

	ecssess := ecs.New(sess)
	ecsapi = ecsiface.ECSAPI(ecssess)

	return ecsapi
}

func NewECSwithAssumeRole(ctx ctx.CTX, ecsRegion string, assumeRole *credentials.Credentials) (ecsapi API) {
	awsConfig := &aws.Config{
		Region: aws.String(ecsRegion),
	}
	if assumeRole != nil {
		awsConfig.Credentials = assumeRole
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": ecsRegion}).Error("new an aws config session failed")
		return
	}

	ecssess := ecs.New(sess)
	ecsapi = ecsiface.ECSAPI(ecssess)

	return ecsapi
}
