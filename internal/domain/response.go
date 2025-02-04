package domain

type PrizeGroupResponse struct {
	ID        int64 `json:"id"`
	ProductID int64 `json:"product_id"`
	PrizeID   int64 `json:"prize_id"`
	Status    bool  `json:"status"`
	Prize     Prize `json:"prize"`
}
