package validate

import (
	"context"

	"golang.org/x/exp/constraints"
)

type NumberValidator[T constraints.Ordered] struct {
	reflectValue
	rules
}

func Number[T constraints.Ordered](ptr *T) *NumberValidator[T] {
	return &NumberValidator[T]{
		reflectValue: newReflectValue(ptr),
	}
}

func min[T constraints.Ordered](value, min T) error {
	if value < min {
		return MinValueError{
			Value: min,
		}
	}

	return nil
}

func max[T constraints.Ordered](value, max T) error {
	if value > max {
		return MaxValueError{
			Value: max,
		}
	}

	return nil
}

func checkRange[T constraints.Ordered](value, min, max T) error {
	if value < min || value > max {
		return RangeError{
			Min: min,
			Max: max,
		}
	}

	return nil
}

func (v *NumberValidator[T]) Min(minValue T) *NumberValidator[T] {
	v.rules = append(v.rules, func(ctx context.Context) error {
		if v.isZero {
			return nil
		}

		return min(v.valueOf.Interface().(T), minValue)
	})

	return v
}

func (v *NumberValidator[T]) Max(maxValue T) *NumberValidator[T] {
	v.rules = append(v.rules, func(ctx context.Context) error {
		if v.isZero {
			return nil
		}

		return max(v.valueOf.Interface().(T), maxValue)
	})

	return v
}

func (v *NumberValidator[T]) Range(minValue, maxValue T) *NumberValidator[T] {
	v.rules = append(v.rules, func(ctx context.Context) error {
		if v.isZero {
			return nil
		}

		return checkRange(v.valueOf.Interface().(T), minValue, maxValue)
	})

	return v
}
