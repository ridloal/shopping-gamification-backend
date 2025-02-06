package postgres

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"math/big"
	"shopping-gamification/internal/domain"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetProducts() ([]domain.Product, error) {
	query := `SELECT id, name, description, price, image_url, stock, status, original_price, stars, sold, reviews, external_link, is_digital FROM products WHERE status = true ORDER BY ID ASC LIMIT 50`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Price, &p.ImageURL, &p.Stock, &p.Status, &p.OriginalPrice, &p.Stars, &p.Sold, &p.Review, &p.ExternalLink, &p.IsDigital)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *Repository) GetProductByID(productID int64) (domain.Product, error) {
	query := `SELECT id, name, description, price, image_url, stock, status, original_price, stars, sold, reviews, external_link, is_digital FROM products WHERE id = $1`
	rows, err := r.db.Query(query, productID)
	if err != nil {
		return domain.Product{}, err
	}
	defer rows.Close()

	var product domain.Product
	for rows.Next() {
		err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.ImageURL, &product.Stock, &product.Status, &product.OriginalPrice, &product.Stars, &product.Sold, &product.Review, &product.ExternalLink, &product.IsDigital)
		if err != nil {
			return domain.Product{}, err
		}
	}
	return product, nil
}

func (r *Repository) GetPrizeGroupsByProductID(productID int64) ([]domain.PrizeGroup, error) {
	query := `
        SELECT pg.id, pg.product_id, pg.prize_id, pg.probability, pg.status,
               p.name, p.description, p.discount_percentage, p.quota, p.remaining_quota, p.status, pg.detail_json, p.image_url
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
			&pg.Prize.Quota, &pg.Prize.RemainingQuota, &pg.Prize.Status, &pg.DetailJson, &pg.Prize.ImageURL,
		)
		if err != nil {
			return nil, err
		}

		pg.Prize.ID = pg.PrizeID

		groups = append(groups, pg)
	}
	return groups, nil
}

func (r *Repository) CreateClaimRequest(req *domain.ClaimRequestInput) (domain.ClaimRequest, error) {

	claimCode := generateClaimCode()
	var claimRequest domain.ClaimRequest

	maxRetries := 3
	retryCount := 0

	for {
		query := `
        INSERT INTO claim_requests 
        (product_id, social_media_username, social_media_platform, post_url, nomor_whatsapp, email, is_liked, is_comment, is_shared, is_follow, claim_code)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id, product_id, social_media_username, social_media_platform, post_url, nomor_whatsapp, email, is_liked, is_comment, is_shared, is_follow, claim_code, verification_status`

		err := r.db.QueryRow(query, req.ProductID, req.SocialMediaUsername,
			req.SocialMediaPlatform, req.PostURL, req.NomorWhatsapp, req.Email, req.IsLiked, req.IsComment, req.IsShared, req.IsFollow, claimCode).Scan(
			&claimRequest.ID, &claimRequest.ProductID, &claimRequest.SocialMediaUsername,
			&claimRequest.SocialMediaPlatform, &claimRequest.PostURL, &claimRequest.NomorWhatsapp, &claimRequest.Email, &claimRequest.IsLiked, &claimRequest.IsComment, &claimRequest.IsShared, &claimRequest.IsFollow, &claimRequest.ClaimCode, &claimRequest.VerificationStatus)
		if err != nil {
			if retryCount < maxRetries {
				// Generate a new claim code and retry
				claimCode = generateClaimCode()
				retryCount++
				continue
			}
			return domain.ClaimRequest{}, err
		}
		break
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

func (r *Repository) GetClaimRequestByClaimCode(claimCode string) (domain.ClaimRequest, error) {
	query := `SELECT id, product_id, social_media_username, social_media_platform, post_url, nomor_whatsapp, email, verification_status, claim_code, is_liked, is_comment, is_shared, is_follow, prize_id, user_id, created_at, updated_at, claimed_at, prize_detail FROM claim_requests WHERE claim_code = $1`
	rows, err := r.db.Query(query, claimCode)

	if err != nil {
		return domain.ClaimRequest{}, err
	}
	defer rows.Close()

	var claimReq domain.ClaimRequest
	for rows.Next() {
		err := rows.Scan(&claimReq.ID, &claimReq.ProductID, &claimReq.SocialMediaUsername, &claimReq.SocialMediaPlatform, &claimReq.PostURL, &claimReq.NomorWhatsapp, &claimReq.Email, &claimReq.VerificationStatus, &claimReq.ClaimCode, &claimReq.IsLiked, &claimReq.IsComment, &claimReq.IsShared, &claimReq.IsFollow, &claimReq.PrizeID, &claimReq.UserID, &claimReq.CreatedAt, &claimReq.UpdatedAt, &claimReq.ClaimedAt, &claimReq.PrizeDetail)
		if err != nil {
			return domain.ClaimRequest{}, err
		}
	}
	return claimReq, nil
}

func (r *Repository) UpdateClaimRequestPrize(claimID int64, prizeID int64, prizeDetail string) error {
	// Validate prizeDetail as a JSON string
	var js json.RawMessage
	if err := json.Unmarshal([]byte(prizeDetail), &js); err != nil {
		prizeDetail = "{}"
	}

	query := `UPDATE claim_requests SET prize_id = $1, prize_detail = $2, verification_status = 'claimed', claimed_at = now() WHERE id = $3 AND verification_status = 'verified'`
	val, err := r.db.Exec(query, prizeID, prizeDetail, claimID)

	if err != nil {
		return err
	}

	rowsAffected, err := val.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("update prize is not eligible for this claim request, verification status must be verified")
	}

	return err
}

func (r *Repository) GetPageHome() (domain.PageHome, error) {
	type result struct {
		topProducts     []domain.Product
		digitalProducts []domain.Product
		socialContents  []domain.SocialContent
		prizes          []domain.Prize
		err             error
	}

	results := make(chan result, 3)

	// Get Top Product and Digital Product
	go func() {
		productList, err := r.GetProducts()
		if err != nil {
			results <- result{err: err}
			return
		}

		topProducts := []domain.Product{}
		digitalProducts := []domain.Product{}
		for i, p := range productList {
			if p.IsDigital {
				digitalProducts = append(digitalProducts, productList[i])
			} else {
				topProducts = append(topProducts, productList[i])
			}
		}

		results <- result{topProducts: topProducts, digitalProducts: digitalProducts}
	}()

	// Get social contents
	go func() {
		socialContentsQuery := `SELECT product_id, title, description, platform, post_url FROM social_contents ORDER BY product_id ASC LIMIT 10`
		socialContentsRows, err := r.db.Query(socialContentsQuery)
		if err != nil {
			results <- result{err: err}
			return
		}
		defer socialContentsRows.Close()

		var socialContents []domain.SocialContent
		for socialContentsRows.Next() {
			var sc domain.SocialContent
			err := socialContentsRows.Scan(&sc.ProductID, &sc.Title, &sc.Description, &sc.Platform, &sc.PostURL)
			if err != nil {
				results <- result{err: err}
				return
			}
			socialContents = append(socialContents, sc)
		}

		results <- result{socialContents: socialContents}
	}()

	// Get prizes
	go func() {
		prizesQuery := `SELECT id, name, description, discount_percentage, quota, remaining_quota, status, image_url FROM prizes WHERE status = true ORDER BY ID ASC LIMIT 50`
		prizesRows, err := r.db.Query(prizesQuery)
		if err != nil {
			results <- result{err: err}
			return
		}
		defer prizesRows.Close()

		var prizes []domain.Prize
		for prizesRows.Next() {
			var p domain.Prize
			err := prizesRows.Scan(&p.ID, &p.Name, &p.Description, &p.DiscountPercentage, &p.Quota, &p.RemainingQuota, &p.Status, &p.ImageURL)
			if err != nil {
				results <- result{err: err}
				return
			}
			prizes = append(prizes, p)
		}

		results <- result{prizes: prizes}
	}()

	var pageHome domain.PageHome
	for i := 0; i < 3; i++ {
		res := <-results
		if res.err != nil {
			return domain.PageHome{}, res.err
		}
		if res.topProducts != nil {
			pageHome.TopProducts = res.topProducts
			pageHome.DigitalProducts = res.digitalProducts
		}
		if res.socialContents != nil {
			pageHome.SocialContents = res.socialContents
		}
		if res.prizes != nil {
			pageHome.Prize = res.prizes
		}
	}

	return pageHome, nil
}

func generateClaimCode() string {
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	codeLength := 8
	code := make([]byte, codeLength)

	for i := range code {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		code[i] = charset[num.Int64()]
	}

	return string(code)
}
