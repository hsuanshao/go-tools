package secretmanager

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	secretIfc "github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"

	"github.com/hsuanshao/go-tools/ctx"
	"github.com/sirupsen/logrus"
)

// NewSM for new secrets manager api based on default aws credentials
func NewSM(ctx ctx.CTX, awsRegion string) (smSrv Service) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": awsRegion}).Error("new an aws session failed")
		return nil
	}

	smAPI := secretsmanager.New(sess)
	smSrv = secretIfc.SecretsManagerAPI(smAPI)
	return smSrv
}

// NewSMWithAssumeRole for new secrets manager api with assume role
func NewSMWithAssumeRole(ctx ctx.CTX, awsRegion, assumeRoleARN string) Service {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(awsRegion),
	})
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": awsRegion}).Error("new an aws session failed")
		return nil
	}

	if !strings.Contains(assumeRoleARN, ":role/") {
		ctx.WithField("assumeRoleARN", assumeRoleARN).Error("invalidated assume role arn had been input")
		return nil
	}

	creds := stscreds.NewCredentials(sess, assumeRoleARN)

	assumeRoleSess := sess
	assumeRoleSess.Config.Credentials = creds
	assumeRoleSess.Config.Region = aws.String(awsRegion)

	smAPI := secretsmanager.New(assumeRoleSess)
	smSrv := secretIfc.SecretsManagerAPI(smAPI)

	return smSrv
}
