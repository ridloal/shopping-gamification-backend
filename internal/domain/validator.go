package domain

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func NewValidator() *validator.Validate {
	validate := validator.New()

	// Custom validation untuk nomor whatsapp Indonesia
	validate.RegisterValidation("indonesian_phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		regex := regexp.MustCompile(`^(\+62|62|0)8[1-9][0-9]{6,9}$`)
		return regex.MatchString(phone)
	})

	return validate
}
