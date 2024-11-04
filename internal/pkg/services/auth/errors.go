package auth

import "errors"

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrSamePassword      = errors.New("same password")
	ErrWrongPassword     = errors.New("wrong password")
	ErrWrongRole         = errors.New("wrong role")
)
