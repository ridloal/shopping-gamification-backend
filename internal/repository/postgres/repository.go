package postgres

import (
	"database/sql"
	"shopping-gamification/internal/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetProducts() ([]domain.Product, error) {
	query := `SELECT id, name, description, price, image_url, stock, status FROM products WHERE status = true`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.Stock, &p.Status)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *Repository) GetProductByID(productID int64) (domain.Product, error) {
	query := `SELECT id, name, description, price, image_url, stock, status FROM products WHERE id = $1`
	rows, err := r.db.Query(query, productID)
	if err != nil {
		return domain.Product{}, err
	}
	defer rows.Close()

	var product domain.Product
	for rows.Next() {
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.ImageURL, &product.Stock, &product.Status)
		if err != nil {
			return domain.Product{}, err
		}
	}
	return product, nil
}

func (r *Repository) GetPrizeGroupsByProductID(productID int64) ([]domain.PrizeGroup, error) {
	query := `
        SELECT pg.id, pg.product_id, pg.prize_id, pg.probability, pg.status,
               p.name, p.description, p.discount_percentage, p.quota, p.remaining_quota, p.status
        FROM prize_groups pg
        JOIN prizes p ON pg.prize_id = p.id
        WHERE pg.product_id = $1 AND pg.status = true`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var groups []domain.PrizeGroup
	for rows.Next() {
		var pg domain.PrizeGroup
		err := rows.Scan(
			&pg.ID, &pg.ProductID, &pg.PrizeID, &pg.Probability, &pg.Status,
			&pg.Prize.Name, &pg.Prize.Description, &pg.Prize.DiscountPercentage,
			&pg.Prize.Quota, &pg.Prize.RemainingQuota, &pg.Prize.Status,
		)
		if err != nil {
			return nil, err
		}
		groups = append(groups, pg)
	}
	return groups, nil
}

func (r *Repository) CreateClaimRequest(req *domain.ClaimRequestInput) (domain.ClaimRequest, error) {
	query := `
		INSERT INTO claim_requests 
		(product_id, social_media_username, social_media_platform, post_url, nomor_whatsapp, email)
		VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, product_id, social_media_username, social_media_platform, post_url, nomor_whatsapp, email`

	var claimRequest domain.ClaimRequest
	err := r.db.QueryRow(query, req.ProductID, req.SocialMediaUsername,
		req.SocialMediaPlatform, req.PostURL, req.NomorWhatsapp, req.Email).Scan(
		&claimRequest.ID, &claimRequest.ProductID, &claimRequest.SocialMediaUsername,
		&claimRequest.SocialMediaPlatform, &claimRequest.PostURL, &claimRequest.NomorWhatsapp, &claimRequest.Email)
	if err != nil {
		return domain.ClaimRequest{}, err
	}

	return claimRequest, nil
}

func (r *Repository) GetClaimRequestByID(claimID int64) (domain.ClaimRequest, error) {
	query := `SELECT id, product_id, social_media_username, social_media_platform, post_url, nomor_whatsapp, email, verification_status, claim_code, is_liked, is_comment, is_shared, is_follow, prize_id, user_id, created_at, updated_at, claimed_at FROM claim_requests WHERE id = $1`
	rows, err := r.db.Query(query, claimID)
	if err != nil {
		return domain.ClaimRequest{}, err
	}
	defer rows.Close()

	var claimReq domain.ClaimRequest
	for rows.Next() {
		err := rows.Scan(&claimReq.ID, &claimReq.ProductID, &claimReq.SocialMediaUsername, &claimReq.SocialMediaPlatform, &claimReq.PostURL, &claimReq.NomorWhatsapp, &claimReq.Email, &claimReq.VerificationStatus, &claimReq.ClaimCode, &claimReq.IsLiked, &claimReq.IsComment, &claimReq.IsShared, &claimReq.IsFollow, &claimReq.PrizeID, &claimReq.UserID, &claimReq.CreatedAt, &claimReq.UpdatedAt, &claimReq.ClaimedAt)
		if err != nil {
			return domain.ClaimRequest{}, err
		}
	}
	return claimReq, nil
}

func (r *Repository) UpdateClaimRequestPrize(claimID int64, prizeID int64) error {
	query := `UPDATE claim_requests SET prize_id = $1 WHERE id = $2`
	_, err := r.db.Exec(query, prizeID, claimID)
	return err
}
