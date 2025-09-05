// File: internal/util/validator.go
package util

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// a single instance of the validator
var validate = validator.New()

// ValidationError wraps the validator's errors into a single error.
type ValidationError struct {
	Errors []string
}

func (v *ValidationError) Error() string {
	return strings.Join(v.Errors, ", ")
}

// ValidateStruct performs validation on a struct's fields based on its tags.
func ValidateStruct(s interface{}) error {
	err := validate.Struct(s)
	if err != nil {
		var validationErrors ValidationError
		for _, err := range err.(validator.ValidationErrors) {
			// Create a more user-friendly error message
			field := err.Field()
			tag := err.Tag()
			message := fmt.Sprintf("Field '%s': validation failed on the '%s' tag", field, tag)
			validationErrors.Errors = append(validationErrors.Errors, message)
		}
		return &validationErrors
	}
	return nil
}