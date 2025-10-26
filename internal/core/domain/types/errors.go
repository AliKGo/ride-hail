package types

import "errors"

var (
	ErrIncorrectPassword = errors.New("incorrect password")
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

var ErrRideNotFound = errors.New("ride not found")
