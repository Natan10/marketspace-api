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

type AnnouncementsController struct {
	Service services.IAnnouncementsService
}

func (ac *AnnouncementsController) UploadAnnouncementPhotos(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(32 << 20)

	if err != nil {
		log.Println("Failed to parse multipart form", err)
		http.Error(w, "Failed to parse multipart form", http.StatusBadRequest)
		return
	}

	path, _ := os.Getwd()
	tempPath := `tmp/products`
	_ = os.MkdirAll(tempPath, os.ModePerm)
	joinPath := filepath.Join(path, tempPath)

	files, _ := r.MultipartForm.File["images"]

	uploadFiles := []int{}

	for i, fileHeader := range files {

		file, err := fileHeader.Open()

		if err != nil {
			log.Println("erro to upload file: " + fileHeader.Filename)
			continue
		}

		defer file.Close()

		// create new file upload
		destination, err := os.Create(fmt.Sprintf("%v/%v", joinPath, fileHeader.Filename))

		defer destination.Close()

		if err != nil {
			log.Println("Failed to create destination file", err)
			continue
		}

		_, err = io.Copy(destination, file)

		if err != nil {
			log.Println("Failed to copy file content", err)
			continue
		}

		uploadFiles = append(uploadFiles, i)
	}

	if len(uploadFiles) == len(files) {
		w.WriteHeader(200)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
}

// @Summary Get Announcement
// @Tags announcements
// @Accept json
// @Produce json
// @Param userId query int true "user id"
// @Param announcementId path int true "announcement id"
// @Success 200 {object} map[string]models.Announcement
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /announcements/{announcementId} [get]
func (ac *AnnouncementsController) Get(w http.ResponseWriter, r *http.Request) {
	param := r.URL.Query().Get("userId")

	if param == "" {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	userId, err := strconv.Atoi(param)

	if err != nil {
		log.Fatalf("Error:%v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	announcementId, err := strconv.Atoi(chi.URLParam(r, "announcementId"))

	if err != nil {
		log.Printf("Error:%v\n", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var announcement *models.Announcement

	announcement, err = ac.Service.GetAnnouncement(int64(userId), int64(announcementId))

	if announcement == nil && err == nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	} else if announcement == nil && err != nil {
		log.Printf("Error:%v\n", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	} else {
		response := map[string]models.Announcement{
			"data": *announcement,
		}

		w.Header().Add("Content-type", "application/json")
		json.NewEncoder(w).Encode(response)
	}
}

// @Summary Get Announcements
// @Tags announcements
// @Accept json
// @Produce json
// @Param userId query int true "user id"
// @Success 200 {object} map[string]models.Announcement
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /announcements [get]
func (ac *AnnouncementsController) GetAll(w http.ResponseWriter, r *http.Request) {
	userIdParam := r.URL.Query().Get("userId")

	var announcements []models.Announcement
	var err error
	var response map[string][]models.Announcement

	if userIdParam != "" {
		userId, err := strconv.Atoi(userIdParam)
		if err != nil {
			log.Fatalf("Error:%v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		announcements, err = ac.Service.GetAllAnnouncementsByUser(int64(userId))

		if err != nil {
			log.Fatalf("Error:%v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

	} else {
		announcements, err = ac.Service.GetAll()

		if err != nil {
			log.Fatalf("Error:%v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
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

// @Summary Create Announcement
// @Tags announcements
// @Accept json
// @Produce json
// @Param request body dtos.AnnouncementDTO true "announcement payload"
// @Success 200 {object} dtos.ResponseDTO
// @Failure 500 {string} string
// @Router /announcements [post]
func (ac *AnnouncementsController) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	var announcement dtos.AnnouncementDTO

	err := json.NewDecoder(r.Body).Decode(&announcement)

	if err != nil {
		log.Fatalf("Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var response dtos.ResponseDTO

	id, err := ac.Service.CreateAnnouncement(announcement)

	if err != nil {
		response = dtos.ResponseDTO{
			Error:   true,
			Message: fmt.Sprintf("Erro ao criar anuncio: %v", err),
		}

	} else {
		response = dtos.ResponseDTO{
			Error:   false,
			Message: fmt.Sprintf("Anuncio criado com sucesso: %v", id),
		}
	}

	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Update Announcement
// @Tags announcements
// @Accept json
// @Produce json
// @Param announcementId path int true "announcement id"
// @Param request body dtos.AnnouncementDTO true "announcement payload"
// @Success 200 {object} dtos.ResponseDTO
// @Failure 500 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Router /announcements/{announcementId} [put]
func (ac *AnnouncementsController) UpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
	announcementId, err := strconv.Atoi(chi.URLParam(r, "announcementId"))

	if err != nil {
		log.Fatalf("Error:%v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var announcement dtos.AnnouncementDTO

	err = json.NewDecoder(r.Body).Decode(&announcement)

	if err != nil {
		log.Fatalf("Error:%v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var response dtos.ResponseDTO

	rows, err := ac.Service.UpdateAnnouncement(int64(announcementId), announcement)

	if err != nil {
		response = dtos.ResponseDTO{
			Error:   true,
			Message: fmt.Sprintf("Erro ao atualizar anuncio: %v", err),
		}
	} else {
		if rows == 0 {
			log.Printf("Anuncio nao encontrado: %v\n", rows)
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		response = dtos.ResponseDTO{
			Error:   false,
			Message: "Anuncio atualizado com sucesso!",
		}
	}

	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// @Summary Delete Announcement
// @Tags announcements
// @Accept json
// @Produce json
// @Param announcementId path int true "announcement id"
// @Success 200 {object} dtos.ResponseDTO
// @Failure 500 {string} string
// @Failure 404 {string} string
// @Failure 400 {string} string
// @Router /announcements/{announcementId} [delete]
func (ac *AnnouncementsController) DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	announcementId, err := strconv.Atoi(chi.URLParam(r, "announcementId"))

	if err != nil {
		log.Fatalf("Error:%v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	rows, err := ac.Service.DeleteAnnouncement(int64(announcementId))

	var response dtos.ResponseDTO

	if err != nil {
		response = dtos.ResponseDTO{
			Error:   true,
			Message: fmt.Sprintf("Erro ao remover anuncion: %v\n", err),
		}
	}

	if rows == 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if rows > 1 {
		log.Printf("Numero de anuncios removidos errado: %v\n", rows)
	}

	response = dtos.ResponseDTO{
		Error:   false,
		Message: "Anuncio removido com sucesso!",
	}

	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}
