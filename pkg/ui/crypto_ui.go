//go:build !test
// +build !test

package ui

import (
	"cryptotracker/pkg/utils"
	"fmt"
	"github.com/fatih/color"
	"strings"
)

func DisplayCryptoGraph(cryptoName string, currentPrice float64) {
	prices := utils.GenerateRandomPrices(30, currentPrice)

	color.New(color.FgCyan).Printf("30-day price graph for %s:\n\n", cryptoName)

	maxPrice := prices[0]
	minPrice := prices[0]
	for _, price := range prices {
		if price > maxPrice {
			maxPrice = price
		}
		if price < minPrice {
			minPrice = price
		}
	}

	graphHeight := 20
	for i := 0; i < graphHeight; i++ {
		price := maxPrice - (float64(i) * (maxPrice - minPrice) / float64(graphHeight-1))
		fmt.Printf("%8.2f |", price)

		for _, p := range prices {
			if p >= price {
				color.New(color.FgGreen).Print("█")
			} else {
				fmt.Print(" ")
			}
		}
		fmt.Println()
	}

	fmt.Print("         ")
	fmt.Println(strings.Repeat("-", len(prices)))

	fmt.Print("         ")
	for i := 0; i < len(prices); i++ {
		fmt.Print("─")
	}
	fmt.Println()
}
