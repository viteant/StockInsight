package repository

import (
	"database/sql"
	"log"
	"time"

	"github.com/viteant/stockinsight/internal/finance/domain"
)

type CockroachStockRepository struct {
	DB *sql.DB
}

func NewCockroachStockRepository(db *sql.DB) *CockroachStockRepository {
	return &CockroachStockRepository{DB: db}
}

func (r *CockroachStockRepository) GetTickersDateRange() ([]domain.TickerRange, error) {
	rows, err := r.DB.Query(`
		SELECT ticker, MIN(created_at), MAX(created_at)
		FROM stocks
		GROUP BY ticker
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []domain.TickerRange
	for rows.Next() {
		var tr domain.TickerRange
		var minTime, maxTime time.Time

		if err := rows.Scan(&tr.Ticker, &minTime, &maxTime); err != nil {
			log.Printf("Error leyendo fila de rango: %v", err)
			continue
		}

		tr.StartDate = minTime
		tr.EndDate = maxTime
		result = append(result, tr)
	}

	return result, nil
}
