package configs

import (
	"log"

	"github.com/joho/godotenv"
)

func Load(env string) {
	var err error

	if env == "development" {
		err = godotenv.Load()
	} else {
		err = godotenv.Load(".env.production")
	}

	if err != nil {
		log.Fatal("Error to loading env variables")
		panic(err)
	}

	log.Printf("Load variables from %v\n", env)
}
