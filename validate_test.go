package validate_test

import (
	"testing"
	"time"

	"github.com/sonalys/validate"
	"github.com/stretchr/testify/require"
)

func Benchmark(b *testing.B) {
	type TestStruct struct {
		Name     string     `json:"name"`
		Email    string     `json:"personal_email"`
		Age      int        `json:"age"`
		Birthday *time.Time `json:"birthday"`
	}

	validate.DefaultFieldFormatter = validate.FieldFormatterTag("json")

	birthday := time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC)

	test := TestStruct{
		Name:     "Joe",
		Email:    "example@domain.com",
		Age:      18,
		Birthday: &birthday,
	}

	validation := validate.Struct(&test,
		validate.String(&test.Name).
			MinLength(3),
		validate.String(&test.Email).
			MinLength(3).
			MaxLength(250),
		validate.Number(&test.Age).
			Range(1, 18),
		validate.Time(&test.Birthday).
			Before(time.Now()).
			After(time.Now().Add(-time.Hour*24*365*18)).
			Between(time.Now().Add(-time.Hour*24*365*18), time.Now()),
	)

	ctx := b.Context()

	b.Run("Benchmark_StructValidator", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_ = validation.Validate(ctx)
		}
	})
}

func Test_Struct(t *testing.T) {
	type TestStruct struct {
		Name     string     `json:"name"`
		Email    string     `json:"personal_email"`
		Age      int        `json:"age"`
		Birthday *time.Time `json:"birthday"`
	}

	validate.DefaultFieldFormatter = validate.FieldFormatterTag("json")

	birthday := time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC)

	test := TestStruct{
		Name:     "Joe",
		Email:    "example@domain.com",
		Age:      18,
		Birthday: &birthday,
	}

	validation := validate.Struct(&test,
		validate.String(&test.Name).
			MinLength(3),
		validate.String(&test.Email).
			MinLength(3).
			MaxLength(250),
		validate.Number(&test.Age).
			Range(1, 18),
		validate.Time(&test.Birthday).
			Between(time.Now().AddDate(-99, 0, 0), time.Now().AddDate(-18, 0, 0)),
	)

	ctx := t.Context()

	err := validation.Validate(ctx)
	require.NoError(t, err)
}
