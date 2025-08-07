package repository

import (
	"database/sql"
	"log"

	"github.com/viteant/stockinsight/internal/finance/domain"
)

type CockroachFinanceRepository struct {
	DB *sql.DB
}

func NewCockroachFinanceRepository(db *sql.DB) *CockroachFinanceRepository {
	return &CockroachFinanceRepository{DB: db}
}

func (r *CockroachFinanceRepository) BulkSave(data []domain.Finance) error {
	tx, err := r.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO finances (
			id, ticker, date, open, high, low, close, volume, source, scraped_at
		) VALUES (
			gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7, $8, $9
		)
		ON CONFLICT (ticker, date) DO UPDATE SET
			open = excluded.open,
			high = excluded.high,
			low = excluded.low,
			close = excluded.close,
			volume = excluded.volume,
			source = excluded.source,
			scraped_at = excluded.scraped_at
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, d := range data {
		_, err := stmt.Exec(
			d.Ticker, d.Date, d.Open, d.High, d.Low, d.Close, d.Volume, d.Source, d.ScrapedAt,
		)
		if err != nil {
			log.Printf("Error insertando %s [%s]: %v", d.Ticker, d.Date.Format("2006-01-02"), err)
			continue
		}
	}

	return tx.Commit()
}
