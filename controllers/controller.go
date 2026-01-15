package controllers

import (
	"net/http"
	"points/models"
	"points/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Controller struct {
	PointsService *services.PointsService
}

func (ctrl *Controller) GetCustomers(c *gin.Context) {

	uuid := c.Query("uuid")
	cod_cpf := c.Query("cod_cpf")
	desc_email := c.Query("desc_email")
	twitch := c.Query("twitch")
	youtube := c.Query("youtube")
	blue_sky := c.Query("blue_sky")
	instagram := c.Query("instagram")

	customers := ctrl.PointsService.PointsRepo.GetCustomer(
		uuid,
		cod_cpf,
		desc_email,
		twitch,
		youtube,
		blue_sky,
		instagram,
	)

	if len(customers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "nenhum resultado encontrado"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

func (ctrl *Controller) GetCustomerByID(c *gin.Context) {

	id := c.Param("id")

	customer, err := ctrl.PointsService.PointsRepo.GetCustomerByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (ctrl *Controller) PostCustomer(c *gin.Context) {

	customer := &models.Customer{}

	newCustomer := models.NewCustomer()

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customer.UUID = newCustomer.UUID

	if err := ctrl.PointsService.PointsRepo.CreateCustomer(customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "created", "customer": customer})
}

func (ctrl *Controller) PutCustomer(c *gin.Context) {

	id := c.Param("id")

	oldCustomer, err := ctrl.PointsService.PointsRepo.GetCustomerByID(id)
	if err != nil {
		if err.Error() == gorm.ErrRecordNotFound.Error() {
			c.JSON(http.StatusNotFound, gin.H{"error": "usuário não encontrado"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newCustomer := &models.Customer{}
	if err := c.ShouldBindJSON(&newCustomer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newCustomer.UUID = oldCustomer.UUID
	newCustomer.NrPoints = oldCustomer.NrPoints
	newCustomer.CreatedAt = oldCustomer.CreatedAt

	if err := ctrl.PointsService.PointsRepo.UpdateCustomer(newCustomer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "cliente atualizado com sucesso"})
}

func (ctrl *Controller) PostTransaction(c *gin.Context) {

	payload := &services.PayloadTransaction{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.PointsService.MakeTransaction(payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "transacao criada e pontos do usuário atualizados"})
}

func (ctrl *Controller) MigrateCustomers(c *gin.Context) {

	payload := &services.PayloadMigrationCustomer{}
	err := c.ShouldBindJSON(&payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}

	if err := ctrl.PointsService.MigrateCustomers(payload); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "customers migrated successfully"})

}

func NewController(pointsService *services.PointsService) *Controller {
	return &Controller{PointsService: pointsService}
}
