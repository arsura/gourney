package validator

import (
	"reflect"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	gpgvalidator "github.com/go-playground/validator/v10"
	entranslations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	Validate *gpgvalidator.Validate
	Trans    ut.Translator
}

func (v *Validator) TransError(err error) []string {
	result := []string{}
	errs := err.(gpgvalidator.ValidationErrors)
	for _, e := range errs {
		result = append(result, e.Translate(v.Trans))
	}
	return result
}

func NewValidator() *Validator {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")
	validator := gpgvalidator.New()
	entranslations.RegisterDefaultTranslations(validator, trans)
	validator.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	return &Validator{validator, trans}
}
