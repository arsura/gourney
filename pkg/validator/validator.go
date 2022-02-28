package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	gpg_validator "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	Validate *gpg_validator.Validate
	Trans    ut.Translator
}

func (v *Validator) TransError(err error) []string {
	result := []string{}
	errs := err.(gpg_validator.ValidationErrors)
	for _, e := range errs {
		result = append(result, e.Translate(v.Trans))
	}
	return result
}

func NewValidator() *Validator {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validator := gpg_validator.New()
	en_translations.RegisterDefaultTranslations(validator, trans)
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &Validator{validator, trans}
}
