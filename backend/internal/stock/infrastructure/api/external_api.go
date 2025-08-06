package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/viteant/stockinsight/internal/stock/domain"
)

type ExternalAPIResponse struct {
	Items    []ExternalAPIItem `json:"items"`
	NextPage string            `json:"next_page"`
}

type ExternalAPIItem struct {
	Ticker     string `json:"ticker"`
	Company    string `json:"company"`
	Brokerage  string `json:"brokerage"`
	Action     string `json:"action"`
	RatingFrom string `json:"rating_from"`
	RatingTo   string `json:"rating_to"`
	TargetFrom string `json:"target_from"`
	TargetTo   string `json:"target_to"`
	Time       string `json:"time"`
}

type ExternalAPIClient struct {
	Endpoint string
	Token    string
}

func NewExternalAPIClient() *ExternalAPIClient {
	return &ExternalAPIClient{
		Endpoint: os.Getenv("API_ENDPOINT"),
		Token:    os.Getenv("API_TOKEN"),
	}
}

func (c *ExternalAPIClient) FetchPage(nextPage string) ([]domain.Stock, string, error) {
	url := c.Endpoint
	if nextPage != "" {
		url += fmt.Sprintf("?next_page=%s", nextPage)
	}

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+c.Token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var apiResp ExternalAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		log.Println("Error al parsear respuesta:", err)
		return nil, "", err
	}

	var result []domain.Stock
	for _, item := range apiResp.Items {
		reportedAt, _ := time.Parse(time.RFC3339Nano, item.Time)

		var tf, tt float32
		fmt.Sscanf(strings.Replace(item.TargetFrom, "$", "", 1), "%f", &tf)
		fmt.Sscanf(strings.Replace(item.TargetTo, "$", "", 1), "%f", &tt)

		result = append(result, domain.Stock{
			Ticker:     item.Ticker,
			Company:    item.Company,
			Brokerage:  item.Brokerage,
			Action:     item.Action,
			RatingFrom: item.RatingFrom,
			RatingTo:   item.RatingTo,
			TargetFrom: tf,
			TargetTo:   tt,
			ReportedAt: reportedAt,
		})
	}

	return result, apiResp.NextPage, nil
}
