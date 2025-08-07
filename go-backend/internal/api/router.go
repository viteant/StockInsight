package api

import (
	"database/sql"

	"github.com/gofiber/fiber/v2"
	stockroutes "github.com/viteant/stockinsight/internal/stock/interfaces"
)

func RegisterRoutes(app *fiber.App, db *sql.DB) {
	apiGroup := app.Group("/api")

	stockroutes.RegisterStockRoutes(apiGroup, db)

	// Luego podrías tener:
	// financeroutes.RegisterFinanceRoutes(apiGroup, db)
}
