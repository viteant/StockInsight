package repository

import (
	"database/sql"
	"log"

	"github.com/viteant/stockinsight/internal/stock/domain"
)

type CockroachStockRepository struct {
	DB *sql.DB
}

func NewCockroachStockRepository(db *sql.DB) *CockroachStockRepository {
	return &CockroachStockRepository{DB: db}
}

func (r *CockroachStockRepository) Save(stock domain.Stock) error {
	_, err := r.DB.Exec(`
		INSERT INTO stocks (
			id, ticker, company, brokerage, action,
			rating_from, rating_to,
			normalize_rating_from, normalize_rating_to,
			target_from, target_to, created_at
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4,
			$5, $6,
			$7, $8,
			$9, $10, $11
		)
		ON CONFLICT (ticker, created_at) DO UPDATE SET
			company = excluded.company,
			brokerage = excluded.brokerage,
			action = excluded.action,
			rating_from = excluded.rating_from,
			rating_to = excluded.rating_to,
			normalize_rating_from = excluded.normalize_rating_from,
			normalize_rating_to = excluded.normalize_rating_to,
			target_from = excluded.target_from,
			target_to = excluded.target_to
	`,
		stock.Ticker,
		stock.Company,
		stock.Brokerage,
		stock.Action,
		stock.RatingFrom,
		stock.RatingTo,
		stock.NormalizedRatingFrom,
		stock.NormalizedRatingTo,
		stock.TargetFrom,
		stock.TargetTo,
		stock.ReportedAt,
	)

	if err != nil {
		log.Printf("Error saving stock %s: %v", stock.Ticker, err)
	} else {
		log.Printf("Stock saved or updated: %s", stock.Ticker)
	}

	return err
}
