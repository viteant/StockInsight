package domain

import "time"

type Stock struct {
	ID                  string    `json:"id"`
	Ticker              string    `json:"ticker"`
	Company             string    `json:"company"`
	Brokerage           string    `json:"brokerage"`
	Action              string    `json:"action"`
	RatingFrom          string    `json:"rating_from"`
	RatingTo            string    `json:"rating_to"`
	NormalizeRatingFrom string    `json:"normalize_rating_from"`
	NormalizeRatingTo   string    `json:"normalize_rating_to"`
	TargetFrom          float32   `json:"target_from"`
	TargetTo            float32   `json:"target_to"`
	ReportedAt          time.Time `json:"created_at"`
}
