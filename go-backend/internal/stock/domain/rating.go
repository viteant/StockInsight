package domain

import "strings"

func NormalizeBrokerRating(raw string) string {
	r := strings.ToLower(strings.TrimSpace(raw))

	switch r {
	case "buy", "strong-buy", "outperform", "outperformer",
		"market outperform", "mkt outperform", "overweight",
		"positive", "sector outperform", "speculative buy", "moderate buy":
		return "buy"

	case "hold", "neutral", "equal weight", "in-line",
		"market perform", "peer perform", "sector perform", "sector weight", "":
		return "hold"

	case "sell", "underweight", "underperform", "underperformer",
		"sector underperform", "reduce", "negative":
		return "sell"

	default:
		return "hold"
	}
}
