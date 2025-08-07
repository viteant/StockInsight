package stocks

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/viteant/stockinsight/internal/stock/domain"
)

func ImportStocksFromJSON(db *sql.DB, filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("no se pudo leer el archivo: %w", err)
	}

	var stocks []domain.Stock
	if err := json.Unmarshal(data, &stocks); err != nil {
		return fmt.Errorf("no se pudo parsear el JSON: %w", err)
	}

	stmt := `
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
	`

	for _, s := range stocks {
		_, err := db.Exec(stmt,
			s.Ticker,
			s.Company,
			s.Brokerage,
			s.Action,
			s.RatingFrom,
			s.RatingTo,
			s.NormalizeRatingFrom,
			s.NormalizeRatingTo,
			s.TargetFrom,
			s.TargetTo,
			s.ReportedAt,
		)
		if err != nil {
			fmt.Printf("Error insertando stock %s (%s): %v\n", s.Ticker, s.ReportedAt.Format(time.RFC3339), err)
		}
	}

	fmt.Printf("Importación completa: %d registros procesados\n", len(stocks))
	return nil
}
