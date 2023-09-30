package validator

import (
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

func Load(log *zap.Logger) {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		err := v.RegisterValidation("notEvil", notEvil)
		if err != nil {
			log.Error("[Validator] Failed to register notEvil validation")
		}
	}
}

var notEvil validator.Func = func(fl validator.FieldLevel) bool {
	input := fl.Field().String()
	return !strings.Contains(strings.ToLower(input), "voldemort")
}
