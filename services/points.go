package services

import (
	"points/models"
	"points/repositories"
)

type PayloadTransaction struct {
	models.Transaction
	Products []models.TransactionProduct `json:"products" binding:"required"`
}

type PayloadMigrationCustomer struct {
	Customer1UUID string `json:"customer_1" binding:"required"`
	Customer2UUID string `json:"customer_2" binding:"required"`
}

type PointsService struct {
	PointsRepo *repositories.PointsRepo
}

func (ps *PointsService) MigrateCustomers(payload *PayloadMigrationCustomer) error {

	c1, err := ps.PointsRepo.GetCustomerByID(payload.Customer1UUID)
	if err != nil {
		return err
	}

	c2, err := ps.PointsRepo.GetCustomerByID(payload.Customer2UUID)
	if err != nil {
		return err
	}

	cTo, cFrom := &models.Customer{}, &models.Customer{}

	if c1.CreatedAt.Unix() > c2.CreatedAt.Unix() {
		cTo, cFrom = c2, c1
	} else {
		cTo, cFrom = c1, c2
	}

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

	return ps.PointsRepo.MigrateCustomerData(cFrom, cTo)
}

func (ps *PointsService) MakeTransaction(payload *PayloadTransaction) error {
	customer, err := ps.PointsRepo.GetCustomerByID(payload.IdCustomer)
	if err != nil {
		return err
	}

	t := models.NewTransaction(payload.IdCustomer, payload.VlPoints, payload.DescSysOrigin)
	trPrds := []models.TransactionProduct{}

	for _, p := range payload.Products {
		trPrds = append(trPrds, *models.NewTransactionProduct(t.UUID, p.CodProduct, p.QtdeProduct, p.VlProduct))
	}

	if err := ps.PointsRepo.CreateTransactionWithProducts(customer, t, trPrds); err != nil {
		return err
	}
	return nil
}

func NewPointsService(pointsRepo *repositories.PointsRepo) *PointsService {
	return &PointsService{
		PointsRepo: pointsRepo,
	}
}
