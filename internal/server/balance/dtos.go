package balance

type WithdrawRequestDTO struct {
	Sum     float64 `json:"sum" validate:"required,gt=0"`
	OrderID string  `json:"order" validate:"required"`
}
