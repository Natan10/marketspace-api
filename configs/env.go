package configs

import (
	"log"

	"github.com/joho/godotenv"
)

func Load() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error to loading env variables")
		panic(err)
	}
}
