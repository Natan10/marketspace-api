package router

import (
	"github.com/go-chi/chi"
	"github.com/natan10/marketspace-api/controllers"
)

func Router() chi.Router {
	r := chi.NewRouter()

	r.Route("/users", func(r chi.Router) {
		r.Post("/", controllers.UserController{}.CreateUser)
		r.Post("/{userId}", controllers.UserController{}.UpdateUser)
	})

	r.Route("/announcements", func(r chi.Router) {
		r.Post("/", controllers.AnnouncementsController{}.CreateAnnouncement)
		r.Post("/{announcementId}", controllers.AnnouncementsController{}.UpdateAnnouncement)
	})

	return r
}
