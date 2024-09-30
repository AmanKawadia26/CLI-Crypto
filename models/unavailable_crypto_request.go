//go:build !test
// +build !test

package models

import "time"

type UnavailableCryptoRequest struct {
	CryptoSymbol string `bson:"crypto_symbol" json:"crypto_symbol"`
	//CryptoName     string    `bson:"crypto_name" json:"crypto_name"`
	UserName       string    `bson:"user_name" json:"user_name"`
	RequestMessage string    `bson:"request_message" json:"request_message"`
	Status         string    `bson:"status" json:"status"`
	Timestamp      time.Time `bson:"timestamp" json:"timestamp"` // Using time.Time for timestamps
}
