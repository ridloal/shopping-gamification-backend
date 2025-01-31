package usecase

import "shopping-gamification/internal/domain"

type PageUsecase interface {
	GetPageHome() (domain.PageHome, error)
}

type pageUsecase struct {
	repo domain.PageRepository
}

func NewPageUsecase(repo domain.PageRepository) PageUsecase {
	return &pageUsecase{repo: repo}
}

func (u *pageUsecase) GetPageHome() (domain.PageHome, error) {
	return u.repo.GetPageHome()
}
