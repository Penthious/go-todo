package domain

import "errors"

var (
	ErrNoResult                     = errors.New("No results")
	ErrUserWithEmailAlreadyExist    = errors.New("User with email already exist")
	ErrUserWithUsernameAlreadyExist = errors.New("User with username already exist")
)
