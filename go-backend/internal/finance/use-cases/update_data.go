package usecases

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/viteant/stockinsight/internal/finance/domain"
)

type UpdateFinanceDataUseCase struct {
	StockRepo   domain.StockRepository
	FinanceRepo domain.FinanceRepository
	Scraper     domain.FinanceScraper
}

func NewUpdateFinanceDataUseCase(
	stockRepo domain.StockRepository,
	financeRepo domain.FinanceRepository,
	scraper domain.FinanceScraper,
) *UpdateFinanceDataUseCase {
	return &UpdateFinanceDataUseCase{
		StockRepo:   stockRepo,
		FinanceRepo: financeRepo,
		Scraper:     scraper,
	}
}

func (u *UpdateFinanceDataUseCase) Execute() error {
	tickers, err := u.StockRepo.GetTickersDateRange()
	if err != nil {
		return err
	}

	throttle := 500 * time.Millisecond
	if v := os.Getenv("THROTTLE_MS"); v != "" {
		if ms, err := strconv.Atoi(v); err == nil {
			throttle = time.Duration(ms) * time.Millisecond
		}
	}

	for _, t := range tickers {
		adjustedStart := t.StartDate.Add(-1 * time.Hour)
		log.Printf("Scrapeando %s desde %s hasta %s", t.Ticker, adjustedStart.Format("2006-01-02"), t.EndDate.Format("2006-01-02"))

		data, err := u.Scraper.GetHistoricalData(t.Ticker, adjustedStart, t.EndDate)
		if err != nil {
			log.Printf("Error scrapeando %s: %v", t.Ticker, err)
		} else if len(data) == 0 {
			log.Printf("Sin datos para %s", t.Ticker)
		} else {
			if err := u.FinanceRepo.BulkSave(data); err != nil {
				log.Printf("Error guardando %s: %v", t.Ticker, err)
			} else {
				log.Printf("%d registros guardados para %s", len(data), t.Ticker)
			}
		}

		time.Sleep(throttle)
	}
	return nil
}
