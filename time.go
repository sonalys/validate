package validate

import (
	"context"
	"time"
)

type TimeValidator struct {
	*rules
}

func Time(ptr any) *TimeValidator {
	return &TimeValidator{
		rules: newRules(ptr),
	}
}

func (v *TimeValidator) Optional() *TimeValidator {
	v.rules.optional = true
	return v
}

func (v *TimeValidator) Before(t time.Time) *TimeValidator {
	v.rules.Append(func(ctx context.Context) error {
		if v.isZero {
			return nil
		}

		if !v.valueOf.Interface().(time.Time).Before(t) {
			return BeforeError{
				Value: t,
			}
		}

		return nil
	})

	return v
}

func (v *TimeValidator) After(t time.Time) *TimeValidator {
	v.rules.Append(func(ctx context.Context) error {
		if v.isZero {
			return nil
		}

		if !v.valueOf.Interface().(time.Time).After(t) {
			return AfterError{
				Value: t,
			}
		}

		return nil
	})

	return v
}

func (v *TimeValidator) Between(start, end time.Time) *TimeValidator {
	v.rules.Append(func(ctx context.Context) error {
		if v.isZero {
			return nil
		}

		if !v.valueOf.Interface().(time.Time).After(start) || !v.valueOf.Interface().(time.Time).Before(end) {
			return RangeError{
				Min: start,
				Max: end,
			}
		}

		return nil
	})

	return v
}
