package validator

import (
	"github.com/go-playground/validator/v10"
)

// CustomValidator implements echo.Validator
type CustomValidator struct {
	validator *validator.Validate
}

// NewCustomValidator creates a new custom validator
func NewCustomValidator() *CustomValidator {
	return &CustomValidator{
		validator: validator.New(),
	}
}

// Validate validates the given struct
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}
