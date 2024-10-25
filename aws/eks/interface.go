package eks

import (
	"github.com/aws/aws-sdk-go/service/eks/eksiface"
)

// API provide EKSAPI
type API interface {
	eksiface.EKSAPI
}
