package validate

import (
	"context"
	"fmt"
	"net/mail"
	"regexp"
	"slices"
)

type StringValidator struct {
	*rules
}

func String(ptr any) *StringValidator {
	return &StringValidator{
		rules: newRules(ptr),
	}
}

func (v *StringValidator) Optional() *StringValidator {
	v.rules.optional = true
	return v
}

func (v *StringValidator) NotEmpty() *StringValidator {
	v.rules.Append(v.ruleNotEmpty())
	return v
}

func (v *StringValidator) MinLength(min int) *StringValidator {
	v.rules.Append(v.ruleMinLength(min))
	return v
}

func (v *StringValidator) MaxLength(max int) *StringValidator {
	v.rules.Append(v.ruleMaxLength(max))
	return v
}

func (v *StringValidator) Length(min, max int) *StringValidator {
	v.rules.Append(v.ruleLength(min, max))
	return v
}

func (v *StringValidator) Matches(pattern string) *StringValidator {
	v.rules.Append(func(ctx context.Context) error {
		if !matches(v.valueOf.String(), pattern) {
			return PatternError{
				ShouldMatch: true,
				Pattern:     pattern,
			}
		}

		return nil
	})

	return v
}

func (v *StringValidator) NotMatches(pattern string) *StringValidator {
	v.rules.Append(func(ctx context.Context) error {
		if matches(v.valueOf.String(), pattern) {
			return fmt.Errorf("must not match pattern %w", PatternError{
				ShouldMatch: false,
				Pattern:     pattern,
			})
		}

		return nil
	})

	return v
}

func matches(value, pattern string) bool {
	regex, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Sprintf("invalid regex pattern: %s", pattern))
	}

	return regex.MatchString(value)
}

func (v *StringValidator) In(values ...string) *StringValidator {
	v.rules.Append(func(ctx context.Context) error {
		if !slices.Contains(values, v.valueOf.String()) {
			return fmt.Errorf("must be one of %v", values)
		}

		return nil
	})

	return v
}

func (v *StringValidator) NotIn(values ...string) *StringValidator {
	v.rules.Append(func(ctx context.Context) error {
		if slices.Contains(values, v.valueOf.String()) {
			return fmt.Errorf("must not be one of %v", values)
		}

		return nil
	})

	return v
}

func (v *StringValidator) IsEmail() *StringValidator {
	v.rules.Append(func(ctx context.Context) error {
		if _, err := mail.ParseAddress(v.valueOf.String()); err != nil {
			return fmt.Errorf("must be a valid email address: %w", err)
		}

		return nil
	})

	return v
}
