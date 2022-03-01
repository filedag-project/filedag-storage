package berrors

import "errors"

var (
	ErrConfigNotFound = errors.New("config file not found")
)
