package interfaces

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	"github.com/viteant/stockinsight/internal/stock/infrastructure/repository"
	"github.com/viteant/stockinsight/internal/stock/use_cases"
)

func RegisterStockRoutes(app fiber.Router, db *sql.DB) {
	stockRepo := repository.NewCockroachStockRepository(db)
	stockService := &use_cases.StockService{Repo: stockRepo}
	stockHandler := NewStockHandler(stockService)

	app.Get("/stocks", stockHandler.GetStocks)
	app.Get("/recommendations", stockHandler.GetRecommendations)
}
