package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	New = errors.New
	Is  = errors.Is
	As  = errors.As

	ErrAuthorization    = errors.New("Authorization error")
	ErrIsNotExist       = errors.New("Object is not exist")
	ErrAlreadyExist     = errors.New("Object already exist")
	ErrAlreadyDeleted   = errors.New("Object already deleted")
	ErrAlreadyUpdated   = errors.New("Already updated")
	ErrIncorrectRequest = errors.New("Request is incorrect")
	ErrIncorrectFormat  = errors.New("Incorrect format")
	ErrEncryption       = errors.New("Encryption error")
	ErrGetRuntimeCaller = errors.New("Can't get runtime.Caller")
	ErrIsEmpty          = errors.New("Is empty")
	ErrCustom           = errors.New("Custom error")
	ErrIsUnsupported    = errors.New("Unsupported")
	ErrTimeout          = errors.New("Timeout")
	ErrNotFound         = errors.New("Not found")
	ErrObjectNotInit    = errors.New("Object is not initialized")
	ErrIsInProcessing   = errors.New("Request is in processing")
	ErrIsNotValid       = errors.New("Is not valid")
	ErrAccessTokenExp   = errors.New("Access token expired")
	ErrRefreshTokenExp  = errors.New("Refresh token expired")
	ErrUnauthorized     = errors.New("Unauthorized")
	ErrAccessDenied     = errors.New("Access denied")
)

func Eq(e error, target error) bool {
	// TODO тупо, но я не знаю лучшего варианта
	return e.Error() == target.Error()
}

func ErrNotFoundFn() error {
	// TODO мы не можем использовать ErrNotFound потому что стек будет другой
	return errors.New("Not found")
}

func ErrFieldIsEmpty(name string) error {
	return fmt.Errorf("%w: %s", ErrIsEmpty, name)
}

func ErrFieldIsIncorrect(name string) error {
	return fmt.Errorf("%w: %s", ErrIncorrectFormat, name)
}

func WithStack(e error) error {
	type stackTracer interface {
		StackTrace() errors.StackTrace
	}

	err, ok := errors.Cause(e).(stackTracer)
	if !ok {
		panic("oops, err does not implement stackTracer")
	}

	st := err.StackTrace()
	maxStEl := len(st)
	if maxStEl > 128 {
		maxStEl = 127
	}

	return errors.WithMessagef(e, "%s, stack: %+v\n", e, st[:maxStEl])
}
