package msk

import (
	"github.com/aws/aws-sdk-go/service/kafka/kafkaiface"
	"github.com/aws/aws-sdk-go/service/kafkaconnect/kafkaconnectiface"
)

// Kafka defineis aws-sdk-go v1 kafka interface
type Kafka interface {
	kafkaiface.KafkaAPI
}

// KafkaConnect defines aws msk kafka connect interface
type KafkaConnect interface {
	kafkaconnectiface.KafkaConnectAPI
}
