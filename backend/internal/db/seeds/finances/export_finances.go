package finances

import (
	"database/sql"
	"encoding/json"
	"log"
	"os"

	"github.com/viteant/stockinsight/internal/finance/domain"
)

func ExportFinanceDataToJSON(db *sql.DB, filepath string) error {
	rows, err := db.Query(`
        SELECT ticker, date, open, high, low, close, volume, source, scraped_at
        FROM finances
    `)
	if err != nil {
		return err
	}
	defer rows.Close()

	var entries []domain.Finance
	for rows.Next() {
		var f domain.Finance
		if err := rows.Scan(
			&f.Ticker,
			&f.Date,
			&f.Open,
			&f.High,
			&f.Low,
			&f.Close,
			&f.Volume,
			&f.Source,
			&f.ScrapedAt,
		); err != nil {
			log.Printf("❌ Error escaneando fila finance: %v", err)
			continue
		}
		entries = append(entries, f)
	}

	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filepath, data, 0644); err != nil {
		return err
	}

	log.Printf("✅ Exportación completa de FinanceData: %d registros escritos en %s\n", len(entries), filepath)
	return nil
}
