package eks

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/eks"
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
	"github.com/sirupsen/logrus"

	"github.com/hsuanshao/go-tools/ctx"
)

func NewEKS(ctx ctx.CTX, eksRegion string) (eksapi API) {
	awsConfig := &aws.Config{
		Region: aws.String(eksRegion),
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": eksRegion}).Error("new an aws config session failed")
		return
	}

	ekssess := eks.New(sess)
	ekssrv := eksiface.EKSAPI(ekssess)

	return ekssrv
}

func NewEKSwithAssumeRole(ctx ctx.CTX, eksRegion string, assumeRole *credentials.Credentials) (eksapi API) {
	awsConfig := &aws.Config{
		Region: aws.String(eksRegion),
	}
	if assumeRole != nil {
		awsConfig.Credentials = assumeRole
	}

	sess, err := session.NewSession(awsConfig)
	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": eksRegion}).Error("new an aws config session failed")
		return
	}

	ekssess := eks.New(sess)
	ekssrv := eksiface.EKSAPI(ekssess)

	return ekssrv
}
