//go:build !test
// +build !test

package models

import (
	"time"
)

type Request struct {
	ID        string    `json:"id" bson:"_id,omitempty"` // bson tag is for MongoDB ObjectId
	Username  string    `json:"username" bson:"username"`
	Symbol    string    `json:"symbol" bson:"symbol"`
	Status    string    `json:"status" bson:"status"` // Pending, Approved, Disapproved
	DateAdded time.Time `json:"date_added" bson:"date_added"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}
