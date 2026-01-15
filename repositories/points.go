package repositories

import (
	"points/models"

	"gorm.io/gorm"
)

type PointsRepo struct {
	db *gorm.DB
}

func (pr *PointsRepo) CreateCustomer(customer *models.Customer) error {

	res := pr.db.Create(customer)
	if res.Error != nil {
		res.Rollback()
		return res.Error
	}

	res.Commit()
	return nil

}

func (pr *PointsRepo) GetCustomerByID(id string) (*models.Customer, error) {

	customer := &models.Customer{}

	result := pr.db.First(&customer, "uuid = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}

	return customer, nil
}

func (pr *PointsRepo) GetCustomer(id, cpf, email, twitch, youtube, bsky, instagram string) []models.Customer {

	var customers []models.Customer

	query := pr.db.Model(&models.Customer{})

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

func (pr *PointsRepo) GetCustomerTransactions(customerID string) ([]models.Transaction, error) {
	return getCustomerTransactions(customerID, pr.db)
}

func (pr *PointsRepo) DeleteCustomer(c *models.Customer) error {
	return deleteCustomer(c, pr.db)
}

func (pr *PointsRepo) UpdateCustomer(c *models.Customer) error {
	return updateCustomer(c, pr.db)
}

func (pr *PointsRepo) CreateTransactionWithProducts(c *models.Customer, t *models.Transaction, ps []models.TransactionProduct) error {

	tx := pr.db.Begin()
	defer tx.Commit()

	if err := tx.Create(t).Error; err != nil {
		tx.Rollback()
		return err
	}

	for _, p := range ps {
		if err := tx.Create(p).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	c.NrPoints += t.VlPoints
	if err := tx.Save(c).Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (pr *PointsRepo) UpdateTransactions(txs []models.Transaction) error {
	return updateTransactions(txs, pr.db)
}

func (pr *PointsRepo) MigrateCustomerData(cFrom, cTo *models.Customer) error {

	tx := pr.db.Begin()
	defer tx.Commit()

	oldTransactions, err := getCustomerTransactions(cFrom.UUID, tx)
	if err != nil {
		return err
	}

	transactionsNew := []models.Transaction{}
	for _, t := range oldTransactions {
		t.IdCustomer = cTo.UUID
		transactionsNew = append(transactionsNew, t)
	}

	if err := updateTransactions(transactionsNew, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := deleteCustomer(cFrom, tx); err != nil {
		tx.Rollback()
		return err
	}

	if err := updateCustomer(cTo, tx); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func getCustomerTransactions(customerID string, db *gorm.DB) ([]models.Transaction, error) {
	transactionsOld := []models.Transaction{}
	if err := db.Model(&models.Transaction{}).Where("id_customer = ?", customerID).Find(&transactionsOld).Error; err != nil {
		return nil, err
	}

	return transactionsOld, nil
}

func updateTransactions(txs []models.Transaction, db *gorm.DB) error {
	for _, t := range txs {
		if err := db.Save(&t).Error; err != nil {
			return err
		}
	}

	return nil
}

func updateCustomer(c *models.Customer, db *gorm.DB) error {
	if c.UUID == "" {
		return gorm.ErrMissingWhereClause
	}

	return db.Save(c).Error
}

func deleteCustomer(c *models.Customer, db *gorm.DB) error {
	if c.UUID == "" {
		return gorm.ErrMissingWhereClause
	}

	return db.Delete(c).Error
}

func NewPointsRepo(db *gorm.DB) *PointsRepo {
	return &PointsRepo{db: db}
}
