package main

import (
	"log"
	"points/db"
	"points/handlers"
	"points/models"

	"github.com/gin-gonic/gin"
)

func main() {

	db, err := db.OpenDBConnection()
	if err != nil {
		log.Fatal("failed to connect database")
	}

	log.Println("Executando migrations...")
	db.AutoMigrate(&models.Customer{}, &models.Transaction{}, &models.TransactionProduct{}, &models.Product{})
	log.Println("ok")

	r := gin.Default()
	r.Use(gin.Recovery())

	r.GET("customers/:id", handlers.GetCustomerByID)
	r.GET("customers/", handlers.GetCustomers)
	r.POST("customers/", handlers.PostCustomer)
	r.PUT("customers/:id", handlers.PutCustomer)

	r.POST("/transactions", handlers.PostTransaction)

	r.Run("0.0.0.0:8081")

}
