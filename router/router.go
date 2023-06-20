package router

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/jwtauth/v5"
	"github.com/natan10/marketspace-api/controllers"
	"github.com/natan10/marketspace-api/services"
)

var (
	announcementService services.AnnouncementsService
	userService         services.UserService
	tokenAuth           = services.TokenAuth
)

func Router() chi.Router {
	r := chi.NewRouter()

	AnnouncementController := controllers.AnnouncementsController{
		Service: &announcementService,
	}

	UserController := controllers.UserController{
		Service: userService,
	}

	AuthController := controllers.AuthController{
		Service: userService,
	}

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(tokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Route("/users", func(r chi.Router) {
			r.Post("/", UserController.CreateUser)
			r.Get("/{userId}", UserController.GetUserInformation)
			r.Put("/{userId}", UserController.UpdateUser)
		})

		r.Route("/announcements", func(r chi.Router) {

			r.Get("/", AnnouncementController.GetAll)
			r.Get("/{announcementId}", AnnouncementController.Get)
			r.Post("/", AnnouncementController.CreateAnnouncement)
			r.Put("/{announcementId}", AnnouncementController.UpdateAnnouncement)
			r.Delete("/{announcementId}", AnnouncementController.DeleteAnnouncement)
		})
	})

	// public
	r.Group(func(r chi.Router) {
		r.Post("/signin", AuthController.SignIn)
	})

	return r
}
