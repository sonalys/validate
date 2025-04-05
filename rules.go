package validate

import (
	"context"
	"fmt"
	"reflect"
)

type reflectValue struct {
	valueOf reflect.Value
	typeOf  reflect.Type
	isZero  bool
	ptr     uintptr
}

func newReflectValue(target any) reflectValue {
	valueOf := reflect.ValueOf(target)
	value, isZero := getValue(target)

	var typeOf reflect.Type
	if !isZero {
		typeOf = value.Type()
	}

	return reflectValue{
		valueOf: value,
		typeOf:  typeOf,
		isZero:  isZero,
		ptr:     valueOf.Pointer(),
	}
}

func (v reflectValue) pointer() uintptr {
	return v.ptr
}

func getValue(target any) (reflect.Value, bool) {
	value := reflect.ValueOf(target)

	for value.Kind() == reflect.Ptr {
		if value.IsNil() {
			return reflect.Value{}, true
		}

		value = value.Elem()
	}

	return value, false
}

func (v reflectValue) ruleNotEmpty() Rule {
	return func(ctx context.Context) error {
		if v.isZero || v.valueOf.Len() == 0 {
			return fmt.Errorf("value must not be empty")
		}
		return nil
	}
}

func (v reflectValue) ruleMinLength(minLength int) Rule {
	length := 0
	if !v.isZero {
		length = v.valueOf.Len()
	}
	return func(ctx context.Context) error {
		if length < minLength {
			return MinLengthError{
				Min:    minLength,
				Length: length,
			}
		}
		return nil
	}
}

func (v reflectValue) ruleMaxLength(maxLength int) Rule {
	length := 0
	if !v.isZero {
		length = v.valueOf.Len()
	}
	return func(ctx context.Context) error {
		if length > maxLength {
			return MaxLengthError{
				Max:    maxLength,
				Length: length,
			}
		}
		return nil
	}
}

func (v reflectValue) ruleLength(minLength, maxLength int) Rule {
	length := 0
	if !v.isZero {
		length = v.valueOf.Len()
	}
	return func(ctx context.Context) error {
		if length < minLength || length > maxLength {
			return LengthError{
				Min:    minLength,
				Max:    maxLength,
				Length: length,
			}
		}
		return nil
	}
}
