package models

import (
	"log"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (c *Customer) Create(db *gorm.DB) error {

	res := db.Create(c)
	if res.Error != nil {
		res.Rollback()
		return res.Error
	}

	res.Commit()
	return nil

}

func MigrateCustomers(c1, c2 *Customer, db *gorm.DB) error {

	cTo, cFrom := &Customer{}, &Customer{}

	if c1.CreatedAt.Unix() > c2.CreatedAt.Unix() {
		cTo, cFrom = c2, c1
	} else {
		cTo, cFrom = c1, c2
	}

	log.Printf("Customer to: %v", cTo)
	log.Printf("Customer from: %v", cFrom)

	tx := db.Begin()

	cTo.NrPoints += cFrom.NrPoints

	if cFrom.DescEmail != nil {
		cTo.DescEmail = cFrom.DescEmail
	}

	if cFrom.IdYouTube != nil {
		cTo.IdYouTube = cFrom.IdYouTube
	}

	if cFrom.IdTwitch != nil {
		cTo.IdTwitch = cFrom.IdTwitch
	}

	if cFrom.IdBlueSky != nil {
		cTo.IdBlueSky = cFrom.IdBlueSky
	}

	if cFrom.IdInstagram != nil {
		cTo.IdInstagram = cFrom.IdInstagram
	}

	log.Printf("Customer to: %v", cTo)
	log.Printf("Customer from: %v", cFrom)

	transactionsOld := []Transaction{}
	if err := tx.Model(&Transaction{}).Where("id_customer = ?", cFrom.UUID).Find(&transactionsOld).Error; err != nil {
		tx.Rollback()
		return err
	}

	transactionsNew := []Transaction{}
	for _, t := range transactionsOld {
		t.IdCustomer = cTo.UUID
		transactionsNew = append(transactionsNew, t)
	}

	if err := UpdateTransactions(transactionsNew, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Delete(cFrom).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Save(cTo).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		return err
	}

	return nil
}

func GetCustomer(id, cpf, email, twitch, youtube, bsky, instagram string, db *gorm.DB) []Customer {

	var customers []Customer

	query := db.Model(&Customer{})

	if id != "" {
		query = query.Where("uuid = ?", id)
	}

	if cpf != "" {
		query = query.Where("cod_cpf = ?", cpf)
	}

	if email != "" {
		query = query.Where("desc_email = ?", email)
	}

	if twitch != "" {
		query = query.Where("id_twitch = ?", twitch)
	}

	if youtube != "" {
		query = query.Where("id_youtube = ?", youtube)
	}

	if bsky != "" {
		query = query.Where("id_blue_sky = ?", bsky)
	}

	if instagram != "" {
		query = query.Where("id_instagram = ?", instagram)
	}

	query.Find(&customers)
	return customers

}
