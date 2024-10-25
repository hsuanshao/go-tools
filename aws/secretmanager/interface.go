package secretmanager

import (
	secretIfc "github.com/aws/aws-sdk-go/service/secretsmanager/secretsmanageriface"
)

type Service interface {
	secretIfc.SecretsManagerAPI
}
