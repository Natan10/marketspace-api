package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/natan10/marketspace-api/dtos"
	"github.com/natan10/marketspace-api/services"
)

type UserController struct{}

func (uc UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dtos.UserDTO

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var response map[string]any

	if id, err := services.CreateUser(user); err != nil {
		response = map[string]any{
			"Error":   true,
			"Message": fmt.Sprintf("Erro ao criar usuario: %v", err),
		}
	} else {
		response = map[string]any{
			"Error":   false,
			"Message": fmt.Sprintf("Usuario criado com sucesso: %v", id),
		}
	}

	json.NewEncoder(w).Encode(response)
}
