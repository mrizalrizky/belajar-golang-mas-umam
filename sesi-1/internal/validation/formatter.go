package validation

import "github.com/go-playground/validator/v10"

func FormatError(err error) map[string]string {
	errors := err.(validator.ValidationErrors)

	errorMap := make(map[string]string)
	for _, e := range errors {
		errorMap[e.Field()] = e.Tag()
	}
	return errorMap
}