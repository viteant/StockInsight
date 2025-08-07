package stocks

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/viteant/stockinsight/internal/stock/domain"
)

func ExportStocksToJSON(db *sql.DB, filepath string) error {
	rows, err := db.Query(`
		SELECT ticker, company, brokerage, action,
		       rating_from, rating_to,
		       normalize_rating_from, normalize_rating_to,
		       target_from, target_to, created_at
		FROM stocks
	`)
	if err != nil {
		return err
	}
	defer rows.Close()

	var stocks []domain.Stock

	for rows.Next() {
		var s domain.Stock
		err := rows.Scan(
			&s.Ticker,
			&s.Company,
			&s.Brokerage,
			&s.Action,
			&s.RatingFrom,
			&s.RatingTo,
			&s.NormalizedRatingFrom,
			&s.NormalizedRatingTo,
			&s.TargetFrom,
			&s.TargetTo,
			&s.ReportedAt,
		)
		if err != nil {
			log.Printf("❌ Error escaneando fila: %v", err)
			continue
		}
		stocks = append(stocks, s)
	}

	data, err := json.MarshalIndent(stocks, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filepath, data, 0644)
	if err == nil {
		log.Printf("✅ Exportación completa: %d registros escritos en %s\n", len(stocks), filepath)
	}
	return err
}
