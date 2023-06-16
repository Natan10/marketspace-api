package configs

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func OpenConn() (*sql.DB, error) {
	host := "localhost"         //os.Getenv("POSTGRES_HOST")
	port := "5432"              //os.Getenv("POSTGRES_PORT")
	user := "postgres"          //os.Getenv("POSTGRES_USER")
	password := "workspace2402" // os.Getenv("POSTGRES_PASSWORD")
	dbName := "workspace"       // os.Getenv("DB_NAME")

	connection := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=disable", host, port, user, password, dbName)

	db, err := sql.Open("postgres", connection)

	if err != nil {
		log.Fatal("Erro to connect db")
		panic(err)
	}

	err = db.Ping()

	return db, err
}
