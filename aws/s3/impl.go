package s3

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/sirupsen/logrus"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/hsuanshao/go-tools/ctx"
)

// NewS3 ...
func NewS3(ctx ctx.CTX, s3Region string) (bkt Bucket) {
	awsConfig := &aws.Config{
		Region: aws.String(s3Region),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": s3Region}).Error("new an aws config session failed")
		return
	}

	s3API := s3.New(sess)
	s3srv := s3iface.S3API(s3API)

	return s3srv
}

// NewS3withAssumeRole ...
func NewS3withAssumeRole(ctx ctx.CTX, s3Region string, assumeRoleARN string) (bkt Bucket) {
	awsConfig := &aws.Config{
		Region: aws.String(s3Region),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": s3Region}).Error("new an aws config session failed")
		return
	}

	if !strings.Contains(assumeRoleARN, ":role/") {
		ctx.WithField("assumeRoleARN", assumeRoleARN).Error("invalidated assume role arn had been input")
		return nil
	}

	creds := stscreds.NewCredentials(sess, assumeRoleARN)

	assumeRoleSess := sess
	assumeRoleSess.Config.Credentials = creds
	assumeRoleSess.Config.Region = aws.String(s3Region)

	s3API := s3.New(assumeRoleSess)
	s3srv := s3iface.S3API(s3API)

	return s3srv
}

// NewS3WithCredential ...
func NewS3WithCredential(ctx ctx.CTX, region, accessKeyID, secretAccessKey string, assumeRoleARN *string) Bucket {
	if strings.TrimSpace(accessKeyID) == "" {
		ctx.Error("access key id is required field")
		return nil
	}
	if strings.TrimSpace(secretAccessKey) == "" {
		ctx.Error("secret access key is required field")
		return nil
	}

	awsConf := &aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessKeyID, secretAccessKey, ""),
	}

	sess, err := session.NewSession(awsConf)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": region}).Error("aws config for new sesssion failed")
		return nil
	}

	if assumeRoleARN != nil {
		if !strings.Contains(*assumeRoleARN, ":role/") {
			ctx.WithField("assume role arn", *assumeRoleARN).Error("input assume role arn format is invalidated")
			return nil
		}

		creds := stscreds.NewCredentials(sess, *assumeRoleARN)
		sess.Config.Credentials = creds
	}

	s3API := s3.New(sess)
	s3Srv := s3iface.S3API(s3API)
	return s3Srv
}
