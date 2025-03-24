package errors

import (
	"errors"
	"fmt"
)

func New(msg string) error {
	return errors.New(msg)
}

func Wrap(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}

func NewF(format string, args ...interface{}) error {
	return fmt.Errorf(format, args...)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func Join(errs ...error) error {
	return errors.Join(errs...)
}
