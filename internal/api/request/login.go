package request

type UserLogin struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Status string `json:"status"`
	Token  string `json:"token"`
}
