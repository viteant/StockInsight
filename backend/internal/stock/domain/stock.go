package domain

import "time"

type Stock struct {
	Ticker     string
	Company    string
	Brokerage  string
	Action     string
	RatingFrom string
	RatingTo   string
	TargetFrom float32
	TargetTo   float32
	ReportedAt time.Time
}
