package domain

import "time"

type TickerRange struct {
	Ticker    string
	StartDate time.Time
	EndDate   time.Time
}

type FinanceRepository interface {
	BulkSave(data []Finance) error
}

type StockRepository interface {
	GetTickersDateRange() ([]TickerRange, error)
}

type FinanceScraper interface {
	GetHistoricalData(ticker string, from, to time.Time) ([]Finance, error)
}
