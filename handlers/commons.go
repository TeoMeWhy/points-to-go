package handlers

import (
	"points/models"

	"gorm.io/gorm"
)

type PayloadMigrationCustomer struct {
	Customer1UUID string `json:"customer_1" binding:"required"`
	Customer2UUID string `json:"customer_2" binding:"required"`
}

type PayloadTransaction struct {
	models.Transaction
	Products []models.TransactionProduct `json:"products" binding:"required"`
}

func TransactionProcess(
	c *models.Customer,
	t *models.Transaction,
	ps []models.TransactionProduct,
	db *gorm.DB) error {

	tx := db.Begin()
	defer tx.Commit()

	res := tx.Create(t)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	for _, p := range ps {
		res := tx.Create(p)
		if res.Error != nil {
			tx.Rollback()
			return res.Error
		}
	}

	c.NrPoints += t.VlPoints
	res = tx.Save(c)
	if res.Error != nil {
		tx.Rollback()
		return res.Error
	}

	return nil
}
