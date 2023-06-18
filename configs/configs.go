package configs

import (
	"log"

	"github.com/joho/godotenv"
)

func Load(env string) {
	var err error

	switch env {
	case "development":
		err = godotenv.Load()
	case "test":
		err = godotenv.Load(".env.test")
	case "production":
		err = godotenv.Load(".env.production")
	default:
		err = godotenv.Load()
	}

	if err != nil {
		log.Fatal("Error to loading env variables")
		panic(err)
	}

	log.Printf("Load variables from %v\n", env)
}
