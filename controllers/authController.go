package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/natan10/marketspace-api/models"
	"github.com/natan10/marketspace-api/services"
)

type AuthController struct {
	Service services.IUserService
}

func (auth AuthController) SignIn(w http.ResponseWriter, r *http.Request) {
	var payload map[string]interface{}

	err := json.NewDecoder(r.Body).Decode(&payload)

	if err != nil {

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email := payload["email"].(string)
	password := payload["password"].(string)

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

	json.NewEncoder(w).Encode(map[string]any{
		"token": token,
	})
}
