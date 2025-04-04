package validate

import (
	"context"
	"reflect"
	"sync"
)

type (
	FieldValidator interface {
		Validate(context.Context) error
		pointer() uintptr
	}

	ValidatorFunc  func(context.Context) error
	FieldFormatter func(value reflect.StructField) string

	rules []ValidatorFunc
)

// FieldFormatterStructName is the default field formatter for new struct validators.
var DefaultFieldFormatter FieldFormatter = FieldFormatterStructName

// DefaultFailFastBehavior determines if the validation should stop on the first error found.
var DefaultFailFastBehavior bool = false

func FieldFormatterStructName(value reflect.StructField) string {
	return value.Name
}

func FieldFormatterTag(tag string) FieldFormatter {
	return func(value reflect.StructField) string {
		if tagValue, ok := value.Tag.Lookup(tag); ok {
			return tagValue
		}
		return FieldFormatterStructName(value)
	}
}

var errorPool = sync.Pool{
	New: func() any {
		return make([]error, 0)
	},
}

func (r rules) Validate(ctx context.Context) error {
	if len(r) == 0 {
		return nil
	}

	errs := errorPool.Get().([]error)
	errs = errs[:0]
	defer errorPool.Put(errs)

	for _, rule := range r {
		if err := rule(ctx); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	multiErr := MultiError(errs)

	return multiErr
}
