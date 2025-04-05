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

	Rule           func(context.Context) error
	FieldFormatter func(value reflect.StructField) string

	rules struct {
		reflectValue
		rules    []Rule
		optional bool
	}
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

func newRules(value any) *rules {
	return &rules{
		reflectValue: newReflectValue(value),
		rules:        make([]Rule, 0),
	}
}

func (r *rules) Validate(ctx context.Context) error {
	if len(r.rules) == 0 || r.isZero && r.optional {
		return nil
	}

	if r.isZero && !r.optional {
		return ErrFieldRequired
	}

	errs := errorPool.Get().([]error)
	errs = errs[:0]
	defer errorPool.Put(errs)

	for _, rule := range r.rules {
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

func (r *rules) Append(rule Rule) {
	r.rules = append(r.rules, rule)
}
