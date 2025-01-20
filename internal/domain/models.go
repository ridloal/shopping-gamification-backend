package domain

import "database/sql"

type Product struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	ImageURL    string  `json:"image_url"`
	Stock       int     `json:"stock"`
	Status      bool    `json:"status"`
}

type Prize struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	DiscountPercentage int    `json:"discount_percentage"`
	Quota              int    `json:"quota"`
	RemainingQuota     int    `json:"remaining_quota"`
	Status             bool   `json:"status"`
}

type PrizeGroup struct {
	ID          int64   `json:"id"`
	ProductID   int64   `json:"product_id"`
	PrizeID     int64   `json:"prize_id"`
	Probability float64 `json:"probability"`
	Status      bool    `json:"status"`
	Prize       Prize   `json:"prize"`
}

type ClaimRequest struct {
	ID                  int64          `json:"id"`
	UserID              *int64         `json:"user_id"`
	ProductID           int64          `json:"product_id"`
	PrizeID             *int64         `json:"prize_id"`
	SocialMediaUsername string         `json:"social_media_username"`
	SocialMediaPlatform string         `json:"social_media_platform"`
	PostURL             string         `json:"post_url"`
	VerificationStatus  string         `json:"verification_status"`
	ClaimCode           *string        `json:"claim_code"`
	NomorWhatsapp       string         `json:"nomor_whatsapp"`
	Email               string         `json:"email"`
	IsLiked             *bool          `json:"is_liked"`
	IsComment           *bool          `json:"is_comment"`
	IsShared            *bool          `json:"is_shared"`
	IsFollow            *bool          `json:"is_follow"`
	CreatedAt           string         `json:"created_at"`
	UpdatedAt           string         `json:"updated_at"`
	ClaimedAt           sql.NullString `json:"claimed_at"`
}

type ProductRepository interface {
	GetProducts() ([]Product, error)
	GetProductByID(productID int64) (Product, error)
	GetPrizeGroupsByProductID(productID int64) ([]PrizeGroup, error)
}

type ClaimRepository interface {
	CreateClaimRequest(req *ClaimRequestInput) (ClaimRequest, error)
	GetClaimRequestByID(claimID int64) (ClaimRequest, error)
	UpdateClaimRequestPrize(claimID int64, prizeID int64) error
}
