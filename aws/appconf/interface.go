package appconf

import (
	"github.com/aws/aws-sdk-go/service/appconfig/appconfigiface"
)

/**
About how to Apply AWS.AppConfig to deploy cotainer to a Fragate
please check this document references

https://aws.amazon.com/blogs/mt/application-configuration-deployment-to-container-workloads-using-aws-appconfig/
*/

type API interface {
	appconfigiface.AppConfigAPI
}
