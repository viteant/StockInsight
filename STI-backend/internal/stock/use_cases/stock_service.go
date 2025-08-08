package use_cases

import (
	"github.com/viteant/stockinsight/internal/stock/domain"
)

type StockRepository interface {
	FetchAllStocks(page, limit int, filters map[string]string, orderBy, orderDir string) ([]domain.Stock, int, error)
	FetchRecommendations() ([]domain.StockRecommendation, error)
}

type StockService struct {
	Repo StockRepository
}

func (s *StockService) GetRecommendations() ([]domain.StockRecommendation, error) {
	return s.Repo.FetchRecommendations()
}

func (s *StockService) GetAllStocks(page, limit int, filters map[string]string, orderBy, orderDir string) ([]domain.Stock, int, error) {
	return s.Repo.FetchAllStocks(page, limit, filters, orderBy, orderDir)
}
