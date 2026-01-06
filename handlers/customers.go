package handlers

import (
	"net/http"
	"points/models"
	"points/myerrors"

	"github.com/gin-gonic/gin"
)

func GetCustomers(c *gin.Context) {

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
		dbConnection,
	)

	if len(customers) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"status": "nenhum resultado encontrado"})
		return
	}

	c.JSON(http.StatusOK, customers)
}

func GetCustomerByID(c *gin.Context) {

	var customer models.Customer
	id := c.Param("id")

	result := dbConnection.First(&customer, "uuid = ?", id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "customer not found"})
		return
	}
	c.JSON(http.StatusOK, customer)
}

func PostCustomer(c *gin.Context) {

	customer := models.NewCustomer()
	customerUUID := customer.UUID

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	customer.UUID = customerUUID

	err := myerrors.GetCreateCustomerErrors(customer.Create(dbConnection))
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

func PutCustomer(c *gin.Context) {

	id := c.Param("id")

	var newCustomer, oldCustomer *models.Customer

	res := dbConnection.First(&oldCustomer, "uuid = ?", id)
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

	if newCustomer.DescCustomerName == "" {
		newCustomer.DescCustomerName = oldCustomer.DescCustomerName
	}

	if newCustomer.CodCPF == nil {
		newCustomer.CodCPF = oldCustomer.CodCPF
	}

	if newCustomer.DescEmail == nil {
		newCustomer.DescEmail = oldCustomer.DescEmail
	}

	if newCustomer.IdTwitch == nil {
		newCustomer.IdTwitch = oldCustomer.IdTwitch
	}

	if newCustomer.IdYouTube == nil {
		newCustomer.IdYouTube = oldCustomer.IdYouTube
	}

	if newCustomer.IdBlueSky == nil {
		newCustomer.IdBlueSky = oldCustomer.IdBlueSky
	}

	if newCustomer.IdInstagram == nil {
		newCustomer.IdInstagram = oldCustomer.IdInstagram
	}

	res = dbConnection.Save(newCustomer)
	if res.Error != nil {
		res.Rollback()
		c.JSON(http.StatusBadRequest, gin.H{"error": res.Error.Error()})
		return
	}

	res.Commit()
	c.JSON(http.StatusOK, gin.H{"status": "cliente atualizado com sucesso"})
}
