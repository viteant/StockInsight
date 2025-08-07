package repository

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/viteant/stockinsight/internal/stock/domain"
)

type PersistenceStockRepository struct {
	DB *sql.DB
}

func NewCockroachStockRepository(db *sql.DB) *PersistenceStockRepository {
	return &PersistenceStockRepository{DB: db}
}

func (r *PersistenceStockRepository) Save(stock domain.Stock) error {
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
		stock.NormalizeRatingFrom,
		stock.NormalizeRatingTo,
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

func (r *PersistenceStockRepository) FetchRecommendations() ([]domain.StockRecommendation, error) {
	query := `
		(
			SELECT
				s.ticker,
				s.company,
				s.brokerage,
				s.action,
				s.normalize_rating_from,
				s.normalize_rating_to,
				b.weight_score
			FROM stocks s
			JOIN broker_evaluation b ON s.brokerage = b.brokerage
			WHERE s.normalize_rating_to = 'buy'
			ORDER BY b.weight_score DESC
			LIMIT 10
		)
		UNION ALL
		(
			SELECT
				s.ticker,
				s.company,
				s.brokerage,
				s.action,
				s.normalize_rating_from,
				s.normalize_rating_to,
				b.weight_score
			FROM stocks s
			JOIN broker_evaluation b ON s.brokerage = b.brokerage
			WHERE s.normalize_rating_to = 'hold'
			ORDER BY b.weight_score DESC
			LIMIT 10
		)
		UNION ALL
		(
			SELECT
				s.ticker,
				s.company,
				s.brokerage,
				s.action,
				s.normalize_rating_from,
				s.normalize_rating_to,
				b.weight_score
			FROM stocks s
			JOIN broker_evaluation b ON s.brokerage = b.brokerage
			WHERE s.normalize_rating_to = 'sell'
			ORDER BY b.weight_score DESC
			LIMIT 10
		);
	`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []domain.StockRecommendation
	for rows.Next() {
		var r domain.StockRecommendation
		if err := rows.Scan(
			&r.Ticker,
			&r.Company,
			&r.Brokerage,
			&r.Action,
			&r.NormalizeRatingFrom,
			&r.NormalizeRatingTo,
			&r.WeightScore,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}

	return results, nil
}

func (r *PersistenceStockRepository) FetchAllStocks(
	page, limit int,
	filters map[string]string,
	orderBy, orderDir string,
) ([]domain.Stock, int, error) {
	offset := (page - 1) * limit

	allowedOrderFields := map[string]bool{
		"created_at":  true,
		"target_to":   true,
		"target_from": true,
		"ticker":      true,
	}

	if !allowedOrderFields[orderBy] {
		orderBy = "created_at"
	}
	if strings.ToLower(orderDir) != "asc" {
		orderDir = "DESC"
	} else {
		orderDir = "ASC"
	}

	whereClauses := []string{}
	args := []interface{}{}
	argIndex := 1

	for _, field := range []string{"ticker", "company", "brokerage"} {
		if val, ok := filters[field]; ok && val != "" {
			whereClauses = append(whereClauses, fmt.Sprintf("%s ILIKE $%d", field, argIndex))
			args = append(args, "%"+val+"%")
			argIndex++
		}
	}

	for key, column := range map[string]string{
		"target_from_min": "target_from >= $%d",
		"target_from_max": "target_from <= $%d",
		"target_to_min":   "target_to >= $%d",
		"target_to_max":   "target_to <= $%d",
		"date_from":       "created_at >= $%d",
		"date_to":         "created_at <= $%d",
	} {
		if val, ok := filters[key]; ok && val != "" {
			whereClauses = append(whereClauses, fmt.Sprintf(column, argIndex))
			args = append(args, val)
			argIndex++
		}
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	query := fmt.Sprintf(`
		SELECT
			id, ticker, company, brokerage, action,
			rating_from, rating_to,
			normalize_rating_from, normalize_rating_to,
			target_from, target_to, created_at
		FROM stocks
		%s
		ORDER BY %s %s
		LIMIT $%d OFFSET $%d;
	`, whereSQL, orderBy, orderDir, argIndex, argIndex+1)

	args = append(args, limit, offset)

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var stocks []domain.Stock
	for rows.Next() {
		var s domain.Stock
		if err := rows.Scan(
			&s.ID,
			&s.Ticker,
			&s.Company,
			&s.Brokerage,
			&s.Action,
			&s.RatingFrom,
			&s.RatingTo,
			&s.NormalizeRatingFrom,
			&s.NormalizeRatingTo,
			&s.TargetFrom,
			&s.TargetTo,
			&s.ReportedAt,
		); err != nil {
			return nil, 0, err
		}
		stocks = append(stocks, s)
	}

	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM stocks %s`, whereSQL)
	countRow := r.DB.QueryRow(countQuery, args[:argIndex-1]...)

	var total int
	if err := countRow.Scan(&total); err != nil {
		return nil, 0, err
	}

	return stocks, total, nil
}
