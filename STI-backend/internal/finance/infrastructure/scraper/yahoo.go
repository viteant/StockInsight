package scraper

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/viteant/stockinsight/internal/finance/domain"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
	"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
	"Mozilla/5.0 (X11; Linux x86_64)",
	// Puedes añadir más variantes
}

type YahooFinanceScraper struct{}

func NewYahooFinanceScraper() *YahooFinanceScraper {
	return &YahooFinanceScraper{}
}

func getRandomUserAgent() string {
	return userAgents[rand.Intn(len(userAgents))]
}

type yahooResponse struct {
	Chart struct {
		Result []struct {
			Timestamp  []int64 `json:"timestamp"`
			Indicators struct {
				Quote []struct {
					Open   []float64 `json:"open"`
					High   []float64 `json:"high"`
					Low    []float64 `json:"low"`
					Close  []float64 `json:"close"`
					Volume []int64   `json:"volume"`
				} `json:"quote"`
			} `json:"indicators"`
		} `json:"result"`
		Error any `json:"error"`
	} `json:"chart"`
}

func (s *YahooFinanceScraper) GetHistoricalData(
	ticker string, from, to time.Time,
) ([]domain.Finance, error) {
	url := fmt.Sprintf(
		"https://query2.finance.yahoo.com/v8/finance/chart/%s?period1=%d&period2=%d&interval=1d&events=history&includeAdjustedClose=true",
		ticker, from.Unix(), to.Unix(),
	)
	log.Printf("🌐 Consultando URL para %s: %s", ticker, url)

	const maxRetries = 3
	delay := 1 * time.Second

	for attempt := 0; attempt <= maxRetries; attempt++ {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", getRandomUserAgent())

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error HTTP al solicitar %s: %w", ticker, err)
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests {
			if attempt == maxRetries {
				return nil, fmt.Errorf("recibido 429 varias veces para %s", ticker)
			}
			log.Printf("429 para %s (intento %d/%d), esperando %v...", ticker, attempt+1, maxRetries, delay)
			time.Sleep(delay)
			delay *= 2
			continue
		}
		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("HTTP %d para %s: %s", resp.StatusCode, ticker, string(bodyBytes))
		}
		if len(bodyBytes) == 0 || (bodyBytes[0] != '{' && bodyBytes[0] != '[') {
			return nil, fmt.Errorf("respuesta inesperada para %s: %s", ticker, string(bodyBytes))
		}

		var yr yahooResponse
		if err := json.Unmarshal(bodyBytes, &yr); err != nil {
			return nil, fmt.Errorf("error parseando JSON para %s: %w", ticker, err)
		}
		if len(yr.Chart.Result) == 0 || len(yr.Chart.Result[0].Timestamp) == 0 {
			return nil, fmt.Errorf("sin datos utiles para %s", ticker)
		}

		quote := yr.Chart.Result[0].Indicators.Quote[0]
		var result []domain.Finance
		for i, ts := range yr.Chart.Result[0].Timestamp {
			if i >= len(quote.Open) {
				break
			}
			result = append(result, domain.Finance{
				Ticker:    strings.ToUpper(ticker),
				Date:      time.Unix(ts, 0).UTC().Truncate(24 * time.Hour),
				Open:      float32(quote.Open[i]),
				High:      float32(quote.High[i]),
				Low:       float32(quote.Low[i]),
				Close:     float32(quote.Close[i]),
				Volume:    quote.Volume[i],
				Source:    "Yahoo",
				ScrapedAt: time.Now(),
			})
		}
		return result, nil
	}

	return nil, fmt.Errorf("falló scrapeo de %s tras %d intentos", ticker, maxRetries)
}
