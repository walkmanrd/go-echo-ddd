package models

import (
	"github.com/lib/pq"
)

// Expense is a model for expense
type Expense struct {
	ID     string         `json:"id"`
	Title  string         `json:"title"`
	Amount float64        `json:"amount"`
	Note   string         `json:"note"`
	Tags   pq.StringArray `json:"tags"`
}
