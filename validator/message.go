package validator

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
)

func Message(err error) []string {
	var ve validator.ValidationErrors
	if err != nil && errors.As(err, &ve) {
		msg := make([]string, len(ve))
		for i, fe := range ve {
			customMsg := DefaultValidatorMessage[fe.Tag()]
			if customMsg == "" {
				customMsg = fe.Error()
			}
			msg[i] = fmt.Sprintf("%s: %s", fe.Field(), customMsg)
		}
		return msg
	}
	return []string{}
}

var (
	DefaultValidatorMessage = map[string]string{
		"required": "This field is required",
		"notEvil":  "Go away you evil person",
	}
)
