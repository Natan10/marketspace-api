package dtos

type ResponseDTO struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
}
