package usecase

import (
	"context"
	"encoding/json"
	"shopping-gamification/internal/domain"
	"shopping-gamification/internal/repository/redis"
)

const REDIS_KEY_PRODUCTS_LIST = "products_list"

type ProductUsecase interface {
	GetProducts(ctx context.Context) ([]domain.Product, error)
	GetProductByID(id int64) (domain.Product, error)
	GetPrizeGroupsByProductID(id int64) ([]domain.PrizeGroup, error)
}

type productUsecase struct {
	repo  domain.ProductRepository
	redis *redis.Repository
}

func NewProductUsecase(repo domain.ProductRepository, redis *redis.Repository) ProductUsecase {
	return &productUsecase{
		repo:  repo,
		redis: redis,
	}
}

func (u *productUsecase) GetProducts(ctx context.Context) ([]domain.Product, error) {

	listProduct := []domain.Product{}
	err := error(nil)

	// Check if redis is available
	if u.redis != nil {
		// Get products from redis
		value, err := u.redis.GetRedisValue(ctx, REDIS_KEY_PRODUCTS_LIST)

		if err == nil {
			err = json.Unmarshal([]byte(value), &listProduct)
			if err == nil {
				return listProduct, nil
			}
		}
	}

	// If redis is not available, use the repository
	listProduct, err = u.repo.GetProducts()
	if err != nil {
		// If there is an error, return the error
		return listProduct, err
	}

	if u.redis != nil {
		// Save products to redis
		data, err := json.Marshal(listProduct)
		if err == nil {
			err = u.redis.SetRedisValue(ctx, REDIS_KEY_PRODUCTS_LIST, string(data))
		}
	}

	return listProduct, err
}

func (u *productUsecase) GetProductByID(id int64) (domain.Product, error) {
	return u.repo.GetProductByID(id)
}

func (u *productUsecase) GetPrizeGroupsByProductID(id int64) ([]domain.PrizeGroup, error) {
	return u.repo.GetPrizeGroupsByProductID(id)
}
