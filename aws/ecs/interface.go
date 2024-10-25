package ecs

import (
	"github.com/aws/aws-sdk-go/service/ecs/ecsiface"
)

// API provide ECSAPI
type API interface {
	ecsiface.ECSAPI
}
