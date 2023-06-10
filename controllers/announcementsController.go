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

type AnnouncementsController struct {
	Service services.IAnnouncementsService
}

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

func (ac *AnnouncementsController) GetAll(w http.ResponseWriter, r *http.Request) {

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

	announcements, err := ac.Service.GetAllAnnouncements(int64(userId))

	fmt.Println(announcements)

	if err != nil {
		log.Fatalf("Error:%v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var response map[string][]models.Announcement

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

func (ac *AnnouncementsController) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	var announcement dtos.AnnouncementDTO

	err := json.NewDecoder(r.Body).Decode(&announcement)

	if err != nil {
		log.Fatalf("Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var response map[string]any

	id, err := ac.Service.CreateAnnouncement(announcement)

	if err != nil {
		response = map[string]any{
			"Error":   true,
			"Message": fmt.Sprintf("Erro ao criar anuncio: %v", err),
		}
	} else {
		response = map[string]any{
			"Error":   false,
			"Message": fmt.Sprintf("Anuncio criado com sucesso: %v", id),
		}
	}

	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

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

	var response map[string]any

	rows, err := ac.Service.UpdateAnnouncement(int64(announcementId), announcement)

	if err != nil {
		response = map[string]any{
			"Error":   true,
			"Message": fmt.Sprintf("Erro ao atualizar anuncio: %v", err),
		}
	}

	if rows == 0 {
		log.Printf("Anuncio nao encontrado: %v\n", rows)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	response = map[string]any{
		"Error":   false,
		"Message": "Anuncio atualizado com sucesso!",
	}

	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (ac *AnnouncementsController) DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	announcementId, err := strconv.Atoi(chi.URLParam(r, "announcementId"))

	if err != nil {
		log.Fatalf("Error:%v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	rows, err := ac.Service.DeleteAnnouncement(int64(announcementId))

	var response map[string]any

	if err != nil {
		response = map[string]any{
			"Error":   true,
			"Message": fmt.Sprintf("Erro ao remover anuncion: %v\n", err),
		}
	}

	if rows == 0 {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if rows > 1 {
		log.Printf("Numero de anuncios removidos errado: %v\n", rows)
	}

	response = map[string]any{
		"Error":   false,
		"Message": "Anuncio removido com sucesso!",
	}

	w.Header().Add("Content-type", "application/json")
	json.NewEncoder(w).Encode(response)
}
