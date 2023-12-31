package services

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/lib/pq"
	"github.com/natan10/marketspace-api/configs"
	"github.com/natan10/marketspace-api/dtos"
	"github.com/natan10/marketspace-api/models"
	"github.com/natan10/marketspace-api/utils"
)

type IAnnouncementsService interface {
	GetAll(params map[string]interface{}) (announcements []models.Announcement, err error)
	GetAnnouncement(useId int64, announcementId int64) (announcement *models.Announcement, err error)
	GetAllAnnouncementsByUser(user_id int64, param string) (announcements []models.Announcement, err error)
	CreateAnnouncement(an dtos.AnnouncementDTO) (id int16, err error)
	UpdateAnnouncement(id int64, an dtos.AnnouncementDTO) (int64, error)
	DeleteAnnouncement(id int64) (int64, error)
}

type AnnouncementsService struct{}

func (s *AnnouncementsService) GetAll(params map[string]interface{}) (announcements []models.Announcement, err error) {
	db, err := configs.OpenConn()

	if err != nil {
		log.Fatalf("Error to connect db: %v", err)
		return
	}

	defer db.Close()

	sqlStatement := `
		SELECT  
			a.id,
			a.title,
			a.description,
			a.is_new,
			a.price,
			a.is_exchangeable,
			a.is_active,
			a.images,
			a.user_id,
			p.boleto,
			p.pix,
			p.cash,
			p.credit_card,
			p.bank_deposit
		FROM announcements a INNER JOIN payment_methods p
		ON p.announcement_id = a.id and a.is_active = 'true'
	`

	filterParams := utils.MountFilterQuery(params)

	if filterParams != "" {
		sqlStatement += fmt.Sprintf("WHERE %v", filterParams)
	}

	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Printf("Error query: %v\n", err)
		return
	}

	for rows.Next() {
		var an models.Announcement

		err := rows.Scan(
			&an.Id,
			&an.Title,
			&an.Description,
			&an.IsNew,
			&an.Price,
			&an.IsExchangeable,
			&an.IsActive,
			pq.Array(&an.Images),
			&an.UserId,
			&an.Boleto,
			&an.Pix,
			&an.Cash,
			&an.CreditCard,
			&an.BankDeposit,
		)

		if err != nil {
			log.Printf("Error scan: %v\n", err)
			continue
		}

		announcements = append(announcements, an)
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error ao buscar registros")
	}
	defer rows.Close()
	return announcements, nil
}

func (s *AnnouncementsService) GetAnnouncement(userId int64, announcementId int64) (announcement *models.Announcement, err error) {
	db, err := configs.OpenConn()

	if err != nil {
		log.Fatalf("Error to connect db: %v", err)
		return
	}

	defer db.Close()

	sqlStatement := `
		SELECT
			a.id,
			a.title,
			a.description,
			a.is_new,
			a.price,
			a.is_exchangeable,
			a.is_active,
			a.images,
			a.user_id,
			p.boleto,
			p.pix,
			p.cash,
			p.credit_card,
			p.bank_deposit
		FROM announcements a INNER JOIN payment_methods p
		ON p.announcement_id = a.id and p.announcement_id=$2
		WHERE a.user_id=$1
	`
	announcement = &models.Announcement{}

	err = db.QueryRow(sqlStatement, userId, announcementId).Scan(
		&announcement.Id,
		&announcement.Title,
		&announcement.Description,
		&announcement.IsNew,
		&announcement.Price,
		&announcement.IsExchangeable,
		&announcement.IsActive,
		pq.Array(&announcement.Images),
		&announcement.UserId,
		&announcement.Boleto,
		&announcement.Pix,
		&announcement.Cash,
		&announcement.CreditCard,
		&announcement.BankDeposit,
	)

	switch err {
	case sql.ErrNoRows:
		log.Printf("Error query: %v\n", err)
		return nil, nil
	case nil:
		return announcement, nil
	default:
		return nil, err
	}
}

