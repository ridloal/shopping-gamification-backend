package usecase

import (
	"math/rand/v2"
	"shopping-gamification/internal/domain"
)

type ClaimUsecase interface {
	CreateClaimRequest(req *domain.ClaimRequestInput) (domain.ClaimRequest, error)
	GetClaimRequestByID(claimID int64) (domain.ClaimRequest, error)
	UpdateClaimRequestPrize(claimID int64, prizeID int64) error
	GetClaimRequestByClaimCode(claimCode string) (domain.ClaimRequest, error)
	ClaimPrize(claimCode string) (domain.PrizeResponse, error)
}

type claimUsecase struct {
	repo        domain.ClaimRepository
	repoProduct domain.ProductRepository
}

func NewClaimUsecase(repo domain.ClaimRepository, repoProduct domain.ProductRepository) ClaimUsecase {
	return &claimUsecase{
		repo:        repo,
		repoProduct: repoProduct,
	}
}

func (u *claimUsecase) CreateClaimRequest(req *domain.ClaimRequestInput) (domain.ClaimRequest, error) {
	return u.repo.CreateClaimRequest(req)
}

func (u *claimUsecase) GetClaimRequestByID(claimID int64) (domain.ClaimRequest, error) {
	return u.repo.GetClaimRequestByID(claimID)
}

func (u *claimUsecase) UpdateClaimRequestPrize(claimID int64, prizeID int64) error {
	return u.repo.UpdateClaimRequestPrize(claimID, prizeID, "{}")
}

func (u *claimUsecase) GetClaimRequestByClaimCode(claimCode string) (domain.ClaimRequest, error) {
	return u.repo.GetClaimRequestByClaimCode(claimCode)
}

func (u *claimUsecase) ClaimPrize(claimCode string) (domain.PrizeResponse, error) {
	claimRequest, err := u.repo.GetClaimRequestByClaimCode(claimCode)
	if err != nil {
		return domain.PrizeResponse{}, err
	}

	prizeGroup, err := u.repoProduct.GetPrizeGroupsByProductID(claimRequest.ProductID)
	if err != nil {
		return domain.PrizeResponse{}, err
	}

	prize := generateRandomPrizeID(prizeGroup)
	err = u.repo.UpdateClaimRequestPrize(claimRequest.ID, prize.PrizeID, prize.DetailJson)
	if err != nil {
		return domain.PrizeResponse{}, err
	}

	prizeResp := domain.PrizeResponse{
		PGID:       prize.ID,
		DetailJson: prize.DetailJson,
		PrizeName:  prize.Prize.Name,
		PrizeDesc:  prize.Prize.Description,
		ImageURL:   prize.Prize.ImageURL,
	}

	return prizeResp, nil
}

// function to generate random prize id by the probability
func generateRandomPrizeID(prizeGroups []domain.PrizeGroup) domain.PrizeGroup {
	var totalProbability float64
	for _, prizeGroup := range prizeGroups {
		totalProbability += prizeGroup.Probability
	}

	randomProbability := totalProbability * rand.Float64()
	for _, prizeGroup := range prizeGroups {
		randomProbability -= prizeGroup.Probability
		if randomProbability <= 0 {
			return prizeGroup
		}
	}

	return domain.PrizeGroup{}
}
