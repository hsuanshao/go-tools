package ses

import (
	"github.com/aws/aws-sdk-go/service/ses/sesiface"
)

type SimpleEmailService interface {
	sesiface.SESAPI
}
