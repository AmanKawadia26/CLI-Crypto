//go:build !test
// +build !test

package models

type Cryptocurrency struct {
	ID        string `json:"id"`
	Symbol    string `json:"symbol"`
	Name      string `json:"name"`
	PriceUSD  string `json:"price_usd"`
	PriceBTC  string `json:"price_btc"`
	MarketCap string `json:"market_cap_usd"`
	Change24h string `json:"percent_change_24h"`
}
