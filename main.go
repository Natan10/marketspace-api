package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/natan10/marketspace-api/configs"
	"github.com/natan10/marketspace-api/router"
)

func init() {
	configs.Load()
}

func main() {
	port := os.Getenv("SERVER_PORT")
	ch := chi.NewRouter()

	ch.Use(middleware.Heartbeat("/healthy"))
	ch.Use(middleware.Logger)

	ch.Route("/v1", func(ch chi.Router) {
		ch.Mount("/", router.Router())
	})

	http.ListenAndServe(fmt.Sprintf(":%v", port), ch)
}

// TODO - AUTH
// TODO - USER CRUD
// TODO - user get
// TODO - user create
// TODO - user update
// TODO - ANNOUCEMENT CRUD
