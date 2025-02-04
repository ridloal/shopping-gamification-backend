package domain

type ClaimRequestInput struct {
	ProductID           int64  `json:"product_id" validate:"required"`
	SocialMediaUsername string `json:"social_media_username" validate:"required,min=3,max=50"`
	SocialMediaPlatform string `json:"social_media_platform" validate:"required,oneof=instagram tiktok youtube facebook"`
	PostURL             string `json:"post_url" default:""`
	NomorWhatsapp       string `json:"nomor_whatsapp" validate:"required,indonesian_phone"`
	Email               string `json:"email" validate:"required,email"`
	IsLiked             bool   `json:"is_liked" default:"false"`
	IsComment           bool   `json:"is_comment" default:"false"`
	IsShared            bool   `json:"is_shared" default:"false"`
	IsFollow            bool   `json:"is_follow" default:"false"`
}
