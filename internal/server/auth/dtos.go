package auth

type RegisterRequestDTO struct {
	Username string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required,min=4,max=32"`
}

type LoginRequestDTO struct {
	Username string `json:"login" validate:"required"`
	Password string `json:"password" validate:"required,min=4,max=32"`
}
