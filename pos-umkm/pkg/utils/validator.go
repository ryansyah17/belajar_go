package utils

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func ValidateStruct(s interface{}) []ValidationError {
	var errors []ValidationError

	err := validate.Struct(s)
	if err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errors = append(errors, ValidationError{
				Field:   strings.ToLower(e.Field()),
				Message: formatMessage(e),
			})
		}
	}
	return errors
}

func formatMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return fmt.Sprintf("%s wajib diisi", e.Field())
	case "email":
		return "Format email tidak valid"
	case "min":
		return fmt.Sprintf("%s minimal %s karakter", e.Field(), e.Param())
	case "max":
		return fmt.Sprintf("%s maksimal %s karakter", e.Field(), e.Param())
	case "gt":
		return fmt.Sprintf("%s harus lebih dari %s", e.Field(), e.Param())
	case "oneof":
		return fmt.Sprintf("%s harus salah satu dari: %s", e.Field(), e.Param())
	default:
		return fmt.Sprintf("%s tidak valid", e.Field())
	}
}