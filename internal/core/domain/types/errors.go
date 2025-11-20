package types

import "errors"

var (
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

var ErrRideNotFound = errors.New("ride not found")

var (
	ErrInternalServiceError = errors.New("internal service error")
	ErrDriverExists         = errors.New("driver already exists")
	ErrDriverNotFound       = errors.New("driver not found")
	ErrDriverOnline         = errors.New("driver is online")
	ErrDriverStatusNotAllow = errors.New("the status of the driver does not allow")
)
