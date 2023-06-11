package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/natan10/marketspace-api/configs"
	"github.com/natan10/marketspace-api/router"

	_ "github.com/natan10/marketspace-api/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func init() {
	configs.Load()
}

// @title MarketSpace Api
// @version 1.0
// @description This is a web server for MarketSpace application.
// @termsOfService http://swagger.io/terms/

// @host localhost:8000
// @BasePath /v1
func main() {
	port := os.Getenv("SERVER_PORT")
	ch := chi.NewRouter()

	ch.Use(middleware.Heartbeat("/healthy"))
	ch.Use(middleware.Logger)

	ch.Route("/v1", func(ch chi.Router) {
		ch.Get("/swagger/*", httpSwagger.Handler(
			httpSwagger.URL(fmt.Sprintf("http://localhost:%v/v1/swagger/doc.json", port)),
		))
		ch.Mount("/", router.Router())
	})

	log.Printf("Server running at %v \n", port)
	http.ListenAndServe(fmt.Sprintf(":%v", port), ch)
}

// TODO - AUTH
