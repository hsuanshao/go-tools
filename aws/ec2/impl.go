package ec2pkg

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
	"github.com/sirupsen/logrus"

	"github.com/hsuanshao/go-tools/aws/entities"
	"github.com/hsuanshao/go-tools/ctx"
)

// InitialEC2Service based on default credentail
func InitialEC2Service(ctx ctx.CTX, region string) Service {
	awsConfig := &aws.Config{
		Region: aws.String(region),
	}
	sess, err := session.NewSession(awsConfig)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "ec2 region": region}).Panic("initial ec2 service failed")
	}

	ec2SVC := ec2.New(sess)
	ec2Srv := ec2iface.EC2API(ec2SVC)

	return ec2Srv
}

// EC2Config defines service
type EC2Config struct {
	Region     string            `json:"region"`
	IAMUser    *entities.IAMUser `json:"iam_user,omitempty"`
	AssumeRole *string           `json:"assume_role_arn,omitempty"`
}

// InitEC2SrvWithAssumeRole for new ec2 api with assume role
func InitEC2SrvWithAssumeRole(ctx ctx.CTX, conf *EC2Config) Service {
	sess, err := session.NewSession(&aws.Config{Region: &conf.Region})
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": conf.Region}).Panic("new aws session failed")
		return nil
	}
	var service *ec2.EC2
	service = ec2.New(sess)

	if conf.AssumeRole != nil && strings.TrimSpace(*conf.AssumeRole) != "" {
		ctx.Info("with assume role")
		creds := stscreds.NewCredentials(sess, *conf.AssumeRole)
		service = ec2.New(sess, &aws.Config{Credentials: creds})
	}

	ec2Srv := ec2iface.EC2API(service)
	return ec2Srv
}
