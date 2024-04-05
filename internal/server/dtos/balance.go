package dtos

import "time"

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
	UserID    int     `json:"-"`
}

type Withdraw struct {
	ID          int       `json:"-"`
	OrderID     string    `json:"order"`
	Amount      float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
	UserID      int       `json:"-"`
}
