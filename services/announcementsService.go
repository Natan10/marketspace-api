package services

import (
	"context"
	"log"

	"github.com/lib/pq"
	"github.com/natan10/marketspace-api/configs"
	"github.com/natan10/marketspace-api/dtos"
)

func CreateAnnouncement(an dtos.AnnouncementDTO) (id int16, err error) {
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

func UpdateAnnouncement(id int64, an dtos.AnnouncementDTO) (int64, error) {
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

func DeleteAnnouncement(id int64) (int64, error) {
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
