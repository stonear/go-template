package validator

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/id"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	id_translations "github.com/go-playground/validator/v10/translations/id"
	"github.com/uptrace/opentelemetry-go-extra/otelzap"
	"go.uber.org/zap"
)

func Load(log *otelzap.Logger) {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		log.Error("[Validator] Failed to get validator engine")
		return
	}

	// Register Indonesian translator
	idLocale := id.New()
	uni := ut.New(idLocale, idLocale)
	trans, ok = uni.GetTranslator("id")
	if !ok {
		log.Error("[Validator] Failed to get Indonesian translator")
		return
	}

	err := id_translations.RegisterDefaultTranslations(v, trans)
	if err != nil {
		log.Error("[Validator] Failed to register Indonesian translations", zap.Error(err))
		return
	}

	// Register custom validation
	err = RegisterCustomValidations(v)
	if err != nil {
		log.Error("[Validator] Failed to register custom validations", zap.Error(err))
	}
}
