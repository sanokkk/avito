package requests

type RegisterDto struct {
	Username string `json:"username" validate:"required,gte=5"`
	Password string `json:"password" validate:"required,gte=8"`
}

type LoginDto struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}
