package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/natan10/marketspace-api/dtos"
	"github.com/natan10/marketspace-api/models"
	"github.com/natan10/marketspace-api/services"
)

type UserController struct {
	Service services.IUserService
}

// @Summary Update User
// @Tags users
// @Accept json
// @Produce json
// @Param userId path int true "user id"
// @Success 200 {object} dtos.UserDTO "response"
// @Router /users/{userId} [get]
func (us *UserController) GetUserInformation(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var user *models.User

	user, err = us.Service.GetUserById(int64(userId))

	if user == nil && err == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(dtos.UserDTO{
		Email:    user.Email,
		Username: user.Name,
		Phone:    user.Phone,
		Photo:    user.Photo,
	})

}

// @Summary Create User
// @Tags users
// @Accept json
// @Produce json
// @Param request body dtos.UserDTO true "user payload"
// @Success 200 {object} dtos.ResponseDTO "response"
// @Router /users [post]
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dtos.UserDTO

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var response dtos.ResponseDTO

	if id, err := uc.Service.CreateUser(user); err != nil {
		response = dtos.ResponseDTO{
			Error:   true,
			Message: fmt.Sprintf("Erro ao criar usuario: %v", err),
		}

	} else {
		response = dtos.ResponseDTO{
			Error:   false,
			Message: fmt.Sprintf("Usuario criado com sucesso: %v", id),
		}
	}

	json.NewEncoder(w).Encode(response)
}

// @Summary Update User
// @Tags users
// @Accept json
// @Produce json
// @Param userId path int true "user id"
// @Param request body dtos.UserDTO true "user payload"
// @Success 200 {object} dtos.ResponseDTO "response"
// @Router /users/{userId} [put]
func (uc *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {

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

	rows, err := uc.Service.UpdateUser(int64(userId), user)

	if err != nil {
		log.Printf("Erro ao atualizar usuario: %v\n", err)

		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if rows > 1 {
		log.Printf("Numero de usarios atualizados incorreto: %v\n", rows)
	}

	response := dtos.ResponseDTO{
		Error:   false,
		Message: "Usuario atualizado com sucesso!",
	}

	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}
