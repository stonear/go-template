package validator

import (
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Load() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("notEvil", notEvil)
	}
}

var notEvil validator.Func = func(fl validator.FieldLevel) bool {
	input := fl.Field().String()
	return !strings.Contains(strings.ToLower(input), "voldemort")
}
