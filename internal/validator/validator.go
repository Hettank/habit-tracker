package validator

import "github.com/go-playground/validator/v10"

var Validate = validator.New()

func FormatValidationErrors(err error) map[string]string {

	errors := make(map[string]string)

	validationErrors := err.(validator.ValidationErrors)

	for _, fieldErr := range validationErrors {

		errors[fieldErr.Field()] = fieldErr.Tag()
	}

	return errors
}
