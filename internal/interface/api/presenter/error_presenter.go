package presenter

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

// ErrorResponse creates an error response
func ErrorResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"error": map[string]interface{}{
			"message": message,
		},
	}
}

// ValidationErrorResponse creates a validation error response
func ValidationErrorResponse(err error) map[string]interface{} {
	// Check if the error is a validator error
	validationErrors, ok := err.(validator.ValidationErrors)
	if !ok {
		return ErrorResponse(err.Error())
	}

	// Create a map of field errors
	fieldErrors := make(map[string]string)
	for _, validationErr := range validationErrors {
		field := strings.ToLower(validationErr.Field())
		switch validationErr.Tag() {
		case "required":
			fieldErrors[field] = "This field is required"
		case "email":
			fieldErrors[field] = "Must be a valid email address"
		case "min":
			fieldErrors[field] = "Value must be at least " + validationErr.Param() + " characters long"
		case "max":
			fieldErrors[field] = "Value must be at most " + validationErr.Param() + " characters long"
		default:
			fieldErrors[field] = "Invalid value"
		}
	}

	return map[string]interface{}{
		"error": map[string]interface{}{
			"message": "Validation failed",
			"fields":  fieldErrors,
		},
	}
}
