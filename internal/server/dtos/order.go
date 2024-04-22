package dtos

import "time"

type Order struct {
	ID     string `json:"number"`
	UserID int    `json:"-"`
	// The accrual points for the order, if available
	// This field is optional in the JSON response
	Accrual   *float64  `json:"accrual,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"uploaded_at"`
	UpdatedAt time.Time `json:"-"`
}
