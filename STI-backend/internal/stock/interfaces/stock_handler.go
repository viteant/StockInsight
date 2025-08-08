package interfaces

import (
	"math"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/viteant/stockinsight/internal/stock/use_cases"
)

type StockHandler struct {
	useCase *use_cases.StockService
}

func NewStockHandler(useCase *use_cases.StockService) *StockHandler {
	return &StockHandler{
		useCase: useCase,
	}
}

// GetRecommendations godoc
// @Summary Recomendaciones de acciones
// @Description Devuelve una lista con 10 acciones recomendadas para comprar, mantener y vender, basadas en la puntuación de los brokers.
// @Tags Recommendations
// @Accept json
// @Produce json
// @Success 200 {array} domain.StockRecommendation
// @Failure 500 {object} map[string]string
// @Router /api/recommendations [get]
func (h *StockHandler) GetRecommendations(c *fiber.Ctx) error {
	recs, err := h.useCase.GetRecommendations()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching recommendations",
		})
	}
	return c.JSON(recs)
}

// GetStocks godoc
// @Summary Lista de acciones
// @Description Devuelve acciones con paginación y filtros
// @Tags Stocks
// @Accept json
// @Produce json
// @Param page query int false "Número de página"
// @Param limit query int false "Cantidad por página"
// @Param orderBy query string false "Columna para ordenar (default: ID)"
// @Param orderDir query string false "Dirección de orden (asc o desc, default": asc)"
// @Param ticker query string false "Filtra por ticker (ILIKE)"
// @Param company query string false "Filtra por nombre de empresa (ILIKE)"
// @Param brokerage query string false "Filtra por brokerage (ILIKE)"
// @Param target_from_min query number false "Filtra por target_from mínimo"
// @Param target_from_max query number false "Filtra por target_from máximo"
// @Param date_from query string false "Fecha mínima (YYYY-MM-DD)"
// @Param date_to query string false "Fecha máxima (YYYY-MM-DD)"
// @Success 200 {object} map[string]interface{}
// @Router /api/stocks [get]
func (h *StockHandler) GetStocks(c *fiber.Ctx) error {
	// Defaults
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))
	orderBy := c.Query("orderBy", "ID")
	orderDir := c.Query("orderDir", "asc")

	// Filtros soportados
	filters := map[string]string{
		"id":              c.Query("id"),
		"ticker":          c.Query("ticker"),
		"company":         c.Query("company"),
		"brokerage":       c.Query("brokerage"),
		"target_from_min": c.Query("target_from_min"),
		"target_from_max": c.Query("target_from_max"),
		"target_to_min":   c.Query("target_to_min"),
		"target_to_max":   c.Query("target_to_max"),
		"date_from":       c.Query("date_from"),
		"date_to":         c.Query("date_to"),
	}

	// Llama al caso de uso
	stocks, total, err := h.useCase.GetAllStocks(page, limit, filters, orderBy, orderDir)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Error fetching stocks",
			"message": err.Error(),
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return c.JSON(fiber.Map{
		"page":        page,
		"limit":       limit,
		"total":       total,
		"total_pages": totalPages,
		"items":       stocks,
	})
}
