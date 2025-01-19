package usecase

import "shopping-gamification/internal/domain"

type ClaimUsecase interface {
	CreateClaimRequest(req *domain.ClaimRequest) error
	UpdateClaimRequestPrize(claimID int64, prizeID int64) error
}

type claimUsecase struct {
	repo domain.ClaimRepository
}

func NewClaimUsecase(repo domain.ClaimRepository) ClaimUsecase {
	return &claimUsecase{repo: repo}
}

func (u *claimUsecase) CreateClaimRequest(req *domain.ClaimRequest) error {
	return u.repo.CreateClaimRequest(req)
}

func (u *claimUsecase) UpdateClaimRequestPrize(claimID int64, prizeID int64) error {
	return u.repo.UpdateClaimRequestPrize(claimID, prizeID)
}
