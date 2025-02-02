package validator

import (
	"errors"
	"fmt"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

var (
	trans ut.Translator

	DefaultValidatorMessage = map[string]string{
		"notEvil": "Go away you evil person",
	}
)

func Message(err error) []string {
	var ve validator.ValidationErrors
	if err != nil && errors.As(err, &ve) {
		msg := make([]string, len(ve))
		for i, fe := range ve {
			// First try custom message
			customMsg := DefaultValidatorMessage[fe.Tag()]
			if customMsg == "" && trans != nil {
				// Then try translation
				customMsg = fe.Translate(trans)
			}
			if customMsg == "" {
				// Finally fallback to default error
				customMsg = fe.Error()
			}
			msg[i] = fmt.Sprintf("%s: %s", fe.Field(), customMsg)
		}
		return msg
	}
	return []string{}
}
