package router

import (
	"github.com/go-chi/chi"
	"github.com/natan10/marketspace-api/controllers"
)

func Router() chi.Router {
	r := chi.NewRouter()

	r.Route("/users", func(r chi.Router) {
		r.Post("/", controllers.UserController{}.CreateUser)
	})

	return r
}
