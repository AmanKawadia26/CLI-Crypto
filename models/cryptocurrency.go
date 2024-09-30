//go:build !test
// +build !test

package models

// QuoteUSD represents the USD quote details for a cryptocurrency
type QuoteUSD struct {
	FullyDilutedMarketCap float64 `json:"fully_diluted_market_cap"`
	PercentChange1H       float64 `json:"percent_change_1h"`
	PercentChange24H      float64 `json:"percent_change_24h"`
	PercentChange30D      float64 `json:"percent_change_30d"`
	PercentChange60D      float64 `json:"percent_change_60d"`
	PercentChange7D       float64 `json:"percent_change_7d"`
	PercentChange90D      float64 `json:"percent_change_90d"`
	Price                 float64 `json:"price"`
}

// Quote represents the quote for a cryptocurrency
type Quote struct {
	USD QuoteUSD `json:"USD"`
}

// Cryptocurrency represents the basic details of a cryptocurrency
type Cryptocurrency struct {
	CMCRank     int    `json:"cmc_rank"`
	DateAdded   string `json:"date_added"`
	ID          int    `json:"id"`
	LastUpdated string `json:"last_updated"`
	Name        string `json:"name"`
	Quote       Quote  `json:"quote"`
	Slug        string `json:"slug"`
	Symbol      string `json:"symbol"`
}
