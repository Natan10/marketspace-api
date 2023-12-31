package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/natan10/marketspace-api/dtos"
	"github.com/natan10/marketspace-api/models"
	"github.com/natan10/marketspace-api/services"
)

type UserController struct {
	UserService         services.IUserService
	AnnouncementService services.IAnnouncementsService
}

func (us *UserController) UploadUserAvatar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(20 << 20)

	if err != nil {
		log.Println("Failed to parse multipart form", err)
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	filename := r.Form.Get("filename")
	file, _, err := r.FormFile("file")

	if err != nil {
		log.Println("Failed to read file", err)
		http.Error(w, "Failed to read file", http.StatusBadRequest)
		return
	}

	defer file.Close()

	// get path and create dir
	path, _ := os.Getwd()
	tempPath := `tmp/avatars`
	_ = os.MkdirAll(tempPath, os.ModePerm)

	// create file
	dstFile, err := os.Create(fmt.Sprintf("%v/%v", filepath.Join(path, tempPath), filename))

	if err != nil {
		log.Println("Failed to create destination file", err)
		http.Error(w, "Failed to create destination file", http.StatusInternalServerError)
		return
	}

	defer dstFile.Close()

	_, err = io.Copy(dstFile, file)

	if err != nil {
		log.Println("Failed to copy file content", err)
		http.Error(w, "Failed to copy file content", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
}

// @Summary Get User Announcements
// @Tags users
// @Accept json
// @Produce json
// @Param userId query int true "user id"
// @Success 200 {object} map[string]models.Announcement
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /users/{userId}/announcements [get]
func (ac *UserController) GetUserAnnouncements(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
	param := r.URL.Query().Get("type")

	if err != nil {
		log.Fatalf("Error:%v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var announcements []models.Announcement
	var response map[string][]models.Announcement

	if err != nil {
		log.Fatalf("Error:%v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	fmt.Println("params: ", param)
	announcements, err = ac.AnnouncementService.GetAllAnnouncementsByUser(int64(userId), param)

	if err != nil {
		log.Fatalf("Error:%v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if len(announcements) > 0 {
		response = map[string][]models.Announcement{
			"data": announcements,
		}
	} else {
		announcements = make([]models.Announcement, 0)
		response = map[string][]models.Announcement{
			"data": announcements,
		}
	}
	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	user, err = us.UserService.GetUserById(int64(userId))

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
// @Router /signup [post]
func (uc *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user dtos.UserDTO

	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatalf("Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var response dtos.ResponseDTO

	if id, err := uc.UserService.CreateUser(user); err != nil {
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

	rows, err := uc.UserService.UpdateUser(int64(userId), user)

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
