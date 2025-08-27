package validation

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

type ValidationError struct {
	Tag   string      `json:"tag"`
	Value interface{} `json:"value"`
}

func init() {
	validate = validator.New()
	registerRules(validate)
}

// ValidateStruct used to validate struct
func ValidateStruct(s interface{}) error {
	return validate.Struct(s)
}

func IsValidateError(err error) bool {
	if _, ok := err.(validator.ValidationErrors); ok {
		return true
	}
	return false
}

func ValidateError(err error) map[string]ValidationError {
	errorMap := map[string]ValidationError{}
	for _, err := range err.(validator.ValidationErrors) {
		field := err.Namespace()
		ind := strings.Index(err.Namespace(), ".")
		if ind != -1 {
			field = field[ind+1:]
		}
		errorMap[field] = ValidationError{
			Tag:   err.Tag(),
			Value: err.Param(),
		}
	}
	return errorMap
}
