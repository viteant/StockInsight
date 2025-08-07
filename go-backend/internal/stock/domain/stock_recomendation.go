package domain

type StockRecommendation struct {
	Ticker              string  `json:"ticker"`
	Company             string  `json:"company"`
	Brokerage           string  `json:"brokerage"`
	Action              string  `json:"action"`
	NormalizeRatingFrom string  `json:"normalize_rating_from"`
	NormalizeRatingTo   string  `json:"normalize_rating_to"`
	WeightScore         float64 `json:"weight_score"`
}
