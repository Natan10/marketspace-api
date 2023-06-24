package services

import (
	"database/sql"
	"log"

	"github.com/natan10/marketspace-api/configs"
	"github.com/natan10/marketspace-api/dtos"
	"github.com/natan10/marketspace-api/models"
)

type IUserService interface {
	GetUserById(userId int64) (user *models.User, err error)
	GetUser(email string, password string) (user *models.User, err error)
	CreateUser(user dtos.UserDTO) (id int64, err error)
	UpdateUser(id int64, user dtos.UserDTO) (int64, error)
}

type UserService struct{}

func (u UserService) GetUserById(userId int64) (user *models.User, err error) {
	db, err := configs.OpenConn()

	if err != nil {
		log.Fatalf("Error to connect db: %v", err)
		return
	}

	defer db.Close()

	sqlStatement := `SELECT u.email, u.username, u.phone, u.photo from users u WHERE u.id=$1`

	user = new(models.User)

	err = db.QueryRow(sqlStatement, userId).Scan(&user.Email, &user.Name, &user.Phone, &user.Photo)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return user, nil
	default:
		return nil, err
	}
}

func (u UserService) GetUser(email string, password string) (user *models.User, err error) {
	db, err := configs.OpenConn()

	if err != nil {
		log.Fatalf("Error to connect db: %v", err)
		return
	}

	defer db.Close()

	sqlStatement := `SELECT u.id, u.email, u.username FROM users u WHERE email=$1 AND password=$2`

	user = new(models.User)

	err = db.QueryRow(sqlStatement, email, password).Scan(&user.Id, &user.Email, &user.Name)

	switch err {
	case sql.ErrNoRows:
		return nil, nil
	case nil:
		return user, nil
	default:
		return nil, err
	}
}

func (u UserService) CreateUser(user dtos.UserDTO) (id int64, err error) {
	db, err := configs.OpenConn()

	if err != nil {
		log.Fatalf("Error to connect db: %v", err)
		return
	}

	defer db.Close()

	sql := `INSERT INTO users (username, email, phone, password, photo) values ($1,$2,$3,$4,$5) RETURNING id`

	err = db.QueryRow(sql, user.Username, user.Email, user.Phone, user.Password, user.Photo).Scan(&id)

	return id, err
}

func (u UserService) UpdateUser(id int64, user dtos.UserDTO) (int64, error) {
	db, err := configs.OpenConn()

	if err != nil {
		log.Fatalf("Error to connect db: %v", err)
		return 0, err
	}

	defer db.Close()

	sql := `UPDATE users SET username=$2,email=$3,phone=$4,photo=$5 WHERE id=$1`

	res, err := db.Exec(sql, id, user.Username, user.Email, user.Phone, user.Photo)

	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
