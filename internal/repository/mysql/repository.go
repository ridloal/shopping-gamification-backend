package mysql

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
	query := `SELECT id, name, description, price, image_url, stock, status FROM products WHERE id = ?`
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
        WHERE pg.product_id = ? AND pg.status = true`

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

func (r *Repository) CreateClaimRequest(req *domain.ClaimRequest) error {
	query := `
        INSERT INTO claim_requests 
        (product_id, social_media_username, social_media_platform, post_url)
        VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(query, req.ProductID, req.SocialMediaUsername,
		req.SocialMediaPlatform, req.PostURL)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	req.ID = id
	return nil
}

func (r *Repository) UpdateClaimRequestPrize(claimID int64, prizeID int64) error {
	query := `UPDATE claim_requests SET prize_id = ? WHERE id = ?`
	_, err := r.db.Exec(query, prizeID, claimID)
	return err
}
