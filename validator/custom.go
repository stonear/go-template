package validator

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// RegisterCustomValidations registers all custom validations
func RegisterCustomValidations(v *validator.Validate) error {
	// notEvil is a custom validator that checks if the input contains the word "voldemort"
	if err := v.RegisterValidation("notEvil", func(fl validator.FieldLevel) bool {
		input := fl.Field().String()
		return !strings.Contains(strings.ToLower(input), "voldemort")
	}); err != nil {
		return err
	}

	// Register more custom validations here...

	return nil
}
