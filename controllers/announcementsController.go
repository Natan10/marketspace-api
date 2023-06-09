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

type AnnouncementsController struct{}

func (ac AnnouncementsController) CreateAnnouncement(w http.ResponseWriter, r *http.Request) {
	var announcement dtos.AnnouncementDTO

	err := json.NewDecoder(r.Body).Decode(&announcement)

	if err != nil {
		log.Fatalf("Error: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var response map[string]any

	id, err := services.CreateAnnouncement(announcement)

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

func (ac AnnouncementsController) UpdateAnnouncement(w http.ResponseWriter, r *http.Request) {
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

	rows, err := services.UpdateAnnouncement(int64(announcementId), announcement)

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

func (ac AnnouncementsController) DeleteAnnouncement(w http.ResponseWriter, r *http.Request) {
	announcementId, err := strconv.Atoi(chi.URLParam(r, "announcementId"))

	if err != nil {
		log.Fatalf("Error:%v", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	rows, err := services.DeleteAnnouncement(int64(announcementId))

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
