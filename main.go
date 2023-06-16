package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/go-chi/jwtauth/v5"
	"github.com/natan10/marketspace-api/configs"
	_ "github.com/natan10/marketspace-api/docs"
	"github.com/natan10/marketspace-api/router"
	httpSwagger "github.com/swaggo/http-swagger"
)

func init() {
	configs.Load()
}

func AuthToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenA, _, err := jwtauth.FromContext(r.Context())

		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// jwtauth.VerifyToken(tokenAuth, token)

		// fmt.Printf("Token is Here: %v", token)
		fmt.Println(tokenA)

		next.ServeHTTP(w, r)
	})
}

// @title MarketSpace Api
// @version 1.0
// @description This is a web server for MarketSpace application.
// @termsOfService http://swagger.io/terms/

// @host localhost:8000
// @BasePath /v1
func main() {
	port := "8000" //os.Getenv("SERVER_PORT")
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
