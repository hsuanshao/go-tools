package ec2pkg

import (
	"github.com/aws/aws-sdk-go/service/ec2/ec2iface"
)

// Service provide EC2 API from aws-sdk-go v1
type Service interface {
	ec2iface.EC2API
}
