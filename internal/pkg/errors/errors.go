package errors

import (
	"errors"
	"fmt"
)

var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
)

// Wrap adds context to an error.
func Wrap(err error, message string) error {
	return fmt.Errorf("%s: %w", message, err)
}

func New(message string) error {
	return errors.New(message)
}
