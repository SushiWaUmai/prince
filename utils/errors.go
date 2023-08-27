package utils

import "errors"

var (
	ErrNotEnoughArguments = errors.New("not enough arguments")
	ErrNotEnoughPermission = errors.New("not enough permission")
)
