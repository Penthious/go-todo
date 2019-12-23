package domain

import (
	"errors"
	"fmt"
)

var (
	ErrNoResult                     = errors.New("No results")
	ErrUserWithEmailAlreadyExist    = errors.New("User with email already exist")
	ErrUserWithUsernameAlreadyExist = errors.New("User with username already exist")
	ErrEmailInvalid                 = errors.New("Email is invalid")
	ErrInvalidCredentials           = errors.New("Email or Password did not match")
)

type ErrNotLongEnough struct {
	field  string
	length int
}

type ErrIsRequired struct {
	field string
}

type ErrMustMatch struct {
	field string
}

func (e ErrNotLongEnough) Error() string {
	return fmt.Sprintf("%v not long enough, %d characters is required", e.field, e.length)
}

func (e ErrIsRequired) Error() string {
	return fmt.Sprintf("%v is a required field", e.field)
}

func (e ErrMustMatch) Error() string {
	return fmt.Sprintf("%v must match", e.field)
}
