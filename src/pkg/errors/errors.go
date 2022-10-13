package errors

import (
	"errors"
	"fmt"
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
	Err500              = errors.New("Something is wrong")
	ErrIsNotValid       = errors.New("Is not valid")
)

func ErrFieldIsEmpty(name string) error {
	return fmt.Errorf("%w: %s", ErrIsEmpty, name)
}

func ErrFieldIsIncorrect(name string) error {
	return fmt.Errorf("%w: %s", ErrIncorrectFormat, name)
}
