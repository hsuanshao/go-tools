package appconf

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/appconfig"
	"github.com/aws/aws-sdk-go/service/appconfig/appconfigiface"
	"github.com/sirupsen/logrus"

	"github.com/hsuanshao/go-tools/ctx"
)

func NewAppConfig(ctx ctx.CTX, appConfRegion string) (appConfigAPI API) {
	awsConfig := &aws.Config{
		Region: aws.String(appConfRegion),
	}

	sess, err := session.NewSession(awsConfig)

	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": appConfRegion}).Error("new an aws config session failed")
		return
	}

	appSess := appconfig.New(sess)
	appConfigAPI = appconfigiface.AppConfigAPI(appSess)

	return appConfigAPI
}

func NewAppConfigWithAssumeRole(ctx ctx.CTX, appConfRegion string, assumeRole *credentials.Credentials) (appConfigAPI API) {
	awsConfig := &aws.Config{
		Region: aws.String(appConfRegion),
	}
	if assumeRole != nil {
		awsConfig.Credentials = assumeRole
	}

	sess, err := session.NewSession(awsConfig)

	if err != nil {
		ctx.WithFields(logrus.Fields{"err": err, "region": appConfRegion}).Error("new an aws config session failed")
		return
	}

	ecssess := appconfig.New(sess)
	appConfigAPI = appconfigiface.AppConfigAPI(ecssess)

	return appConfigAPI
}
