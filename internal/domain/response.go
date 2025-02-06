package domain

type PrizeGroupResponse struct {
	ID        int64 `json:"id"`
	ProductID int64 `json:"product_id"`
	PrizeID   int64 `json:"prize_id"`
	Status    bool  `json:"status"`
	Prize     Prize `json:"prize"`
}

type PrizeResponse struct {
	PGID       int64  `json:"pg_id"`
	DetailJson string `json:"detail_json"`
	PrizeName  string `json:"prize_name"`
	PrizeDesc  string `json:"prize_desc"`
	ImageURL   string `json:"image_url"`
}
