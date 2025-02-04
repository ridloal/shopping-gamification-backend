package domain

import "database/sql"

type Product struct {
	ID            int64   `json:"id"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	OriginalPrice float64 `json:"original_price"`
	Price         float64 `json:"price"`
	ImageURL      string  `json:"image_url"`
	Stock         int     `json:"stock"`
	Stars         float32 `json:"stars"`
	Sold          int     `json:"sold"`
	Review        int     `json:"review"`
	ExternalLink  string  `json:"external_link"`
	IsDigital     bool    `json:"is_digital"`
	Status        bool    `json:"status"`
}

type Prize struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name"`
	Description        string `json:"description"`
	DiscountPercentage int    `json:"discount_percentage"`
	Quota              int    `json:"quota"`
	RemainingQuota     int    `json:"remaining_quota"`
	Status             bool   `json:"status"`
	ImageURL           string `json:"image_url"`
}

type PrizeGroup struct {
	ID          int64   `json:"id"`
	ProductID   int64   `json:"product_id"`
	PrizeID     int64   `json:"prize_id"`
	Probability float64 `json:"probability"`
	Status      bool    `json:"status"`
	Prize       Prize   `json:"prize"`
	DetailJson  string  `json:"detail_json"`
}

type ClaimRequest struct {
	ID                  int64          `json:"id"`
	UserID              *int64         `json:"user_id"`
	ProductID           int64          `json:"product_id"`
	PrizeID             *int64         `json:"prize_id"`
	PrizeDetail         string         `json:"prize_detail"`
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

type SocialContent struct {
	ProductID   int64  `json:"product_id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Platform    string `json:"platform"`
	PostURL     string `json:"post_url"`
}

type PageHome struct {
	TopProducts     []Product       `json:"top_products"`
	DigitalProducts []Product       `json:"digital_products"`
	SocialContents  []SocialContent `json:"social_contents"`
	Prize           []Prize         `json:"prize"`
}

type ProductRepository interface {
	GetProducts() ([]Product, error)
	GetProductByID(productID int64) (Product, error)
	GetPrizeGroupsByProductID(productID int64) ([]PrizeGroup, error)
}

type ClaimRepository interface {
	CreateClaimRequest(req *ClaimRequestInput) (ClaimRequest, error)
	GetClaimRequestByID(claimID int64) (ClaimRequest, error)
	UpdateClaimRequestPrize(claimID int64, prizeID int64, prizeDetail string) error
	GetClaimRequestByClaimCode(claimCode string) (ClaimRequest, error)
}

type PageRepository interface {
	GetPageHome() (PageHome, error)
}
