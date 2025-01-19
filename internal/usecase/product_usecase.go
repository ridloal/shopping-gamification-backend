package usecase

import "shopping-gamification/internal/domain"

type ProductUsecase interface {
	GetProducts() ([]domain.Product, error)
	GetProductByID(id int64) (domain.Product, error)
	GetPrizeGroupsByProductID(id int64) ([]domain.PrizeGroup, error)
}

type productUsecase struct {
	repo domain.ProductRepository
}

func NewProductUsecase(repo domain.ProductRepository) ProductUsecase {
	return &productUsecase{repo: repo}
}

func (u *productUsecase) GetProducts() ([]domain.Product, error) {
	return u.repo.GetProducts()
}

func (u *productUsecase) GetProductByID(id int64) (domain.Product, error) {
	return u.repo.GetProductByID(id)
}

func (u *productUsecase) GetPrizeGroupsByProductID(id int64) ([]domain.PrizeGroup, error) {
	return u.repo.GetPrizeGroupsByProductID(id)
}
