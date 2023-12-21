package safe

import (
	"runtime/debug"

	"github.com/partyzanex/shortlink/pkg/logger"
	"github.com/pkg/errors"
)

// Recover wraps recover for defer.
func Recover() {
	if r := recover(); r != nil {
		logger.GetLogger().Errorf("panic recovered: %s; stack %s", r, debug.Stack())
	}
}

// RecoverToError writes recover result to error.
func RecoverToError(err *error) { //nolint:gocritic
	if r := recover(); r != nil {
		*err = errors.Errorf("panic recovered: %s; stack %s", r, debug.Stack())
	}
}

// Go run goroutine with recover.
func Go(fn func()) {
	go func() {
		defer Recover()

		fn()
	}()
}

// Fn returns function with recover.
func Fn(fn func() error) func() error {
	return func() (err error) {
		defer RecoverToError(&err)

		return fn()
	}
}
