package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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

func (uc UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {

	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))

	if err != nil {
		log.Fatalf("Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var user dtos.UserDTO
	err = json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var response map[string]any

	rows, err := services.UpdateUser(int64(userId), user)

	if err != nil {
		log.Printf("Erro ao atualizar usuario: %v\n", err)

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if rows > 1 {
		log.Printf("Numero de usarios atualizados incorreto: %v\n", rows)
	}

	response = map[string]any{
		"Error":   false,
		"Message": "Usuario atualizado com sucesso!",
	}

	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}
