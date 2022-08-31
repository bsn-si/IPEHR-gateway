package errors

import (
	"errors"
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
)
