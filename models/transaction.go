package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Transaction struct {
	UUID          string    `json:"transaction_id" gorm:"primaryKey"`
	IdCustomer    string    `json:"customer_id" binding:"required"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	VlPoints      int64     `json:"points" binding:"required"`
	DescSysOrigin string    `json:"system_origin" binding:"required"`
}

func UpdateTransactions(txs []Transaction, db *gorm.DB) error {

	for _, t := range txs {
		if err := db.Save(&t).Error; err != nil {
			return err
		}
	}

	return nil
}

func UpdateTransaction(t *Transaction, db *gorm.DB) error {

	if t.UUID == "" {
		return gorm.ErrMissingWhereClause
	}

	return db.Save(t).Error
}

func NewTransaction(idCustomer string, points int64, origin string) *Transaction {
	return &Transaction{
		UUID:          uuid.New().String(),
		IdCustomer:    idCustomer,
		VlPoints:      points,
		DescSysOrigin: origin,
	}
}
