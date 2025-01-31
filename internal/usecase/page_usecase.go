package usecase

import (
	"context"
	"log"
	"shopping-gamification/internal/domain"
	"shopping-gamification/internal/repository/postgres"
	"shopping-gamification/internal/repository/redis"
)

type PageUsecase interface {
	GetPageHome(ctx context.Context) (domain.PageHome, error)
}

type pageUsecase struct {
	postgresRepo *postgres.Repository
	redisRepo    *redis.Repository
}

func NewPageUsecase(postgresRepo *postgres.Repository, redisRepo *redis.Repository) PageUsecase {
	return &pageUsecase{postgresRepo: postgresRepo, redisRepo: redisRepo}
}

func (u *pageUsecase) GetPageHome(ctx context.Context) (domain.PageHome, error) {
	var pageHome domain.PageHome
	var err error

	if u.redisRepo != nil {
		pageHome, err = u.redisRepo.GetPageHome(ctx)
		if err != nil {
			log.Println("Failed to get data from redis:", err)
			return pageHome, err
		}

		if len(pageHome.TopProducts) != 0 {
			return pageHome, nil
		}
	}

	pageHome, err = u.postgresRepo.GetPageHome()
	if err != nil {
		return pageHome, err
	}

	if u.redisRepo != nil {
		err = u.redisRepo.SetPageHome(ctx, pageHome)
		if err != nil {
			log.Println("Failed to set data to redis:", err)
		}
	}

	return pageHome, nil
}
