package router

import (
	"github.com/go-chi/chi"
	"github.com/natan10/marketspace-api/controllers"
	"github.com/natan10/marketspace-api/services"
)

var announcementService services.AnnouncementsService

func Router() chi.Router {
	r := chi.NewRouter()

	r.Route("/users", func(r chi.Router) {
		r.Post("/", controllers.UserController{}.CreateUser)
		r.Post("/{userId}", controllers.UserController{}.UpdateUser)
	})

	AnnouncementController := controllers.AnnouncementsController{
		Service: &announcementService,
	}

	r.Route("/announcements", func(r chi.Router) {
		r.Get("/", AnnouncementController.GetAll)
		r.Post("/", AnnouncementController.CreateAnnouncement)
		r.Put("/{announcementId}", AnnouncementController.UpdateAnnouncement)
		r.Delete("/{announcementId}", AnnouncementController.DeleteAnnouncement)
	})

	return r
}
