package validate

import (
	"fmt"
	"strings"
)

type (
	// FieldValidator represents an error that occurred while validating a field.
	// It contains the field name and the error message.
	// Err might be a MultiError.
	FieldError struct {
		Field string
		Value any
		Err   error
	}

	// MultiError represents a collection of errors that occurred during validation.
	MultiError []error

	MinValueError struct {
		Value any
	}

	MaxValueError struct {
		Value any
	}

	RangeError struct {
		Min any
		Max any
	}

	PatternError struct {
		ShouldMatch bool
		Pattern     string
	}

	MinLengthError struct {
		Min     int
		Current int
	}

	MaxLengthError struct {
		Max     int
		Current int
	}

	LengthError struct {
		Min     int
		Max     int
		Current int
	}

	BeforeError struct {
		Value any
	}

	AfterError struct {
		Value any
	}
)

func (e FieldError) Error() string {
	return fmt.Sprintf("%s (%v)", e.Field, e.Err)
}

func (e MultiError) Error() string {
	messages := make([]string, 0, len(e))

	for _, err := range e {
		if err == nil {
			continue
		}

		messages = append(messages, err.Error())
	}

	return strings.Join(messages, "; ")
}

func (e MultiError) Unwrap() []error {
	return e
}

func (e MinValueError) Error() string {
	return fmt.Sprintf("must be at least %v", e.Value)
}

func (e MaxValueError) Error() string {
	return fmt.Sprintf("must be smaller than %v", e.Value)
}

func (e RangeError) Error() string {
	return fmt.Sprintf("must be between %v and %v", e.Min, e.Max)
}

func (e PatternError) Error() string {
	if e.ShouldMatch {
		return fmt.Sprintf("must match pattern %s", e.Pattern)
	}

	return fmt.Sprintf("must not match pattern %s", e.Pattern)
}

func (e MinLengthError) Error() string {
	return fmt.Sprintf("must be at least %d characters long", e.Min)
}

func (e MaxLengthError) Error() string {
	return fmt.Sprintf("must be at most %d characters long", e.Max)
}

func (e LengthError) Error() string {
	return fmt.Sprintf("must be between %d and %d characters long", e.Min, e.Max)
}

func (e BeforeError) Error() string {
	return fmt.Sprintf("must be before %v", e.Value)
}

func (e AfterError) Error() string {
	return fmt.Sprintf("must be after %v", e.Value)
}
