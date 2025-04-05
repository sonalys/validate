# Validate

A lightweight and extensible Go validation library for validating structs and fields with ease.

## Installation

```bash
go get github.com/sonalys/validate
```

## Features

- Validate strings, numbers, and time fields with various rules.
- Support for custom field name formatting.
- Fail-fast or collect-all-errors behavior.
- Easy-to-use API for struct validation.

## Usage

### Validating Strings

```go
package main

import (
	"context"
	"fmt"
	"github.com/sonalys/validate"
)

func main() {
	name := "John"
	validation := validate.String(&name).
		NotEmpty().
		MinLength(3).
		MaxLength(50)

	err := validation.Validate(context.Background())
	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed!")
	}
}
```

### Validating Numbers

```go
package main

import (
	"context"
	"fmt"
	"github.com/sonalys/validate"
)

func main() {
	age := 25
	validation := validate.Number[int](&age).
		Min(18).
		Max(60)

	err := validation.Validate(context.Background())
	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed!")
	}
}
```

### Validating Time

```go
package main

import (
	"context"
	"fmt"
	"time"
	"github.com/sonalys/validate"
)

func main() {
	birthday := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	validation := validate.Time(&birthday).
		Before(time.Now()).
		After(time.Now().AddDate(-100, 0, 0))

	err := validation.Validate(context.Background())
	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed!")
	}
}
```

### Validating Structs

```go
package main

import (
	"context"
	"fmt"
	"time"
	"github.com/sonalys/validate"
)

type User struct {
	Name     string     `json:"name"`
	Email    string     `json:"email"`
	Age      int        `json:"age"`
	Birthday *time.Time `json:"birthday"`
}

func main() {
	birthday := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	user := User{
		Name:     "Alice",
		Email:    "alice@example.com",
		Age:      30,
		Birthday: &birthday,
	}

	validation := validate.Struct(&user,
		validate.String(&user.Name).NotEmpty().MinLength(3),
		validate.String(&user.Email).IsEmail(),
		validate.Number[int](&user.Age).Range(18, 60),
		validate.Time(&user.Birthday).Optional().Before(time.Now()),
	)

	err := validation.Validate(context.Background())
	if err != nil {
		fmt.Println("Validation failed:", err)
	} else {
		fmt.Println("Validation passed!")
	}
}
```

## License

This project is licensed under the MIT License.