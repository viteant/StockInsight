package use_cases

import (
	"log"

	"github.com/viteant/stockinsight/internal/stock/domain"
)

type StockFetcher interface {
	FetchPage(nextPage string) ([]domain.Stock, string, error)
}

type StockSaver interface {
	Save(stock domain.Stock) error
}

type SyncService struct {
	Fetcher StockFetcher
	Repo    StockSaver
}

func NewSyncService(fetcher StockFetcher, repo StockSaver) *SyncService {
	return &SyncService{
		Fetcher: fetcher,
		Repo:    repo,
	}
}

func (s *SyncService) Sync() error {
	next := ""

	for {
		stocks, nextPage, err := s.Fetcher.FetchPage(next)
		if err != nil {
			return err
		}

		for _, stock := range stocks {
			if err := s.Repo.Save(stock); err != nil {
				log.Printf("Error al guardar stock [%s]: %v", stock.Ticker, err)
			}
		}

		if nextPage == "" {
			break
		}
		next = nextPage
	}

	log.Println("Sincronización completada")
	return nil
}
