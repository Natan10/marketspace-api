package services

import (
	"log"

	"github.com/natan10/marketspace-api/configs"
	"github.com/natan10/marketspace-api/dtos"
)

func CreateUser(user dtos.UserDTO) (id int64, err error) {
	db, err := configs.OpenConn()

	if err != nil {
		log.Fatalf("Error to connect db: %v", err)
		return
	}

	defer db.Close()

	sql := `INSERT INTO users (username, email, phone) values ($1,$2,$3) RETURNING id`

	err = db.QueryRow(sql, user.Username, user.Email, user.Phone).Scan(&id)

	return id, err
}
