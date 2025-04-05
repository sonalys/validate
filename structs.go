package validate

import (
	"context"
	"fmt"
)

type (
	StructValidator struct {
		reflectValue
		rules           []FieldValidator
		fieldNameGetter FieldFormatter
		failFast        bool
		optional        bool
	}
)

func Struct(target any, rules ...FieldValidator) *StructValidator {
	reflectValue := newReflectValue(target)

	validateFields(reflectValue, rules...)

	return &StructValidator{
		reflectValue:    reflectValue,
		rules:           rules,
		fieldNameGetter: DefaultFieldFormatter,
		failFast:        DefaultFailFastBehavior,
	}
}

func validateFields(value reflectValue, rules ...FieldValidator) {
	startOffset := value.ptr

	var biggestOffset uintptr

	for i := range value.typeOf.NumField() {
		fieldTypeOf := value.typeOf.Field(i)
		if fieldTypeOf.Offset > biggestOffset {
			biggestOffset = fieldTypeOf.Offset
		}
	}

	endOffset := startOffset + biggestOffset

	for i, rule := range rules {
		ptr := rule.pointer()
		if ptr < startOffset || ptr > endOffset {
			panic(fmt.Sprintf("field %d does not belong to the struct", i))
		}
	}
}

func (v *StructValidator) SetFieldNameFormatter(formatter FieldFormatter) *StructValidator {
	v.fieldNameGetter = formatter
	return v
}

func (v *StructValidator) SetFailFast(failFast bool) *StructValidator {
	v.failFast = failFast
	return v
}

func (v *StructValidator) getFieldNameValue(fieldPtr uintptr) (string, any) {
	for i := range v.typeOf.NumField() {
		fieldTypeOf := v.typeOf.Field(i)
		if fieldTypeOf.Offset+v.ptr != fieldPtr {
			continue
		}

		fieldValueOf := v.valueOf.Field(i)
		return v.fieldNameGetter(fieldTypeOf), fieldValueOf.Interface()
	}

	return "", nil
}

func (v *StructValidator) Optional() *StructValidator {
	v.optional = true
	return v
}

func (v *StructValidator) Validate(ctx context.Context) error {
	if len(v.rules) == 0 || v.optional && v.isZero {
		return nil
	}

	if !v.optional && v.isZero {
		return ErrFieldRequired
	}

	errs := errorPool.Get().([]error)
	errs = errs[:0]
	defer errorPool.Put(errs)

	for _, rule := range v.rules {
		err := rule.Validate(ctx)
		if err == nil {
			continue
		}

		name, value := v.getFieldNameValue(rule.pointer())

		err = FieldError{
			Field: name,
			Value: value,
			Err:   err,
		}

		if v.failFast {
			return err
		}

		errs = append(errs, err)
	}

	if len(errs) == 0 {
		return nil
	}

	return MultiError(errs)
}
