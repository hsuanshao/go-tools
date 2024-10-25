/*
Package ctx extends standard context to support logging.
For context detail, see https://golang.org/pkg/context/
*/
package ctx

//
import (
	"context"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/hsuanshao/go-tools/validator"
)

var (
	// logExcept set keys that should be except from log
	logExcept = []string{"accessToken", "password", "authorization", "Authorization"}

	// settleLogger as shared parameter for handle and control logger as expected setup
	settleLogger *logrus.Logger
)

// CTX extends Google's context to support logging methods.
type CTX struct {
	context.Context
	logrus.FieldLogger
}

// Background returns a non-nil, empty Context. It is never canceled, has no values, and
// has no deadline. It is typically used by the main function, initialization, and tests,
// and as the top-level Context for incoming requests
func Background() CTX {
	if settleLogger == nil {
		settleLogger = logrus.StandardLogger()
		h, err := NewHook(&Config{LogFormat: &logrus.TextFormatter{ForceColors: true}})
		if err != nil {
			settleLogger.WithField("err", err).Fatal("prepare log hook failed")
		}
		settleLogger.AddHook(h)
	}
	return CTX{
		Context:     context.Background(),
		FieldLogger: settleLogger,
	}
}

// InjectContext return a copy of native context
func InjectContext(parent context.Context, ctx CTX) CTX {

	return CTX{
		Context:     context.WithValue(parent, nil, nil),
		FieldLogger: ctx.FieldLogger,
	}
}

// WithValue returns a copy of parent in which the value associated with key is val.
func WithValue(parent CTX, key string, val interface{}) CTX {
	var keyi interface{} = key
	if validator.IsInStringSlice(logExcept, key) {
		return CTX{
			Context:     context.WithValue(parent, keyi, val),
			FieldLogger: parent.FieldLogger,
		}
	}

	return CTX{
		Context:     context.WithValue(parent, keyi, val),
		FieldLogger: parent.FieldLogger.WithField(key, val),
	}
}

// WithValues returns a copy of parent in which the values associated with keys are vals.
func WithValues(parent CTX, kvs map[string]interface{}) CTX {
	c := parent
	for k, v := range kvs {
		c = WithValue(c, k, v)
	}

	return c
}

// WithCancel returns a copy of parent with added cancel function
func WithCancel(parent CTX) (CTX, context.CancelFunc) {
	newCtx, cFunc := context.WithCancel(parent)
	return CTX{
		Context:     newCtx,
		FieldLogger: parent.FieldLogger,
	}, cFunc
}

// WithTimeout returns a copy of parent with timeout condition
// and cancel function
func WithTimeout(parent CTX, d time.Duration) (CTX, context.CancelFunc) {
	newCtx, cFunc := context.WithTimeout(parent, d)
	return CTX{
		Context:     newCtx,
		FieldLogger: parent.FieldLogger,
	}, cFunc
}
