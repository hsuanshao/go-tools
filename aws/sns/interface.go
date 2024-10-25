package sns

import (
	"github.com/aws/aws-sdk-go/service/sns/snsiface"
)

type Notification interface {
	snsiface.SNSAPI
}
