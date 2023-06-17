package dtos

type AuthUserDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthUserResponseDTO struct {
	Token string `json:"token"`
}
