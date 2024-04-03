package dtos

import "time"

type Order struct {
	ID        string    `json:"number"`
	UserID    int       `json:"-"`
	Accrual   *float64  `json:"accrual,omitempty"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"-"`
}
