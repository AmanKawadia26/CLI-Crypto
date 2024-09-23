//go:build !test
// +build !test

package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PriceNotification struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	CryptoID    int                `bson:"crypto_id" json:"crypto_id"`
	Crypto      string             `bson:"crypto" json:"crypto"`
	TargetPrice float64            `bson:"target_price" json:"target_price"`
	Username    string             `bson:"username" json:"username"`
	AskedAt     string             `bson:"asked_at" json:"asked_at"`
	Status      string             `bson:"status" json:"status"`
	ServedAt    string             `bson:"served_at,omitempty" json:"served_at,omitempty"`
}
