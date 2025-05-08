package handler

type UserRegistrationDTO struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm-password"`
}

type UserLoginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
