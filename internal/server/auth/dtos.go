package auth

type RegisterRequestDTO struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required,min=4,max=32"`
}