func (s *AnnouncementsService) GetAllAnnouncementsByUser(useId int64, param string) (announcements []models.Announcement, err error) {
	db, err := configs.OpenConn()

	if err != nil {
		log.Fatalf("Error to connect db: %v", err)
		return
	}

	defer db.Close()

	sqlStatement := `
		SELECT  
			a.id,
			a.title,
			a.description,
			a.is_new,
			a.price,
			a.is_exchangeable,
			a.is_active,
			a.images,
			a.user_id,
			p.boleto,
			p.pix,
			p.cash,
			p.credit_card,
			p.bank_deposit
		FROM announcements a INNER JOIN payment_methods p
		ON p.announcement_id = a.id
		WHERE a.user_id=$1
	`

	var paramFilter string = ""

	if param != "" {
		switch param {
		case "used":
			paramFilter = "is_new='false'"
		case "new":
			paramFilter = "is_new='true'"
		default:
		}
	}

	if paramFilter != "" {
		sqlStatement += " and " + paramFilter
	}

	rows, err := db.Query(sqlStatement, useId)

	if err != nil {
		log.Printf("Error query: %v\n", err)
		return
	}

	for rows.Next() {
		var an models.Announcement

		err := rows.Scan(
			&an.Id,
			&an.Title,
			&an.Description,
			&an.IsNew,
			&an.Price,
			&an.IsExchangeable,
			&an.IsActive,
			pq.Array(&an.Images),
			&an.UserId,
			&an.Boleto,
			&an.Pix,
			&an.Cash,
			&an.CreditCard,
			&an.BankDeposit,
		)

		if err != nil {
			log.Printf("Error scan: %v\n", err)
			continue
		}

		announcements = append(announcements, an)
	}

	if err := rows.Err(); err != nil {
		fmt.Println("Error ao buscar registros")
	}

	defer rows.Close()

	return announcements, nil
}

func (s *AnnouncementsService) CreateAnnouncement(an dtos.AnnouncementDTO) (id int16, err error) {
	db, err := configs.OpenConn()

	if err != nil {
		log.Fatalf("Error to connect db: %v", err)
		return
	}

	defer db.Close()

	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		log.Printf("Erro to start transaction: %v", err)

		return 0, err
	}

	anSql := `INSERT into 
		announcements (title, description, is_new, price, is_exchangeable, is_active, images, user_id)
		VALUES ($1,$2,$3,$4,$5,$6,$7, $8) RETURNING id
	`

	err = tx.QueryRow(anSql, an.Title, an.Description, an.IsNew, an.Price, an.IsExchangeable, an.IsActive, pq.Array(an.Images), an.UserId).Scan(&id)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	pmSql := `INSERT into 
		payment_methods (boleto, pix, cash, credit_card, bank_deposit, announcement_id)
		VALUES ($1,$2,$3,$4,$5,$6)
	`
	_, err = tx.Exec(pmSql, an.Boleto, an.Pix, an.Cash, an.CreditCard, an.BankDeposit, id)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Erro to process transaction: %v", err)

		return 0, err
	}

	return id, nil
}

func (s *AnnouncementsService) UpdateAnnouncement(id int64, an dtos.AnnouncementDTO) (int64, error) {
	db, err := configs.OpenConn()

	if err != nil {
		log.Fatalf("Error to connect db: %v", err)
		return 0, err
	}

	defer db.Close()

	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)

	if err != nil {
		log.Printf("Erro to start transaction: %v", err)

		return 0, err
	}

	anSql := `UPDATE announcements 
		SET title=$2,description=$3,is_new=$4,price=$5,is_exchangeable=$6,images=$7,is_active=$8 
		WHERE id=$1
	`

	paymentsSql := `UPDATE payment_methods
		SET boleto=$2, pix=$3, cash=$4, credit_card=$5, bank_deposit=$6
		WHERE announcement_id=$1
	`

	// update announcement
	rows, err := tx.Exec(anSql, id, an.Title, an.Description, an.IsNew, an.Price, an.IsExchangeable, pq.Array(an.Images), an.IsActive)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	if numberRowsAffected, err := rows.RowsAffected(); numberRowsAffected > 1 {
		tx.Rollback()
		return 0, err
	}

	// update payment methods based on announcement id
	if _, err := tx.Exec(paymentsSql, id, an.Boleto, an.Pix, an.Cash, an.CreditCard, an.BankDeposit); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		log.Fatalf("Error to process transaction: %v", err)
		return 0, err
	}

	return rows.RowsAffected()
}

func (s *AnnouncementsService) DeleteAnnouncement(id int64) (int64, error) {
	db, err := configs.OpenConn()

	if err != nil {
		log.Fatalf("Error to connect db:%v", err)
		return 0, err
	}

	defer db.Close()

	rows, err := db.Exec("DELETE from announcements WHERE id=$1", id)

	if err != nil {
		return 0, err
	}

	return rows.RowsAffected()
}
