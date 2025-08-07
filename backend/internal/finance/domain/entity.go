package domain

import "time"

type Finance struct {
	Ticker    string
	Date      time.Time
	Open      float32
	High      float32
	Low       float32
	Close     float32
	Volume    int64
	Source    string
	ScrapedAt time.Time
}
