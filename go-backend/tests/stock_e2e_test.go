package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"

	"github.com/viteant/stockinsight/internal/api"
	"github.com/viteant/stockinsight/internal/db"
)

func setupE2EApp() *fiber.App {
	_ = godotenv.Load("../.env") // Asegúrate que DATABASE_URI esté cargado

	dbConn := db.NewCockroachDB()
	app := fiber.New()

	api.RegisterRoutes(app, dbConn)
	return app
}

func TestGetStocksE2E_WithDynamicFilters(t *testing.T) {
	app := setupE2EApp()

	// Paso 1: Consulta general sin filtros
	req1 := httptest.NewRequest("GET", "/api/stocks?limit=5", nil)
	req1.Header.Set("Content-Type", "application/json")

	resp1, err := app.Test(req1, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp1.StatusCode)

	var body1 map[string]interface{}
	err = json.NewDecoder(resp1.Body).Decode(&body1)
	assert.NoError(t, err)

	items1, ok := body1["items"].([]interface{})
	assert.True(t, ok)
	assert.Greater(t, len(items1), 0, "Se requiere al menos un registro para validar los filtros")

	// Paso 2: Extraer campos del primer item
	first := items1[0].(map[string]interface{})

	ticker := url.QueryEscape(first["ticker"].(string))
	company := url.QueryEscape(first["company"].(string))
	brokerage := url.QueryEscape(first["brokerage"].(string))
	targetFrom := first["target_from"].(float64)
	createdAt := first["created_at"].(string)

	// Paso 3: Construcción de la URL con escape
	filteredURL := fmt.Sprintf(
		"/api/stocks?limit=10&ticker=%s&company=%s&brokerage=%s&target_from_min=%.2f&target_from_max=%.2f&date_from=%s&date_to=%s",
		ticker, company, brokerage, targetFrom-1, targetFrom+1, createdAt, createdAt)

	t.Logf("🌐 URL generada: %s", filteredURL)

	// Paso 4: Ejecutar consulta filtrada
	req2 := httptest.NewRequest("GET", filteredURL, nil)
	resp2, err := app.Test(req2, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)

	var body2 map[string]interface{}
	err = json.NewDecoder(resp2.Body).Decode(&body2)
	assert.NoError(t, err)

	items2, ok := body2["items"].([]interface{})
	assert.True(t, ok)

	// Paso 5: Validar que los datos coincidan con los filtros
	for _, item := range items2 {
		s := item.(map[string]interface{})
		assert.Contains(t, s["ticker"], first["ticker"])
		assert.Contains(t, s["company"], first["company"])
		assert.Contains(t, s["brokerage"], first["brokerage"])

		if tf, ok := s["target_from"].(float64); ok {
			assert.GreaterOrEqual(t, tf, targetFrom-1)
			assert.LessOrEqual(t, tf, targetFrom+1)
		}

		if dateStr, ok := s["created_at"].(string); ok {
			assert.Equal(t, dateStr, createdAt)
		}
	}
}

func TestGetRecommendationsE2E(t *testing.T) {
	app := setupE2EApp()

	req := httptest.NewRequest("GET", "/api/recommendations", nil)
	resp, err := app.Test(req, -1)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var body []map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&body)
	assert.NoError(t, err)

	// Validamos que haya exactamente 30 recomendaciones
	assert.Len(t, body, 30, "deben haber exactamente 30 recomendaciones")

	// Validamos que haya 10 por cada categoría (buy, hold, sell)
	counts := map[string]int{}
	for _, rec := range body {
		assert.NotEmpty(t, rec["ticker"])
		assert.NotEmpty(t, rec["brokerage"])

		// Validamos y contamos cada categoría
		rating, ok := rec["normalize_rating_to"].(string)
		assert.True(t, ok, "normalize_rating_to debe ser string")
		counts[rating]++
	}

	assert.Equal(t, 10, counts["buy"], "deben haber 10 recomendaciones 'buy'")
	assert.Equal(t, 10, counts["hold"], "deben haber 10 recomendaciones 'hold'")
	assert.Equal(t, 10, counts["sell"], "deben haber 10 recomendaciones 'sell'")
}
