// Package validation for validate data
package validation

import (
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

// CustomValidator struct for validator
type CustomValidator struct {
	Validator *validator.Validate
}

// Validate for validating data
func (cv *CustomValidator) Validate(i interface{}) error {
	if err := cv.Validator.Struct(i); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}
