package errors

import (
	"errors"
	"fmt"
	"github.com/hashicorp/go-multierror"
)

func New(message string, a ...any) error {
	return errors.New(fmt.Sprintf(message, a...))
}

func NewMulti(errs []error, message string, a ...any) error {
	if len(errs) == 0 {
		return nil
	}

	var multiError *multierror.Error
	for _, err := range errs {
		multiError = multierror.Append(err)
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(message, a...), multiError)
}

func Wrap(err error, message string, a ...any) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", fmt.Sprintf(message, a...), err)
}

func ToString(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}
