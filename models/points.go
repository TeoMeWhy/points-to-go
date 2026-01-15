package models

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	UUID             string    `json:"uuid" gorm:"primaryKey"`
	DescCustomerName string    `json:"customer_name"`
	CodCPF           *string   `json:"cpf" gorm:"unique"`
	DescEmail        *string   `json:"email" gorm:"unique"`
	IdTwitch         *string   `json:"twitch" gorm:"unique"`
	IdYouTube        *string   `json:"youtube" gorm:"unique"`
	IdBlueSky        *string   `json:"bluesky" gorm:"unique"`
	IdInstagram      *string   `json:"instagram" gorm:"unique"`
	NrPoints         int64     `json:"points"`
	CreatedAt        time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt        time.Time `json:"updated_at" gorm:"autoUpdateTime:true"`
}

func NewCustomer() *Customer {

	newUUID := uuid.New()
	strUUID := newUUID.String()

	return &Customer{
		UUID: strUUID,
	}

}

type Product struct {
	UUID                   string `gorm:"primaryKey"`
	DescProduct            string
	DescProductDescription string
	DescProductCategory    string
}

func NewProduct(name, description, category string) *Product {

	return &Product{
		UUID:                   uuid.New().String(),
		DescProduct:            name,
		DescProductDescription: description,
		DescProductCategory:    category,
	}

}

type Transaction struct {
	UUID          string    `json:"transaction_id" gorm:"primaryKey"`
	IdCustomer    string    `json:"customer_id" binding:"required"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
	VlPoints      int64     `json:"points" binding:"required"`
	DescSysOrigin string    `json:"system_origin" binding:"required"`
}

func NewTransaction(idCustomer string, points int64, origin string) *Transaction {
	return &Transaction{
		UUID:          uuid.New().String(),
		IdCustomer:    idCustomer,
		VlPoints:      points,
		DescSysOrigin: origin,
	}
}

type TransactionProduct struct {
	UUID          string `json:"transaction_product_id" gorm:"primaryKey"`
	IdTransaction string `json:"transaction_id"`
	CodProduct    string `json:"product_id" binding:"required"`
	QtdeProduct   int64  `json:"product_qtd" binding:"required"`
	VlProduct     int64  `json:"points" binding:"required"`
}

func NewTransactionProduct(
	idTransaction, idProduct string,
	qtdProduct, points int64) *TransactionProduct {

	return &TransactionProduct{
		UUID:          uuid.New().String(),
		IdTransaction: idTransaction,
		CodProduct:    idProduct,
		QtdeProduct:   qtdProduct,
		VlProduct:     points,
	}

}
