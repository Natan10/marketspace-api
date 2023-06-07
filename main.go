package main

import (
	"fmt"
	"os"

	"github.com/natan10/marketspace-api/configs"
)

func init() {
	configs.Load()
}

func main() {
	fmt.Println("Iniciando")

	_, err := configs.OpenConn()

	fmt.Println(err)

	fmt.Println("PORT:", os.Getenv("PORT"))
}
