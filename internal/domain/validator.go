package domain

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

type ClaimRequestInput struct {
	ProductID           int64  `json:"product_id" validate:"required"`
	SocialMediaUsername string `json:"social_media_username" validate:"required,min=3,max=50"`
	SocialMediaPlatform string `json:"social_media_platform" validate:"required,oneof=instagram tiktok"`
	PostURL             string `json:"post_url" validate:"required,url"`
	NomorWhatsapp       string `json:"nomor_whatsapp" validate:"required,e164"`
	Email               string `json:"email" validate:"required,email"`
}

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
