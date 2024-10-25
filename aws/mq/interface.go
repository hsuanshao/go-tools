package mq

import (
	awsmq "github.com/aws/aws-sdk-go/service/mq/mqiface"
)

// API provide mockable interface for test code
type API interface {
	awsmq.MQAPI
}
