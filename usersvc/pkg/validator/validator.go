package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	v10 "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// Validator is the validator interface
type Validator interface {
	Validate(any) error
}

// validator is the validator struct
type validator struct {
	validator *v10.Validate
	trans     ut.Translator
}

// ValidationError is the validator error
type ValidationError struct {
	Message error
	Errors  map[string]string
}

// Error returns the error message
func (e *ValidationError) Error() string {
	return e.Message.Error()
}

// New creates a new validator
func New() Validator {
	validate := v10.New()
	language := en.New()
	uni := ut.New(language, language)
	trans, _ := uni.GetTranslator("language")

	_ = en_translations.RegisterDefaultTranslations(validate, trans)

	// Register custom validation
	if err := validate.RegisterValidation("numberQueryParam", ValidateNumberOnQueryParam); err != nil {
		return nil
	}

	return &validator{
		validator: validate,
		trans:     trans,
	}
}

// Validate validates the data
func (c *validator) Validate(data any) error {
	if err := c.validator.Struct(data); err != nil {
		result := make(map[string]string)
		for _, e := range err.(v10.ValidationErrors) {
			fieldName := ConvertSnakeCase(e.Field())

			switch e.Tag() {
			case "numberQueryParam":
				result[fieldName] = fmt.Sprintf("The %s field must be a valid number", fieldName)
			default:
				result[fieldName] = strings.Replace(e.Translate(c.trans), e.Field(), fieldName, -1)
			}
		}

		// Extract the first error message from the result map
		var firstErrorMessage string
		for _, msg := range result {
			firstErrorMessage = msg
			break
		}

		return &ValidationError{
			Message: errors.New(firstErrorMessage),
			Errors:  result,
		}
	}

	return nil
}

func ValidateNumberOnQueryParam(fl v10.FieldLevel) bool {
	value := fl.Field().String()

	if value == "" || value == "<int Value>" {
		// Empty string is allowed, as the field is optional
		return true
	}

	// If already integer, return true
	if reflect.TypeOf(value) == reflect.TypeOf(int(0)) {
		return true
	}

	// Try to parse the value to an integer
	if _, err := strconv.Atoi(value); err != nil {
		return false
	}

	return true
}

// ConvertSnakeCase converts the field name to snake case
func ConvertSnakeCase(fieldName string) string {
	var re = regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := re.ReplaceAllString(fieldName, "${1}_${2}")
	return strings.ToLower(snake)
}
