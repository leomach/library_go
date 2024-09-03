package models

import (
	"library_api/pkg/db"
	"time"
)

type Payment struct {
	ID            int       `json:"id"`
	BookID        int       `json:"book_id"`
	PaymentValue  float64   `json:"payment_value"`
	PaymentMethod string    `json:"payment_method"`
	PaymentDate   time.Time `json:"payment_date"`
}

func (p *Payment) CreatePayment() error {
	sqlPayment := `
    INSERT INTO payments (book_id, payment_value, payment_method, payment_date)
    VALUES (?, ?, ?, NOW())`

	stmtPay, err := db.DB.Prepare(sqlPayment)
	if err != nil {
		return err
	}
	defer stmtPay.Close()
	_, err = stmtPay.Exec(p.BookID, p.PaymentValue, p.PaymentMethod)
	if err != nil {
		return err
	}
	return nil
}