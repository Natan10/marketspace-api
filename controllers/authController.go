package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/natan10/marketspace-api/dtos"
	"github.com/natan10/marketspace-api/models"
	"github.com/natan10/marketspace-api/services"
)

type AuthController struct {
	Service services.IUserService
}

// @Summary Auth User
// @Tags authentication
// @Accept json
// @Produce json
// @Param request body dtos.AuthUserDTO true "auth payload"
// @Success 200 {object} dtos.AuthUserResponseDTO "response"
// @Router /signin [post]
func (auth AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	var payload dtos.AuthUserDTO

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email := payload.Email
	password := payload.Password

	var user *models.User

	user, err = auth.Service.GetUser(email, password)

	if user == nil && err == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// generate token
	_, token, err := services.TokenService{}.EncodeToken(user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tokenResponse := dtos.AuthUserResponseDTO{
		Token: token,
	}

	json.NewEncoder(w).Encode(tokenResponse)
}
