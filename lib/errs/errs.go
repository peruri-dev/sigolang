package errs

import (
	"fmt"
	"log/slog"
	"reflect"

	"github.com/pkg/errors"
)

// for https://github.com/pkg/errors
type stackTracer interface {
	StackTrace() errors.StackTrace
}

func ErrorField(err error) slog.Attr {
	var stack string
	if serr, ok := err.(stackTracer); ok {
		st := serr.StackTrace()
		stack = fmt.Sprintf("%+v", st)
		if len(stack) > 0 && stack[0] == '\n' {
			stack = stack[1:]
		}
	}
	cause := errors.Cause(err)
	return slog.Group("error",
		"kind", reflect.TypeOf(cause).String(),
		"stack", stack,
		"message", err.Error(),
	)
}
