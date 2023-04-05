package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

var (
	New    = errors.New
	Errorf = errors.Errorf
	Is     = errors.Is
	As     = errors.As
	Wrap   = errors.Wrap

	ErrAuthorization       = errors.New("Authorization error")
	ErrIsNotExist          = errors.New("Object is not exist")
	ErrAlreadyExist        = errors.New("Object already exist")
	ErrAlreadyDeleted      = errors.New("Object already deleted")
	ErrAlreadyUpdated      = errors.New("Already updated")
	ErrIncorrectRequest    = errors.New("Request is incorrect")
	ErrIncorrectFormat     = errors.New("Incorrect format")
	ErrEncryption          = errors.New("Encryption error")
	ErrGetRuntimeCaller    = errors.New("Can't get runtime.Caller")
	ErrIsEmpty             = errors.New("Is empty")
	ErrCustom              = errors.New("Custom error")
	ErrIsUnsupported       = errors.New("Unsupported")
	ErrTimeout             = errors.New("Timeout")
	ErrNotFound            = errors.New("Not found")
	ErrInternalServerError = errors.New("Internal Server Error")
	ErrObjectNotInit       = errors.New("Object is not initialized")
	ErrIsInProcessing      = errors.New("Request is in processing")
	ErrIsNotValid          = errors.New("Is not valid")
	ErrNotImplemented      = errors.New("Not implemented")
	ErrAccessTokenExp      = errors.New("Access token expired")
	ErrRefreshTokenExp     = errors.New("Refresh token expired")
	ErrUnauthorized        = errors.New("Unauthorized")
	ErrAccessDenied        = errors.New("Access denied")
	ErrTypeNotValid        = errors.New("Type is not valid")
)

func ErrFieldIsEmpty(name string) error {
	return fmt.Errorf("%w '%s'", ErrIsEmpty, name)
}

func ErrFieldIsIncorrect(name string) error {
	return fmt.Errorf("%w: %s", ErrIncorrectFormat, name)
}

func ErrObjectWithIDIsNotExist(name string, ID string) error {
	return fmt.Errorf("%w: %s '%s'", ErrIsNotExist, name, ID)
}
