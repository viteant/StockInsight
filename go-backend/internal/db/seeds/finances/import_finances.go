package finances

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/viteant/stockinsight/internal/finance/domain"
)

func ImportFinanceDataFromJSON(db *sql.DB, filepath string) error {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("no se pudo leer el archivo: %w", err)
	}

	var entries []domain.Finance
	if err := json.Unmarshal(data, &entries); err != nil {
		return fmt.Errorf("no se pudo parsear el JSON: %w", err)
	}

	stmt := `
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
    `

	for _, f := range entries {
		_, err := db.Exec(stmt,
			f.Ticker,
			f.Date,
			f.Open,
			f.High,
			f.Low,
			f.Close,
			f.Volume,
			f.Source,
			f.ScrapedAt,
		)
		if err != nil {
			log.Printf("Error insertando FinanceData %s [%s]: %v",
				f.Ticker, f.Date.Format("2006-01-02"), err)
		}
	}

	log.Printf("Importación completa FinanceData: %d registros procesados\n", len(entries))
	return nil
}
