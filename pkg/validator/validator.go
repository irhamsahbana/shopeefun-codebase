package validator

import (
	"reflect"
	"strings"

	// "github.com/go-playground/locales/en"
	// ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	// en_translations "github.com/go-playground/validator/v10/translations/en"
	"github.com/rs/zerolog/log"
)

type Validator struct {
	// trans     ut.Translator
	validator *validator.Validate
}

func NewValidator() *Validator {
	validatorCustom := &Validator{}

	// en := en.New()
	// uni := ut.New(en, en)
	// trans, _ := uni.GetTranslator("en")

	v := validator.New()
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		var name string

		name = strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("query"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("form"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("params"), ",", 2)[0]
		}

		if name == "" {
			name = strings.SplitN(fld.Tag.Get("prop"), ",", 2)[0]
		}

		if name == "-" {
			return ""
		}

		return name
	})

	// en_translations.RegisterDefaultTranslations(v, trans)
	if err := v.RegisterValidation("email_blacklist", isEmailBlacklist); err != nil {
		log.Fatal().Err(err).Msg("Error while registering email_blacklist validator")
	}
	if err := v.RegisterValidation("strong_password", isStrongPassword); err != nil {
		log.Fatal().Err(err).Msg("Error while registering strong_password validator")
	}
	if err := v.RegisterValidation("unique_in_slice", isUniqueInSlice); err != nil {
		log.Fatal().Err(err).Msg("Error while registering unique validator")
	}

	validatorCustom.validator = v
	// validatorCustom.trans = trans

	return validatorCustom
}

func (v *Validator) Validate(i any) error {
	return v.validator.Struct(i)
}

// blacklist email validator
func isEmailBlacklist(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	disallowedDomains := []string{"gmail", "yahoo", "outlook", "hotmail", "aol", "live", "inbox", "icloud", "mail", "gmx", "yandex"}

	for _, domain := range disallowedDomains {
		if strings.Contains(email, domain) {
			return false
		}
	}

	return true
}

func isStrongPassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 12 {
		return false
	}

	hasUppercase := false
	hasLowercase := false
	hasNumber := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUppercase = true
		case char >= 'a' && char <= 'z':
			hasLowercase = true
		case char >= '0' && char <= '9':
			hasNumber = true
		}
	}

	return hasUppercase && hasLowercase && hasNumber
}

func isUniqueInSlice(fl validator.FieldLevel) bool {
	// Get the slice from the FieldLevel interface
	val := fl.Field()

	// Ensure the field is a slice
	if val.Kind() != reflect.Slice {
		return false
	}

	// Use a map to check for duplicates
	elements := make(map[interface{}]bool)
	for i := 0; i < val.Len(); i++ {
		elem := val.Index(i).Interface()
		if _, found := elements[elem]; found {
			return false // Duplicate found
		}
		elements[elem] = true
	}
	return true
}
