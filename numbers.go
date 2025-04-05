package validate

import (
	"cmp"
	"context"
)

type NumberValidator[T cmp.Ordered] struct {
	*rules
}

func Number[T cmp.Ordered, N **T | *T](ptr N) *NumberValidator[T] {
	return &NumberValidator[T]{
		rules: newRules(ptr),
	}
}

func (v *NumberValidator[T]) Optional() *NumberValidator[T] {
	v.rules.optional = true
	return v
}

func min[T cmp.Ordered](value, min T) error {
	if value < min {
		return MinValueError{
			Value: min,
		}
	}

	return nil
}

func max[T cmp.Ordered](value, max T) error {
	if value > max {
		return MaxValueError{
			Value: max,
		}
	}

	return nil
}

func checkRange[T cmp.Ordered](value, min, max T) error {
	if value < min || value > max {
		return RangeError{
			Min: min,
			Max: max,
		}
	}

	return nil
}

func (v *NumberValidator[T]) Min(minValue T) *NumberValidator[T] {
	v.rules.Append(func(ctx context.Context) error {
		if v.isZero {
			return nil
		}

		return min(v.valueOf.Interface().(T), minValue)
	})

	return v
}

func (v *NumberValidator[T]) Max(maxValue T) *NumberValidator[T] {
	v.rules.Append(func(ctx context.Context) error {
		if v.isZero {
			return nil
		}

		return max(v.valueOf.Interface().(T), maxValue)
	})

	return v
}

func (v *NumberValidator[T]) Range(minValue, maxValue T) *NumberValidator[T] {
	v.rules.Append(func(ctx context.Context) error {
		if v.isZero {
			return nil
		}

		return checkRange(v.valueOf.Interface().(T), minValue, maxValue)
	})

	return v
}
