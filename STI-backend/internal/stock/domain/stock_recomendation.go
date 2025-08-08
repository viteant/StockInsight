package domain

type StockRecommendation struct {
	ID                  string  `json:"id"`
	Ticker              string  `json:"ticker"`
	Company             string  `json:"company"`
	Brokerage           string  `json:"brokerage"`
	Action              string  `json:"action"`
	TargetFrom          float32 `json:"target_from"`
	TargetTo            float32 `json:"target_to"`
	NormalizeRatingFrom string  `json:"normalize_rating_from"`
	NormalizeRatingTo   string  `json:"normalize_rating_to"`
	WeightScore         float64 `json:"weight_score"`
}
