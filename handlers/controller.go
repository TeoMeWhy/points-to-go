package handlers

import (
	"net/http"
	"points/models"
	"points/myerrors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	db *gorm.DB
}

func (ctrl *Controller) GetCustomers(c *gin.Context) {

	uuid := c.Query("uuid")
	cod_cpf := c.Query("cod_cpf")
	desc_email := c.Query("desc_email")
	twitch := c.Query("twitch")
	youtube := c.Query("youtube")
	blue_sky := c.Query("blue_sky")
	instagram := c.Query("instagram")

	customers := models.GetCustomer(
		uuid,
		cod_cpf,
		desc_email,
		twitch,
		youtube,
		blue_sky,
		instagram,
		ctrl.db,
	)

	if len(customers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "nenhum resultado encontrado"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

func (ctrl *Controller) GetCustomerByID(c *gin.Context) {

	var customer models.Customer
	id := c.Param("id")

	result := ctrl.db.First(&customer, "uuid = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
}

func (ctrl *Controller) PostCustomer(c *gin.Context) {

	customer := models.NewCustomer()
	customerUUID := customer.UUID

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customer.UUID = customerUUID

	err := myerrors.GetCreateCustomerErrors(customer.Create(ctrl.db))
	if err != nil {

		if err.Error() == myerrors.EmailCreateError.Error() {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusBadRequest, gin.H{"error": "solicitação incorreta"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"status": "created", "customer": customer})
}

func (ctrl *Controller) PutCustomer(c *gin.Context) {

	id := c.Param("id")

	var newCustomer, oldCustomer *models.Customer

	res := ctrl.db.First(&oldCustomer, "uuid = ?", id)
	if res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}

	if err := c.ShouldBindJSON(&newCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCustomer.UUID = oldCustomer.UUID
	newCustomer.NrPoints = oldCustomer.NrPoints
	newCustomer.CreatedAt = oldCustomer.CreatedAt

	res = ctrl.db.Save(newCustomer)
	if res.Error != nil {
		res.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error.Error()})
		return
	}

	res.Commit()
	c.JSON(http.StatusOK, gin.H{"status": "cliente atualizado com sucesso"})
}

func (ctrl *Controller) PostTransaction(c *gin.Context) {

	payload := &PayloadTransaction{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	customer := &models.Customer{}
	res := ctrl.db.First(&customer, "uuid = ?", payload.IdCustomer)
	if res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
		return
	}

	t := models.NewTransaction(payload.IdCustomer, payload.VlPoints, payload.DescSysOrigin)
	ps := []models.TransactionProduct{}

	for _, p := range payload.Products {
		ps = append(ps, *models.NewTransactionProduct(t.UUID, p.CodProduct, p.QtdeProduct, p.VlProduct))
	}

	if err := TransactionProcess(customer, t, ps, ctrl.db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "transacao criada e pontos do usuário atualizados"})
}

func (ctrl *Controller) MigrateCustomers(c *gin.Context) {

	payloads := &PayloadMigrationCustomer{}
	err := c.ShouldBindJSON(&payloads)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	c1, c2 := &models.Customer{}, &models.Customer{}

	res := ctrl.db.First(&c1, "uuid = ?", payloads.Customer1UUID)
	if res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer 1 not found"})
		return
	}

	res = ctrl.db.First(&c2, "uuid = ?", payloads.Customer2UUID)
	if res.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer 2 not found"})
		return
	}

	if err := models.MigrateCustomers(c1, c2, ctrl.db); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "customers migrated successfully"})

}

func NewController(db *gorm.DB) *Controller {
	return &Controller{db: db}
}
