package validate

import (
	"context"
)

type (
	StructValidator struct {
		reflectValue
		rules           []FieldValidator
		fieldNameGetter FieldFormatter
		failFast        bool
	}
)

func Struct(target any, rules ...FieldValidator) *StructValidator {
	return &StructValidator{
		reflectValue:    newReflectValue(target),
		rules:           rules,
		fieldNameGetter: DefaultFieldFormatter,
		failFast:        DefaultFailFastBehavior,
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

func (v *StructValidator) Validate(ctx context.Context) error {
	if len(v.rules) == 0 {
		return nil
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
